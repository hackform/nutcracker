package nutcracker

import (
	"strings"
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

func parseArg(text string) (*nodeArg, string, error) {
	nodes := []Node{}
	i := 0
	for i < len(text) {
		ch := text[i]
		if ch == '\\' {
			if i+1 >= len(text) {
				return nil, "", errInvalidEscape
			}
			i += 2
		} else if isSpace(ch) {
			if i > 0 {
				nodes = append(nodes, newNodeText(text[0:i]))
				text = text[i:]
			}
			text = trimLSpace(text)
			i = 0
			break
		} else if ch == '"' {
			if i > 0 {
				nodes = append(nodes, newNodeText(text[0:i]))
				text = text[i:]
			}
			n, next, err := parseDoubleQuote(text)
			if err != nil {
				return nil, "", err
			}
			nodes = append(nodes, n)
			text = next
			i = 0
		} else {
			i++
		}
	}
	if i > 0 {
		nodes = append(nodes, newNodeText(text[0:i]))
		text = text[i:]
		i = 0
	}
	return newNodeArg(nodes), text, nil
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

func parseDoubleQuote(text string) (*nodeStrI, string, error) {
	nodes := []Node{}
	i := 1
	found := false
	for i < len(text) {
		ch := text[i]
		if ch == '\\' {
			if i+1 >= len(text) {
				return nil, "", errInvalidEscape
			}
			i += 2
		} else if ch == '"' {
			found = true
			if i > 1 {
				s, err := unquoteStrI(text[1:i])
				if err != nil {
					return nil, "", err
				}
				nodes = append(nodes, nodeText{
					text: s,
				})
			}
			text = text[i+1:]
			i = 0
			break
		} else {
			i++
		}
	}
	if !found {
		return nil, "", errUnclosedStrI
	}
	return newNodeStrI(nodes), text, nil
}
