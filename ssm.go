package gosm

import "fmt"

type StateMachine struct {
	current     string
	transitions transitions
}

type transitions map[node]string

type event struct {
	from string
	to   string
}

type node struct {
	event string
	from  string
}

type eventDesc struct {
	Event string
	From  []string
	To    string
}

type Events []eventDesc

func New(initial string, events Events) *StateMachine {
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

func (sm *StateMachine) Current() string {
	return sm.current
}

func (sm *StateMachine) Event(event string, args ...interface{}) error {
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

func (sm *StateMachine) Can(event string) bool {
	_, ok := sm.transitions[node{event, sm.current}]
	return ok
}

func (sm *StateMachine) Cannot(event string) bool {
	return !sm.Can(event)
}
