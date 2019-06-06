package nutcracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_isSpace(t *testing.T) {
	assert := assert.New(t)

	assert.True(isSpace(' '), "space is a space")
	assert.True(isSpace('\n'), "newline is a space")
	assert.True(isSpace('\t'), "tab is a space")
	assert.True(isSpace('\r'), "carriage return is a space")
	assert.False(isSpace('a'), "alphanumeric is not a space")
}

func Test_trimLSpace(t *testing.T) {
	assert := assert.New(t)

	{
		arg := "\t\r\n hello "
		s := trimLSpace(arg)
		assert.Equal("hello ", s, "spaces should be removed from the left of the string only")
	}
}

func Test_unquoteArg(t *testing.T) {
	assert := assert.New(t)

	{
		arg := `hello\ world`
		s, err := unquoteArg(arg)
		assert.NoError(err, "unquote should not error")
		assert.Equal("hello world", s, "string should be unquoted")
	}
	{
		arg := `hello world\`
		_, err := unquoteArg(arg)
		assert.Equal(errInvalidEscape, err, "unquote should error on invalid escapes")
	}
	{
		arg := `hello\
 world\ `
		s, err := unquoteArg(arg)
		assert.NoError(err, "unquote should not error")
		assert.Equal("hello world ", s, "string should have newline removed")
	}
}

func Test_unquoteStrI(t *testing.T) {
	assert := assert.New(t)

	{
		arg := `hello\ world\$\"\\ `
		s, err := unquoteStrI(arg)
		assert.NoError(err, "unquote should not error")
		assert.Equal("hello\\ world$\"\\ ", s, "unquote should only be removed if before a special char")
	}
	{
		arg := `hello world\`
		_, err := unquoteStrI(arg)
		assert.Equal(errInvalidEscape, err, "unquote should error on invalid escapes")
	}
	{
		arg := `hello\
world\ `
		s, err := unquoteStrI(arg)
		assert.NoError(err, "unquote should not error")
		assert.Equal("helloworld\\ ", s, "string should have newline removed")
	}
}
