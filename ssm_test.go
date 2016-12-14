package ssm

import (
	"testing"

	testify "github.com/stretchr/testify/assert"
)

var (
	sm *StateMachine
)

func init() {
	Reset()
}

func Reset() {
	sm = New("a",
		Events{
			{"a-b", States{"a"}, "b"},
			{"a-loop", States{"a"}, "a"},
			{"b-c", States{"b"}, "c"},
		},
	)
}

func TestCan(t *testing.T) {
	assert := testify.New(t)

	Reset()

	assert.True(sm.Can("a-b"))
	assert.False(sm.Can("b-c"))
}

func TestTransition(t *testing.T) {
	assert := testify.New(t)

	Reset()

	assert.NoError(sm.Event("a-b"))
	assert.Equal("b", sm.Current())

	assert.NoError(sm.Event("b-c"))
	assert.Equal("c", sm.Current())
}

func TestLoopTransition(t *testing.T) {
	assert := testify.New(t)

	Reset()

	assert.NoError(sm.Event("a-loop"))
	assert.Equal("a", sm.Current())
}
