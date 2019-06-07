package nutcracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Parse(t *testing.T) {
	assert := assert.New(t)

	{
		arg := `echo $hello`
		n, err := Parse(arg)
		assert.NoError(err, "Parse should not error")
		assert.Equal(newNodeCmd([]Node{newNodeArg([]Node{newNodeText("echo")}), newNodeArg([]Node{newNodeEnvVar("hello", nil)})}, true), n, "only the first argument should be parsed")
		v, err := n.Value(Env{Envfunc: func(s string) string {
			if s == "hello" {
				return "world"
			}
			return ""
		}})
		assert.NoError(err, "cmd should not error")
		assert.Equal("world\n", v, "cmd output should be correct")
	}
	{
		arg := `echo $hello\`
		_, err := Parse(arg)
		assert.Equal(ErrInvalidEscape, err, "Parse should error on invalid argument")
	}
}
