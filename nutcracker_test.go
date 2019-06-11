package nutcracker

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Parse(t *testing.T) {
	assert := assert.New(t)

	exec := NewExecutor()
	{
		b := bytes.Buffer{}
		arg := `echo $hello`
		n, err := Parse(arg)
		assert.NoError(err, "Parse should not error")
		assert.Equal([]Node{newNodeArg([]Node{newNodeText("echo")}), newNodeArg([]Node{newNodeEnvVar("hello", nil)})}, n.args, "all arguments should be parsed")
		err = n.Exec(Env{Envfunc: func(s string) string {
			if s == "hello" {
				return "world"
			}
			return ""
		}, Ex: exec, Stdout: &b})
		assert.NoError(err, "cmd should not error")
		assert.Equal("world\n", b.String(), "cmd stdout output should be correct")
	}
	{
		arg := `echo $hello\`
		_, err := Parse(arg)
		assert.Equal(ErrInvalidEscape, err, "Parse should error on invalid argument")
	}
	{
		arg := `echo $(bogus)`
		n, err := Parse(arg)
		assert.NoError(err, "Parse should not error for valid syntax")
		err = n.Exec(Env{Ex: exec})
		assert.Error(err, "Parse should error on command error")
	}
	{
		b := bytes.Buffer{}
		arg := ``
		n, err := Parse(arg)
		assert.NoError(err, "Parse should not error for valid syntax")
		err = n.Exec(Env{Ex: exec, Stdout: &b})
		assert.NoError(err, "Exec should not error on empty command")
		assert.Equal("", b.String(), "Exec should not write to stdout on empty command")
	}
	{
		arg := `bogus hello`
		n, err := Parse(arg)
		assert.NoError(err, "Parse should not error for valid syntax")
		err = n.Exec(Env{Ex: exec})
		assert.Error(err, "Parse should error on command error")
	}
}
