package nutcracker

type (
	Error int
)

const (
	_               = iota
	ErrUnclosedStrI = Error(iota)
	ErrUnclosedStrL
	ErrUnclosedParen
	ErrInvalidEscape
	ErrInvalidCloseParen
	ErrInvalidArgMode
	ErrInvalidExec
)

func (e Error) Error() string {
	switch e {
	case ErrUnclosedStrI:
		return "unclosed double quote"
	case ErrUnclosedStrL:
		return "unclosed single quote"
	case ErrUnclosedParen:
		return "unclosed parenthesis"
	case ErrInvalidEscape:
		return "invalid escape"
	case ErrInvalidCloseParen:
		return "invalid close parenthesis"
	case ErrInvalidArgMode:
		return "invalid argument mode"
	case ErrInvalidExec:
		return "invalid command to execute"
	default:
		return "nutcracker error"
	}
}
