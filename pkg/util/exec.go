package util

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ExecCmdWrapper struct {
	*exec.Cmd
	env []string
}

func NewExecCmdWrapper(cmd *exec.Cmd, env ...string) *ExecCmdWrapper {
	c := &ExecCmdWrapper{cmd, env}

	// Must, if os.Environ is not assigned to cmd.Env, some important os native environment variables may be lost,
	// and thus influence the cmd's execute behaviour
	c.Cmd.Env = os.Environ()
	c.Cmd.Env = append(c.Cmd.Env, env...)

	return c
}

// String returns a human-readable description of c.
// It have environment variables prefixed to the string of exec.Cmd
func (c *ExecCmdWrapper) String() string {
	b := new(strings.Builder)

	// don't loop c.Cmd.Env, it contains os native environment variables
	for _, e := range c.env {
		b.WriteString(e)
		b.WriteByte(' ')
	}

	b.WriteString(c.Cmd.String())

	return b.String()
}

// TimeoutExec like
func TimeoutExec(timeout int64, cmdname string, arg ...string) (stdout []byte, stderr []byte, err error) {
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, cmdname, arg...)

	// use this method to capture stdout and stderr
	// If don't want to capture, use os.Stdout, os.Stderr as the right value
	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e

	if err := cmd.Run(); err != nil {
		// We want to check the context error to see if the timeout was executed.
		// The error returned by cmd.Output() will be OS specific based on what
		// happens when a process is killed.
		if ctx.Err() == context.DeadlineExceeded {
			errInfo := fmt.Errorf("Command timed out")
			return nil, nil, errInfo
		}

		return o.Bytes(), e.Bytes(), err

	}

	return o.Bytes(), e.Bytes(), nil
}
