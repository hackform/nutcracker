package nutcracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Error_Error(t *testing.T) {
	assert := assert.New(t)

	assert.NotEqual("", ErrUnclosedStrI.Error(), "error should not be empty")
	assert.NotEqual("", ErrUnclosedStrL.Error(), "error should not be empty")
	assert.NotEqual("", ErrUnclosedParen.Error(), "error should not be empty")
	assert.NotEqual("", ErrInvalidEscape.Error(), "error should not be empty")
	assert.NotEqual("", ErrInvalidCloseParen.Error(), "error should not be empty")
	assert.NotEqual("", ErrInvalidArgMode.Error(), "error should not be empty")
	assert.NotEqual("", ErrInvalidExec.Error(), "error should not be empty")
	assert.NotEqual("", internalError(0).Error(), "error should not be empty")
}
