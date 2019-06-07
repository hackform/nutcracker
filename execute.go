package nutcracker

import (
	"os/exec"
)

type (
	Executor interface {
		Exec(args []string, env Env) error
	}

	executor struct {
	}
)

func NewExecutor() Executor {
	return &executor{}
}

func (e executor) Exec(args []string, env Env) error {
	if args == nil || len(args) < 1 {
		return ErrInvalidExec
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = env.Stdin
	cmd.Stdout = env.Stdout
	cmd.Stderr = env.Stderr
	cmd.Env = env.Envvar
	return cmd.Run()
}
