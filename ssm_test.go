package gosm

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
			{"a-b", []string{"a"}, "b"},
			{"a-loop", []string{"a"}, "a"},
			{"b-c", []string{"b"}, "c"},
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
