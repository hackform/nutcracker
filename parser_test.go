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
