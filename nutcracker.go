package nutcracker

type (
	token struct {
		id  int
		val string
	}

	matcher func(string) ([]token, string, error)
)

func tokenize(directive string) ([]*nodeArg, error) {
	args := []*nodeArg{}
	for text := trimLSpace(directive); len(text) > 0; text = trimLSpace(text) {
		t, next, err := parseArg(text, argModeNorm)
		if err != nil {
			return nil, err
		}
		args = append(args, t)
		text = next
	}
	return args, nil
}
