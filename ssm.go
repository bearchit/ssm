package ssm

import "fmt"

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

func New(initial state, events Events) *StateMachine {
	sm := StateMachine{
		current:     initial,
		transitions: make(transitions),
	}

	for _, e := range events {
		for _, from := range e.From {
			sm.transitions[node{e.Event, from}] = e.To
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
		return fmt.Errorf("Invalid transition: current: %s, event: %s", event, sm.Current())
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
