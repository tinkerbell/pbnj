// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package ipmitool

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/server/httpsvr/evlog"
	"github.com/tinkerbell/pbnj/server/httpsvr/interfaces/power"
)

var (
	statusRE = regexp.MustCompile(`Chassis Power is (on|off)`)
	prompt   = `ipmitool> `
)

// Shell represents an open remote shell.
type Shell struct {
	Address string
	tx      *evlog.Tx

	cmd   *exec.Cmd
	stdin io.WriteCloser

	mu     sync.Mutex
	outBuf bytes.Buffer
	errBuf bytes.Buffer
	last   string
	status power.Status
	err    *multierror.Error

	prompt  chan struct{}
	exitCh  chan struct{}
	exitErr error
}

func (s *Shell) debugEvent(event string, fields ...interface{}) {
	fields = append([]interface{}{
		"device", s.Address,
	}, fields...)
	s.tx.Debug(event, fields...)
}

func (s *Shell) errorEvent(event string, fields ...interface{}) {
	fields = append([]interface{}{
		"device", s.Address,
	}, fields...)
	s.tx.Error(event, fields...)
}

func (s *Shell) infoEvent(event string, fields ...interface{}) {
	fields = append([]interface{}{
		"device", s.Address,
	}, fields...)
	s.tx.Info(event, fields...)
}

func (s *Shell) noticeEvent(event string, fields ...interface{}) {
	fields = append([]interface{}{
		"device", s.Address,
	}, fields...)
	s.tx.Notice(event, fields...)
}

// Shell starts a remote shell.
func (o Options) Shell(ctx context.Context) (*Shell, error) {
	cmd := o.buildCommand("shell")

	s := &Shell{
		Address: o.Address,
		tx:      elog.TxFromContext(ctx),
	}
	s.stdin, _ = cmd.StdinPipe()
	cmd.Stdout = funcWriter(s.stdout)
	cmd.Stderr = funcWriter(s.stderr)

	s.mu.Lock()
	s.debugEvent("shell_start")

	if err := cmd.Start(); err != nil {
		s.stdin = nil
		s.mu.Unlock()
		return nil, errors.Wrap(err, "error starting ipmitool shell")
	}
	go s.waitForExit()

	s.cmd = cmd
	s.exitCh = make(chan struct{})
	s.prompt = make(chan struct{}, 1)
	s.mu.Unlock()

	if err := s.waitForPrompt(); err != nil {
		_ = s.Close()
		return nil, errors.WithMessage(err, "error waiting for shell prompt")
	}
	// TODO(betawaffle): Add a finalizer to ensure we don't leak?
	// TODO(mmlb): finalizers are not guaranteed to be invoked, figure out a different way to ensure no leaks
	return s, nil
}

// Close stops the remote shell.
func (s *Shell) Close() error {
	return s.Run("quit")
}

// Run executes a command on the remote shell.
func (s *Shell) Run(cmd string) (err error) {
	s.mu.Lock()
	s.last = cmd
	s.mu.Unlock()

	fields := []interface{}{
		"device", s.Address,
		"cmd", cmd,
	}
	trace := s.tx.Trace("run_ipmi_cmd", fields...)
	defer trace.Stop(&err)

	fmt.Fprintln(s.stdin, cmd)

	return s.waitForPrompt()
}

// LastStatus returns the previous power status.
func (s *Shell) LastStatus() power.Status {
	s.mu.Lock()
	status := s.status
	s.mu.Unlock()
	return status
}

func (s *Shell) stdout(p []byte) (int, error) {
	s.mu.Lock()

	if VerboseDebug {
		s.debugEvent("stdout_read", "buf", p)
	}

	eachLine(&s.outBuf, p, s.stdoutLine)

	if pending := s.outBuf.Bytes(); len(pending) > 0 {
		s.stdoutPending(pending)
	}

	s.mu.Unlock()
	return len(p), nil
}

func (s *Shell) stdoutLine(line []byte) {
	if len(line) == 0 {
		return
	}
	if len(s.last) == len(line) && string(line) == s.last {
		s.last = ""
		return
	}

	if m := statusRE.FindSubmatch(line); m != nil {
		// FIXME(betawaffle): This is brittle if we decide to change the power.Status constants.
		status := power.Status(m[1])
		if status != s.status {
			s.status = status
			s.noticeEvent("power_is_" + string(status))
		} else {
			s.debugEvent("power_is_still_" + string(status))
		}
		return
	}

	s.infoEvent("stdout_line", "line", string(line))
}

func (s *Shell) stdoutPending(pending []byte) {
	if VerboseDebug {
		s.debugEvent("stdout_pending", "buf", pending)
	}

	if n := len(prompt); len(pending) >= n && prompt == string(pending[:n]) {
		s.outBuf.Next(n)
		select {
		case s.prompt <- struct{}{}:
		default:
		}
	}
}

func (s *Shell) stderr(p []byte) (int, error) {
	s.mu.Lock()

	if VerboseDebug {
		s.debugEvent("stderr_read", "buf", p)
	}

	eachLine(&s.errBuf, p, s.stderrLine)

	s.mu.Unlock()
	return len(p), nil
}

func (s *Shell) stderrLine(line []byte) {
	if len(line) == 0 {
		return
	}
	s.errorEvent("stderr_line", "line", string(line))

	if err := parseStderrLine(line); err != nil {
		s.err = multierror.Append(s.err, err)
	}
}

func (s *Shell) takeErr() error {
	s.mu.Lock()
	err := s.err
	s.err = nil
	s.mu.Unlock()

	if err == nil {
		return nil
	}
	switch len(err.Errors) {
	case 1:
		if errors.Is(err.Errors[0], ErrChannelCipherSuites) {
			// ignore cipher suites error, ipmitool sends this even though it successfully downgrades
			return nil
		}
		return errors.Wrap(err.Errors[0], "shell error")
	case 0:
		return nil
	}
	return errors.Wrap(err, "shell error")
}

func (s *Shell) waitForExit() {
	s.exitErr = s.cmd.Wait()

	var (
		logger = s.infoEvent
		proc   = s.cmd.ProcessState
		sys    = proc.SystemTime()
		usr    = proc.UserTime()
		fields = append(make([]interface{}, 0, 6), "system_time", sys, "user_time", usr)
	)
	if !proc.Success() {
		logger = s.errorEvent
		fields = append(fields, "err", proc)
	}
	logger("shell_exit", fields...)

	close(s.exitCh)
}

func (s *Shell) waitForPrompt() error {
	select {
	case <-s.prompt:
		return s.takeErr()
	case <-s.exitCh:
		if err := s.takeErr(); err != nil {
			return errors.WithMessage(err, "error waiting for exitCh prompt")
		}
		return errors.Wrap(s.exitErr, "exit error")
	}
}
