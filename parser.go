package nutcracker

import (
	"bytes"
	"io"
	"strings"
)

const (
	argModeNorm = iota
	argModeCmd
	argModeSub
	argModeVar
)

type (
	EnvFunc func(string) string

	Env struct {
		Envvar  []string
		Envfunc EnvFunc
		Stdin   io.Reader
		Stdout  io.Writer
		Stderr  io.Writer
		Ex      Executor
	}

	Node interface {
		Value(env Env) (string, error)
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

func (n nodeText) Value(env Env) (string, error) {
	return n.text, nil
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

func (n nodeArg) Value(env Env) (string, error) {
	s := strings.Builder{}
	for _, i := range n.nodes {
		v, err := i.Value(env)
		if err != nil {
			return "", err
		}
		s.WriteString(v)
	}
	return s.String(), nil
}

// parseArg parses one argument in the current mode
// takes in a string not beginning with whitespace
func parseArg(text string, mode int) (*nodeArg, string, error) {
	switch mode {
	case argModeNorm, argModeCmd, argModeSub, argModeVar:
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
		} else if isSpace(ch) || ch == ')' || ch == '}' || ch == '"' || ch == '\'' || ch == '$' {
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
				switch mode {
				case argModeNorm, argModeVar:
					return nil, "", ErrInvalidCloseParen
				}
				break
			} else if ch == '}' {
				switch mode {
				case argModeNorm, argModeCmd, argModeSub:
					return nil, "", ErrInvalidCloseBrace
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
			} else if ch == '$' {
				n, next, err := parseVar(text)
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

// parseArgText consumes the first i bytes to create a text node
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

func (n nodeStrI) Value(env Env) (string, error) {
	s := strings.Builder{}
	for _, i := range n.nodes {
		v, err := i.Value(env)
		if err != nil {
			return "", err
		}
		s.WriteString(v)
	}
	return s.String(), nil
}

// parseStrI parses interpolated strings.
// takes in a string beginning with '"'
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
		} else if ch == '"' || ch == '$' {
			if i > 0 {
				s, err := unquoteStrI(text[0:i])
				if err != nil {
					return nil, "", err
				}
				nodes = append(nodes, newNodeText(s))
				text = text[i:]
				i = 0
			}
			if ch == '"' {
				text = text[1:]
				return newNodeStrI(nodes), text, nil
			} else if ch == '$' {
				n, next, err := parseVar(text)
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

func (n nodeStrL) Value(env Env) (string, error) {
	return n.text, nil
}

// parseStrL parses literal strings.
// takes in a string beginning with '\''
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

type (
	nodeEnvVar struct {
		name   string
		defval []Node
	}
)

func newNodeEnvVar(name string, defval []Node) *nodeEnvVar {
	return &nodeEnvVar{
		name:   name,
		defval: defval,
	}
}

func (n nodeEnvVar) Value(env Env) (string, error) {
	if env.Envfunc != nil {
		k := env.Envfunc(n.name)
		if len(k) > 0 {
			return k, nil
		}
	}
	if n.defval == nil {
		return "", nil
	}
	s := strings.Builder{}
	first := true
	for _, i := range n.defval {
		v, err := i.Value(env)
		if err != nil {
			return "", err
		}
		if first {
			first = false
		} else {
			s.WriteByte(' ')
		}
		s.WriteString(v)
	}
	return s.String(), nil
}

// parseVar parses env vars and command substitutions.
// takes in a string beginning with '$'
func parseVar(text string) (Node, string, error) {
	if len(text) < 2 {
		return nil, "", ErrInvalidVar
	}
	k := parseTopEnvVar(text[1:])
	if k > 0 {
		text = text[1:]
		name := text[0:k]
		text = text[k:]
		return newNodeEnvVar(name, nil), text, nil
	}
	ch := text[1]
	if ch == '{' {
		return parseVarLong(text)
	} else if ch == '(' {
		return parseCmd(text)
	}
	return nil, "", ErrInvalidVar
}

// parseVarLong parses long long env vars.
// takes in a string beginning with '${'
func parseVarLong(text string) (Node, string, error) {
	text = text[2:]
	k := parseTopEnvVar(text)
	name := text[0:k]
	text = text[k:]
	if len(text) < 1 {
		return nil, "", ErrUnclosedBrace
	}
	if text[0] == '}' {
		text = text[1:]
		return newNodeEnvVar(name, nil), text, nil
	}
	if len(text) < 2 || text[0:2] != ":-" {
		return nil, "", ErrInvalidVar
	}

	nodes := []Node{}
	text = trimLSpace(text[2:])
	for len(text) > 0 {
		ch := text[0]
		if ch == '}' {
			text = text[1:]
			return newNodeEnvVar(name, nodes), text, nil
		}
		n, next, err := parseArg(text, argModeVar)
		if err != nil {
			return nil, "", err
		}
		nodes = append(nodes, n)
		text = next
	}
	return nil, "", ErrUnclosedBrace
}

type (
	nodeCmd struct {
		nodes []Node
		top   bool
	}
)

func newNodeCmd(nodes []Node, top bool) *nodeCmd {
	return &nodeCmd{
		nodes: nodes,
		top:   top,
	}
}

func (n nodeCmd) Value(env Env) (string, error) {
	if len(n.nodes) == 0 {
		return "", nil
	}
	k := make([]string, 0, len(n.nodes))
	for _, i := range n.nodes {
		v, err := i.Value(env)
		if err != nil {
			return "", err
		}
		k = append(k, v)
	}
	b := bytes.Buffer{}
	if !n.top {
		env.Stdout = &b
	}
	if err := env.Ex.Exec(k, env); err != nil {
		return "", err
	}
	if !n.top {
		return parseTextNodes(b.String()), nil
	}
	return "", nil
}

func parseTextNodes(text string) string {
	nodes := []string{}
	text = trimLSpace(text)
	for len(text) > 0 {
		k := nextSpace(text)
		if k < 0 {
			k = len(text)
		}
		nodes = append(nodes, text[0:k])
		text = trimLSpace(text[k:])
	}
	return strings.Join(nodes, " ")
}

// parseCmd parses a command substitution.
// takes in a string beginning with '$('
func parseCmd(text string) (Node, string, error) {
	text = trimLSpace(text[2:])
	nodes := []Node{}
	for len(text) > 0 {
		ch := text[0]
		if ch == ')' {
			text = text[1:]
			return newNodeCmd(nodes, false), text, nil
		}
		n, next, err := parseArg(text, argModeCmd)
		if err != nil {
			return nil, "", err
		}
		nodes = append(nodes, n)
		text = next
	}
	return nil, "", ErrUnclosedParen
}
