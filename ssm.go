package ssm

import "errors"

type StateMachine struct {
	current     state
	transitions transitions
}

type state interface{}
type States []state

type event interface{}

type transitions map[node]event

type node struct {
	event event
	from  state
}

type eventDesc struct {
	Event event
	From  States
	To    state
}

type Events []eventDesc

type loopEvent struct {
	Event event
	Stay  States
}

type LoopEvents []loopEvent

var (
	ErrInvalidTransition = errors.New("Invalid transition")
)

func New(initial state, events Events, loopEvents LoopEvents) *StateMachine {
	sm := StateMachine{
		current:     initial,
		transitions: make(transitions),
	}

	for _, e := range events {
		for _, from := range e.From {
			sm.transitions[node{e.Event, from}] = e.To
		}
	}

	for _, e := range loopEvents {
		for _, stay := range e.Stay {
			sm.transitions[node{e.Event, stay}] = stay
		}
	}

	return &sm
}

func (sm *StateMachine) Current() state {
	return sm.current
}

func (sm *StateMachine) Event(event event, args ...interface{}) error {
	dst, ok := sm.transitions[node{event, sm.Current()}]
	if !ok {
		return ErrInvalidTransition
	}

	if dst == sm.Current() {
		return nil
	}

	sm.current = dst
	return nil
}

func (sm *StateMachine) Can(event event) bool {
	_, ok := sm.transitions[node{event, sm.current}]
	return ok
}

func (sm *StateMachine) Cannot(event event) bool {
	return !sm.Can(event)
}
