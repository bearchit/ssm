package ssm

import (
	"fmt"
	"regexp"
)

type StateMachine struct {
	current     state
	transitions transitions
	cbEvent     eventCallbacks
	cbState     stateCallbacks
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

type loopDesc struct {
	Event event
	Stay  States
}

type LoopEvents []loopDesc

type callbackFn func(...interface{}) error

type eventCallbacks map[string]map[event]callbackFn
type stateCallbacks map[string]map[state]callbackFn

type callbackDesc struct {
	Name     string
	Callback callbackFn
}

type Callbacks []callbackDesc

var (
	patternEventCallback = regexp.MustCompile(`(before|after)_(.+)`)
	patternStateCallback = regexp.MustCompile(`(enter|leave)_(.+)`)
)

func New(initial state, events Events, loopEvents LoopEvents, callbacks Callbacks) *StateMachine {
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

	for _, cb := range callbacks {
		r := patternEventCallback.FindStringSubmatch(cb.Name)
		if len(r) == 3 {
			if sm.cbEvent[r[1]] == nil {
				sm.cbEvent[r[1]] = make(map[event]callbackFn)
			}
			sm.cbEvent[r[1]][r[2]] = cb.Callback
		} else {
			r := patternStateCallback.FindStringSubmatch(cb.Name)
			if len(r) == 3 {
				if sm.cbState[r[1]] == nil {
					sm.cbState[r[1]] = make(map[state]callbackFn)
				}
				sm.cbState[r[1]][r[2]] = cb.Callback
			}
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
		return fmt.Errorf("Invalid transition current: %s, action: %s", sm.Current(), event)
	}

	if cb, ok := sm.cbEvent["before"][event]; ok {
		if err := cb(args...); err != nil {
			return err
		}
	}

	if cb, ok := sm.cbState["enter"][dst]; ok {
		if err := cb(args...); err != nil {
			return err
		}
	}

	if cb, ok := sm.cbState["leave"][sm.Current()]; ok {
		if err := cb(args...); err != nil {
			return err
		}
	}

	if dst == sm.Current() {
		return nil
	}

	sm.current = dst

	if cb, ok := sm.cbEvent["after"][event]; ok {
		if err := cb(args...); err != nil {
			return err
		}
	}

	return nil
}

func (sm *StateMachine) Can(event event) bool {
	_, ok := sm.transitions[node{event, sm.current}]
	return ok
}

func (sm *StateMachine) Cannot(event event) bool {
	return !sm.Can(event)
}
