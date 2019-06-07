package nutcracker

import (
	"strings"
)

const (
	argModeNorm = iota
	argModeCmd
	argModeSub
)

type (
	Node interface {
		Value() string
	}
)

type (
	nodeText struct {
		text string
	}
)

func newNodeText(text string) *nodeText {
	return &nodeText{
		text: text,
	}
}

func (n nodeText) Value() string {
	return n.text
}

type (
	nodeArg struct {
		nodes []Node
	}
)

func newNodeArg(nodes []Node) *nodeArg {
	return &nodeArg{
		nodes: nodes,
	}
}

func (n nodeArg) Value() string {
	s := strings.Builder{}
	for _, i := range n.nodes {
		s.WriteString(i.Value())
	}
	return s.String()
}

func parseArg(text string, mode int) (*nodeArg, string, error) {
	switch mode {
	case argModeNorm, argModeCmd, argModeSub:
	default:
		return nil, "", ErrInvalidArgMode
	}

	nodes := []Node{}
	i := 0
	for i < len(text) {
		ch := text[i]
		if ch == '\\' {
			if i+1 >= len(text) {
				return nil, "", ErrInvalidEscape
			}
			i += 2
		} else if isSpace(ch) || ch == ')' || ch == '"' || ch == '\'' {
			if i > 0 {
				n, next, err := parseArgText(text, i)
				if err != nil {
					return nil, "", err
				}
				nodes = append(nodes, n)
				text = next
				i = 0
			}
			if ch == ')' {
				if mode == argModeNorm {
					return nil, "", ErrInvalidCloseParen
				}
				break
			} else if isSpace(ch) {
				text = trimLSpace(text)
				break
			} else if ch == '"' {
				n, next, err := parseStrI(text)
				if err != nil {
					return nil, "", err
				}
				nodes = append(nodes, n)
				text = next
			} else if ch == '\'' {
				n, next, err := parseStrL(text)
				if err != nil {
					return nil, "", err
				}
				nodes = append(nodes, n)
				text = next
			}
		} else {
			i++
		}
	}

	if i > 0 {
		n, next, err := parseArgText(text, i)
		if err != nil {
			return nil, "", err
		}
		nodes = append(nodes, n)
		text = next
		i = 0
	}

	return newNodeArg(nodes), text, nil
}

func parseArgText(text string, i int) (*nodeText, string, error) {
	k, err := unquoteArg(text[0:i])
	if err != nil {
		return nil, "", err
	}
	return newNodeText(k), text[i:], nil
}

type (
	nodeStrI struct {
		nodes []Node
	}
)

func newNodeStrI(nodes []Node) *nodeStrI {
	return &nodeStrI{
		nodes: nodes,
	}
}

func (n nodeStrI) Value() string {
	s := strings.Builder{}
	for _, i := range n.nodes {
		s.WriteString(i.Value())
	}
	return s.String()
}

func parseStrI(text string) (*nodeStrI, string, error) {
	nodes := []Node{}
	text = text[1:]
	i := 0
	for i < len(text) {
		ch := text[i]
		if ch == '\\' {
			if i+1 >= len(text) {
				return nil, "", ErrInvalidEscape
			}
			i += 2
		} else if ch == '"' {
			if i > 0 {
				s, err := unquoteStrI(text[0:i])
				if err != nil {
					return nil, "", err
				}
				nodes = append(nodes, newNodeText(s))
			}
			text = text[i+1:]
			i = 0
			return newNodeStrI(nodes), text, nil
		} else {
			i++
		}
	}
	return nil, "", ErrUnclosedStrI
}

type (
	nodeStrL struct {
		text string
	}
)

func newNodeStrL(s string) *nodeStrL {
	return &nodeStrL{
		text: s,
	}
}

func (n nodeStrL) Value() string {
	return n.text
}

func parseStrL(text string) (*nodeStrL, string, error) {
	text = text[1:]
	i := 0
	for i < len(text) {
		ch := text[i]
		if ch == '\'' {
			k := text[0:i]
			text = text[i+1:]
			return newNodeStrL(k), text, nil
		} else {
			i++
		}
	}
	return nil, "", ErrUnclosedStrL
}
