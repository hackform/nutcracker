package nutcracker

import (
	"os/exec"
)

type (
	Executor interface {
		Exec() error
	}

	executor struct {
		cmd *exec.Cmd
	}
)

func newExecutor(args []string, env Env) (Executor, error) {
	if args == nil || len(args) < 1 {
		return nil, ErrInvalidExec
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = env.Stdin
	cmd.Stdout = env.Stdout
	cmd.Stderr = env.Stderr
	cmd.Env = env.Envvar
	return &executor{
		cmd: cmd,
	}, nil
}

func (e *executor) Exec() error {
	return e.cmd.Run()
}
