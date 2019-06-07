package nutcracker

type (
	Cmd struct {
		node Node
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
		node: newNodeCmd(args, true),
	}, nil
}

func (c Cmd) Exec(env Env) error {
	_, err := c.node.Value(env)
	return err
}
