package nutcracker

func Parse(shellcmd string) (Node, error) {
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
	return newNodeCmd(args, true), nil
}
