package nutcracker

type (
	Error int
)

const (
	_               = iota
	ErrUnclosedStrI = Error(iota)
	ErrUnclosedStrL
	ErrInvalidEscape
	ErrInvalidCloseParen
	ErrInvalidArgMode
)

func (e Error) Error() string {
	switch e {
	case ErrUnclosedStrI:
		return "unclosed double quote"
	case ErrUnclosedStrL:
		return "unclosed single quote"
	case ErrInvalidEscape:
		return "invalid escape"
	case ErrInvalidCloseParen:
		return "invalid close parenthesis"
	case ErrInvalidArgMode:
		return "invalid argument mode"
	default:
		return "nutcracker error"
	}
}
