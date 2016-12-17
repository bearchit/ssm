package ssm

import "fmt"

type StateMachine struct {
	current     state
	transitions transitions
	cbEvent     eventCallbacks
	cbState     stateCallbacks
}

type state interface{}
type States []state

type event interface{}

type transitions map[node]state

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

type loopDesc struct {
	Event event
	Stay  States
}

type LoopEvents []loopDesc

type callbackFn func(...interface{}) error

type eventCallbacks map[int]map[event]callbackFn
type stateCallbacks map[int]map[state]callbackFn

type eventCallbackDesc struct {
	Type     int
	Event    event
	Callback callbackFn
}

type stateCallbackDesc struct {
	Type     int
	State    state
	Callback callbackFn
}

type EventCallbacks []eventCallbackDesc
type StateCallbacks []stateCallbackDesc

const (
	Before = iota
	After
	Enter
	Leave
)

func New(initial state, events Events, loopEvents LoopEvents,
	ecb EventCallbacks, scb StateCallbacks) *StateMachine {

	sm := StateMachine{
		current:     initial,
		transitions: make(transitions),
		cbEvent:     make(eventCallbacks),
		cbState:     make(stateCallbacks),
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

	for _, cb := range ecb {
		if sm.cbEvent[cb.Type] == nil {
			sm.cbEvent[cb.Type] = make(map[event]callbackFn)
		}
		sm.cbEvent[cb.Type][cb.Event] = cb.Callback
	}

	for _, cb := range scb {
		if sm.cbState[cb.Type] == nil {
			sm.cbState[cb.Type] = make(map[state]callbackFn)
		}
		sm.cbState[cb.Type][cb.State] = cb.Callback
	}

	return &sm
}

func (sm *StateMachine) Current() state {
	return sm.current
}

func (sm *StateMachine) Event(e event, args ...interface{}) error {
	dst, ok := sm.transitions[node{e, sm.Current()}]
	if !ok {
		return fmt.Errorf("Invalid transition current: %s, action: %s", sm.Current(), e)
	}

	if cb, ok := sm.cbEvent[Before][e]; ok {
		if err := cb(args...); err != nil {
			return err
		}
	}

	if cb, ok := sm.cbState[Enter][dst]; ok {
		if err := cb(args...); err != nil {
			return err
		}
	}

	if cb, ok := sm.cbState[Leave][sm.Current()]; ok {
		if err := cb(args...); err != nil {
			return err
		}
	}

	if dst == sm.Current() {
		return nil
	}

	sm.current = dst

	if cb, ok := sm.cbEvent[After][e]; ok {
		if err := cb(args...); err != nil {
			return err
		}
	}

	return nil
}

// TODO: Callback, Condition 구분 필요
func (sm *StateMachine) Can(e event, args ...interface{}) bool {
	dst, ok := sm.transitions[node{e, sm.Current()}]
	if !ok {
		return false
	}

	if cb, ok := sm.cbEvent[Before][e]; ok {
		if err := cb(args...); err != nil {
			return false
		}
	}

	if cb, ok := sm.cbState[Enter][dst]; ok {
		if err := cb(args...); err != nil {
			return false
		}
	}

	if cb, ok := sm.cbState[Leave][sm.Current()]; ok {
		if err := cb(args...); err != nil {
			return false
		}
	}

	return true
}

func (sm *StateMachine) Cannot(e event, args ...interface{}) bool {
	return !sm.Can(e, args...)
}
