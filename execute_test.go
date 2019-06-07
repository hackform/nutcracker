package nutcracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newExecutor(t *testing.T) {
	assert := assert.New(t)

	{
		_, err := newExecutor([]string{}, Env{})
		assert.Equal(ErrInvalidExec, err, "executor must be initialized with a command")
	}
}
