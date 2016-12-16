package ssm

import (
	"testing"

	"fmt"

	"reflect"

	testify "github.com/stretchr/testify/assert"
)

var (
	sm *StateMachine
)

func init() {
	Reset()
}

type State string
type Event string

const (
	EventAtoB Event = "a-b"
	EventBtoC Event = "b-c"
	EventLoop Event = "loop"

	StateA State = "a"
	StateB State = "b"
	StateC State = "c"
)

func Reset() {
	sm = New(StateA,
		Events{
			{EventAtoB, States{StateA}, StateB},
			{EventBtoC, States{StateB}, StateC},
		},
		LoopEvents{
			{EventLoop, States{StateA, StateB}},
		},
		EventCallbacks{
			{Type: Before, Event: EventAtoB, Callback: func(args ...interface{}) error {
				fmt.Println("before_a-b")
				return nil
			}},
			{Type: After, Event: EventAtoB, Callback: func(args ...interface{}) error {
				fmt.Println("after_a-b")
				return nil
			}},
		},
		StateCallbacks{
			{Type: Enter, State: StateB, Callback: func(args ...interface{}) error {
				fmt.Println("enter_b")
				return nil
			}},
			{Type: Leave, State: StateB, Callback: func(args ...interface{}) error {
				fmt.Println("leave_b")
				return nil
			}},
		},
	)
}

func TestCan(t *testing.T) {
	assert := testify.New(t)

	Reset()

	assert.True(sm.Can(EventAtoB))
	assert.False(sm.Can(EventBtoC))
}

func TestTransition(t *testing.T) {
	assert := testify.New(t)

	Reset()

	assert.NoError(sm.Event(EventAtoB))
	assert.Equal(StateB, sm.Current())

	assert.NoError(sm.Event(EventBtoC))
	assert.Equal(StateC, sm.Current())
}

func TestLoopTransition(t *testing.T) {
	assert := testify.New(t)

	Reset()

	assert.NoError(sm.Event(EventLoop))
	assert.Equal(StateA, sm.Current())

	assert.NoError(sm.Event(EventAtoB))
	assert.Equal(StateB, sm.Current())

	assert.NoError(sm.Event(EventLoop))
	assert.Equal(StateB, sm.Current())
}

func TestCustomTypeEquality(t *testing.T) {
	assert := testify.New(t)

	assert.NotEqual(StateA, "a")
	assert.Equal(StateA, State("a"))

	for k := range sm.cbEvent[Before] {
		t.Log(reflect.TypeOf(k))
	}
}
