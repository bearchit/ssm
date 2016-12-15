package ssm

import (
	"testing"

	"fmt"

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
			{"b-c", States{"b"}, "c"},
		},
		LoopEvents{
			{"loop", States{"a", "b"}},
		},
		Callbacks{
			{"before_a-b", func(args ...interface{}) error {
				fmt.Println("before_a-b")
				return nil
			}},
			{"after_a-b", func(args ...interface{}) error {
				fmt.Println("after_a-b")
				return nil
			}},
			{"enter_b", func(args ...interface{}) error {
				fmt.Println("enter_b")
				return nil
			}},
			{"leave_b", func(args ...interface{}) error {
				fmt.Println("leave_b")
				return nil
			}},
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

	assert.NoError(sm.Event("loop"))
	assert.Equal("a", sm.Current())

	assert.NoError(sm.Event("a-b"))
	assert.Equal("b", sm.Current())

	assert.NoError(sm.Event("loop"))
	assert.Equal("b", sm.Current())
}
