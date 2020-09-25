// Copyright 2020 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package racadm

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/tinkerbell/pbnj/evlog"
	"github.com/tinkerbell/pbnj/interfaces/power"
	"golang.org/x/crypto/ssh"
)

// Shell represents an open remote shell
type Shell struct {
	Address string
	tx      *evlog.Tx
	client  *ssh.Client

	lock       sync.RWMutex
	lastStatus power.Status
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

func (opts Options) AuthKeyboardInteractive(user, instruction string, questions []string, echos []bool) ([]string, error) {
	answers := make([]string, len(questions))
	for i := range answers {
		answers[i] = opts.Password
	}

	return answers, nil
}

// Shell starts a remote shell
func (opts Options) Shell(ctx context.Context) (*Shell, error) {
	s := &Shell{
		Address: opts.Address,
		tx:      elog.TxFromContext(ctx),
	}

	config := &ssh.ClientConfig{
		User: opts.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(opts.Password),
			ssh.KeyboardInteractive(opts.AuthKeyboardInteractive),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", opts.Address+":22", config)
	if err != nil {
		s.errorEvent("failed to dial", "action", "dial", "error", err)
		return nil, errors.Wrap(err, "failed to dial")
	}
	s.client = client

	s.debugEvent("shell_start")
	return s, nil
}

// Close stops the remote shell
func (s *Shell) Close() error {
	return s.client.Close()
}

// Run executes a command on the remote shell
func (s *Shell) Run(cmd string) error {
	var err error
	defer s.tx.Trace("run_racadm_cmd", "device", s.Address, "cmd", cmd).Stop(&err)

	session, err := s.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create new ssh client session")
	}
	defer func() {
		_ = session.Run("exit")
		_ = session.Close()
	}()

	outB, err := session.CombinedOutput(cmd)
	out := string(outB)
	if err != nil {
		if out != "" {
			return errors.Wrap(err, strings.TrimSpace(out))
		}
		return errors.Wrap(err, "error combining ssh session output")
	}

	s.infoEvent("output", "out", out)
	if strings.HasPrefix("ERROR: ", out) {
		return errors.New(out)
	}

	return err
}

// Output executes a command on the remote shell and returns any output
func (s *Shell) Output(cmd string) (string, error) {
	var err error
	defer s.tx.Trace("run_racadm_cmd", "device", s.Address, "cmd", cmd).Stop(&err)

	session, err := s.client.NewSession()
	if err != nil {
		return "", errors.Wrap(err, "failed to create new ssh client session")
	}
	defer func() {
		_ = session.Run("exit")
		_ = session.Close()
	}()

	outB, err := session.CombinedOutput(cmd)
	out := string(outB)
	if err != nil {
		if out != "" {
			return "", errors.Wrap(err, strings.TrimSpace(out))
		}
		return "", errors.Wrap(err, "error combining ssh session output")
	}

	s.infoEvent("output", "out", out)
	if strings.HasPrefix(out, "ERROR:") {
		err = errors.New(strings.TrimSpace(out))
		return "", err
	}

	if strings.HasPrefix(out, "[Key=:") {
		// ignore lines that begin with [Key= because we don't care about them for the commands we call
		index := strings.Index(out, "\n")
		if index != -1 {
			out = out[index+1:]
		}
	}

	return out, nil
}
