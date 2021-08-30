// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package power

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	// PollingInterval specifies how often to poll for power action.
	PollingInterval = 250 * time.Millisecond
)

// Operation specifies which operations are available.
type Operation func(context.Context, Driver, Options) error

// OperationBySlug lists the names of available operations.
var OperationBySlug = map[string]Operation{
	"cycle":    doCycle,
	"reset":    doReset,
	"soft":     doSoftOff,
	"soft_off": doSoftOff,
	"hard_off": doHardOff,
	"off":      doTurnOff,
	"turn_off": doTurnOff,
	"on":       doTurnOn,
	"turn_on":  doTurnOn,
}

// UnmarshalText unmarshals an Operation from a textual representation.
func (o *Operation) UnmarshalText(text []byte) error {
	if v, ok := OperationBySlug[string(text)]; ok {
		*o = v
		return nil
	}
	return errors.Errorf("unsupported power action: %q", text)
}

func doAction(ctx context.Context, driver Driver, action Action, timeout time.Duration, target Status) error {
	if err := driver.Power(action); err != nil {
		return errors.WithMessage(err, "unable to initiate power action")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return waitFor(ctx, driver, target)
}

func doCycle(ctx context.Context, driver Driver, opts Options) (err error) {
	defer elog.TxFromContext(ctx).Trace("power_cycle").Stop(&err)

	status, err := driver.PowerStatus()
	if err != nil {
		return errors.WithMessage(err, "error retrieving power status")
	}

	if status != Off {
		if err := doTurnOff(ctx, driver, opts); err != nil {
			return errors.WithMessage(err, "error turning off")
		}

		doSleep(ctx, opts.OffDuration)

		if isDone(ctx) {
			return ctx.Err()
		}
	}

	return doTurnOn(ctx, driver, opts)
}

func doHardOff(ctx context.Context, driver Driver, opts Options) (err error) {
	defer elog.TxFromContext(ctx).Trace("power_hard_off").Stop(&err)

	return doAction(ctx, driver, HardOff, opts.OffTimeout, Off)
}

func doReset(ctx context.Context, driver Driver, _ Options) (err error) {
	defer elog.TxFromContext(ctx).Trace("power_reset").Stop(&err)

	return driver.Power(Reset)
}

func doSoftOff(ctx context.Context, driver Driver, opts Options) (err error) {
	defer elog.TxFromContext(ctx).Trace("power_soft_off").Stop(&err)

	return doAction(ctx, driver, SoftOff, opts.SoftTimeout, Off)
}

func doTurnOff(ctx context.Context, driver Driver, opts Options) (err error) {
	defer elog.TxFromContext(ctx).Trace("power_turn_off").Stop(&err)

	if err := doSoftOff(ctx, driver, opts); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return errors.WithMessage(err, "error initiating soft off")
	}

	if isDone(ctx) {
		return ctx.Err()
	}

	return doHardOff(ctx, driver, opts)
}

func doTurnOn(ctx context.Context, driver Driver, opts Options) (err error) {
	defer elog.TxFromContext(ctx).Trace("power_turn_on").Stop(&err)

	return doAction(ctx, driver, TurnOn, opts.OnTimeout, On)
}

func doSleep(ctx context.Context, duration time.Duration) {
	elog.TxFromContext(ctx).Info("pausing", "duration", duration)

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	<-ctx.Done()
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func attemPowerStatus(ctx context.Context, driver Driver, delay time.Duration) error {
	var err error
	for attempts := 0; attempts < 30; attempts++ {
		_, err = driver.PowerStatus()
		if err == nil || errors.Cause(err).Error() != "Unable to get Chassis Power Status" {
			break
		}
		delay += 250 * time.Millisecond
		doSleep(ctx, delay)

		if isDone(ctx) {
			return ctx.Err()
		}
	}

	if err != nil {
		return errors.WithMessage(err, "error attempting power status")
	}

	return err
}

func waitFor(ctx context.Context, driver Driver, target Status) (err error) {
	defer elog.TxFromContext(ctx).Trace("wait_for_power_" + string(target)).Stop(&err)

	switch target {
	case Off:
		doSleep(ctx, 9*time.Second)
	case On:
		doSleep(ctx, 10*time.Second)
	default:
		return fmt.Errorf("unexpected target: %v", target)
	}

	ticker := time.NewTicker(PollingInterval) // Throttle calls to the driver.
	defer ticker.Stop()

poll:
	delay := 500 * time.Millisecond
	err = attemPowerStatus(ctx, driver, delay)
	if err != nil {
		return err
	}

	if target == AnyStatus || target == driver.LastStatus() {
		return nil
	}

	select {
	case <-ticker.C:
		goto poll
	case <-ctx.Done():
		return ctx.Err()
	}
}
