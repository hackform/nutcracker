package nutcracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseArg(t *testing.T) {
	assert := assert.New(t)

	{
		arg := `hello world `
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("world ", next, "only the first argument should be parsed")
		assert.Equal(newNodeArg([]Node{newNodeText("hello")}), n, "only the first argument should be parsed")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal("hello", v, "value returns correct arg value")
	}
	{
		arg := `hello\ world\
! kevin `
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("kevin ", next, "escape will escape spaces and eliminate newline")
		assert.Equal(newNodeArg([]Node{newNodeText("hello world!")}), n, "escape will escape spaces and eliminate newline")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal("hello world!", v, "value returns correct arg value")
	}
	{
		arg := `hello\ world`
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("", next, "escape will escape spaces")
		assert.Equal(newNodeArg([]Node{newNodeText("hello world")}), n, "escape will escape spaces")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal("hello world", v, "value returns correct arg value")
	}
	{
		arg := `"hello\ 'world"`
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("", next, "interpolated string will include spaces and single quotes")
		assert.Equal(newNodeArg([]Node{newNodeStrI([]Node{newNodeText("hello\\ 'world")})}), n, "interpolated string will include spaces and single quotes")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal("hello\\ 'world", v, "value returns correct arg value")
	}
	{
		arg := `"hello\
'world" kevin `
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("kevin ", next, "interpolated string will eliminate escaped newline")
		assert.Equal(newNodeArg([]Node{newNodeStrI([]Node{newNodeText("hello'world")})}), n, "interpolated string will eliminate escaped newline")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal("hello'world", v, "value returns correct arg value")
	}
	{
		arg := `"hello\$ world"\$ kevin `
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("kevin ", next, "parse arg will include adjacent nodes")
		assert.Equal(newNodeArg([]Node{newNodeStrI([]Node{newNodeText("hello$ world")}), newNodeText("$")}), n, "parse arg will include adjacent nodes")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal("hello$ world$", v, "value returns correct arg value")
	}
	{
		arg := `'hello\$ world'\$ kevin `
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("kevin ", next, "text in literal quote remains unchanged")
		assert.Equal(newNodeArg([]Node{newNodeStrL("hello\\$ world"), newNodeText("$")}), n, "text in literal quote remains unchanged")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal("hello\\$ world$", v, "value returns correct arg value")
	}
	{
		arg := `$hello\ ${world}kevin `
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("", next, "all variables should be consumed")
		assert.Equal(newNodeArg([]Node{newNodeEnvVar("hello", nil), newNodeText(" "), newNodeEnvVar("world", nil), newNodeText("kevin")}), n, "text in literal quote remains unchanged")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal(" kevin", v, "value returns correct arg value")
		v, err = n.Value(Env{Envfunc: func(s string) string {
			if s == "hello" {
				return "kevin"
			} else if s == "world" {
				return "wang"
			}
			return ""
		}})
		assert.NoError(err, "node value should not error")
		assert.Equal("kevin wangkevin", v, "value returns correct arg value")
	}
	{
		arg := `${world:-  some   default      value}kevin `
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("", next, "all variables should be consumed")
		assert.Equal(newNodeArg([]Node{newNodeEnvVar("world", []Node{newNodeArg([]Node{newNodeText("some")}), newNodeArg([]Node{newNodeText("default")}), newNodeArg([]Node{newNodeText("value")})}), newNodeText("kevin")}), n, "default value is parsed by arguments")
		v, err := n.Value(Env{})
		assert.NoError(err, "node value should not error")
		assert.Equal("some default valuekevin", v, "value returns correct arg value")
	}
	{
		arg := `${world:-$hello}kevin`
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("", next, "all variables should be consumed")
		assert.Equal(newNodeArg([]Node{newNodeEnvVar("world", []Node{newNodeArg([]Node{newNodeEnvVar("hello", nil)})}), newNodeText("kevin")}), n, "default value is parsed as arg")
		v, err := n.Value(Env{Envfunc: func(s string) string {
			if s == "hello" {
				return "greetings"
			}
			return ""
		}})
		assert.NoError(err, "node value should not error")
		assert.Equal("greetingskevin", v, "value returns correct arg value")
	}
	{
		arg := `"${world:-$hello }  $hello"kevin`
		n, next, err := parseArg(arg, argModeNorm)
		assert.NoError(err, "parse arg should not error")
		assert.Equal("", next, "all variables should be consumed")
		assert.Equal(newNodeArg([]Node{newNodeStrI([]Node{newNodeEnvVar("world", []Node{newNodeArg([]Node{newNodeEnvVar("hello", nil)})}), newNodeText("  "), newNodeEnvVar("hello", nil)}), newNodeText("kevin")}), n, "args in strings are parsed")
		v, err := n.Value(Env{Envfunc: func(s string) string {
			if s == "hello" {
				return "greetings"
			}
			return ""
		}})
		assert.NoError(err, "node value should not error")
		assert.Equal("greetings  greetingskevin", v, "value returns correct arg value")
	}
	{
		arg := `hello\ world\`
		_, _, err := parseArg(arg, -1)
		assert.Equal(ErrInvalidArgMode, err, "parse arg should error on invalid mode")
	}
	{
		arg := `hello\ world\`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidEscape, err, "parse arg should error on invalid escape")
	}
	{
		arg := `hello\ $`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidVar, err, "parse arg should error on invalid var")
	}
	{
		arg := `hello\ ${hello`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrUnclosedBrace, err, "parse arg should error on invalid var")
	}
	{
		arg := `hello\ ${hello:-`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrUnclosedBrace, err, "parse arg should error on invalid var")
	}
	{
		arg := `hello\ "$"`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidVar, err, "parse arg should error on invalid var in string")
	}
	{
		arg := `hello\ "${"`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidVar, err, "parse arg should error on invalid var in string")
	}
	{
		arg := `hello\ "${hello"`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidVar, err, "parse arg should error on invalid var in string")
	}
	{
		arg := `hello\ "${hello:-`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrUnclosedBrace, err, "parse arg should error on invalid var in string")
	}
	{
		arg := `hello\ "${hello:- $}"`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidVar, err, "parse arg should error on invalid arg in default value")
	}
	{
		arg := `"hello\$ world\`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidEscape, err, "parse arg should error on invalid escape")
	}
	{
		arg := `hello) world`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidCloseParen, err, "parse arg should error on invalid mode")
	}
	{
		arg := `hello} world`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrInvalidCloseBrace, err, "parse arg should error on invalid mode")
	}
	{
		arg := `'hello\$ world\`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrUnclosedStrL, err, "parse arg should error on unclosed literal string")
	}
	{
		arg := `"hello\$ world`
		_, _, err := parseArg(arg, argModeNorm)
		assert.Equal(ErrUnclosedStrI, err, "parse arg should error on unclosed interpolated string")
	}
}

func Test_parseArgText(t *testing.T) {
	assert := assert.New(t)

	{
		arg := `hello \`
		_, _, err := parseArgText(arg, len(arg))
		assert.Equal(ErrInvalidEscape, err, "parse arg text should error on invalid escape")
	}
}
