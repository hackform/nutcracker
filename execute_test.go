package nutcracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Executor_Exec(t *testing.T) {
	assert := assert.New(t)

	{
		exec := newExecutor()
		err := exec.Exec([]string{}, Env{})
		assert.Equal(ErrInvalidExec, err, "executor must be run with a command")
	}
}
