package nutcracker

type (
	Cmd struct {
		args []Node
	}
)

func Parse(shellcmd string) (*Cmd, error) {
	args := []Node{}
	text := trimLSpace(shellcmd)
	for len(text) > 0 {
		n, next, err := parseArg(text, argModeNorm)
		if err != nil {
			return nil, err
		}
		args = append(args, n)
		text = next
	}
	return &Cmd{
		args: args,
	}, nil
}

func (c Cmd) Exec(env Env) error {
	if len(c.args) == 0 {
		return nil
	}
	k := make([]string, 0, len(c.args))
	for _, i := range c.args {
		v, err := i.Value(env)
		if err != nil {
			return err
		}
		k = append(k, v)
	}
	if err := env.Ex.Exec(k, env); err != nil {
		return err
	}
	return nil
}
