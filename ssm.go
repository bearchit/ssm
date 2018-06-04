package ssm

type StateMachine struct {
	current     State
	transitions transitions
	cbEvent     eventCallbacks
	cbState     stateCallbacks
}

type State interface{}
type States []State

type Event interface{}

type transitions map[node]State

type node struct {
	event Event
	from  State
}

type eventDesc struct {
	Event Event
	From  States
	To    State
}

type Events []eventDesc

type loopDesc struct {
	Event Event
	Stay  States
}

type LoopEvents []loopDesc

type callbackFn func(current State, args ...interface{}) error

type eventCallbacks map[int]map[Event]callbackFn
type stateCallbacks map[int]map[State]callbackFn

type eventCallbackDesc struct {
	Type     int
	Event    Event
	Callback callbackFn
}

type stateCallbackDesc struct {
	Type     int
	State    State
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

func New(initial State, events Events, loopEvents LoopEvents,
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
			sm.cbEvent[cb.Type] = make(map[Event]callbackFn)
		}
		sm.cbEvent[cb.Type][cb.Event] = cb.Callback
	}

	for _, cb := range scb {
		if sm.cbState[cb.Type] == nil {
			sm.cbState[cb.Type] = make(map[State]callbackFn)
		}
		sm.cbState[cb.Type][cb.State] = cb.Callback
	}

	return &sm
}

func (sm *StateMachine) Current() State {
	return sm.current
}

func (sm *StateMachine) Event(e Event, args ...interface{}) error {
	dst, ok := sm.transitions[node{e, sm.Current()}]
	if !ok {
		return &InvalidTransitionError{Event: e, From: sm.Current()}
	}

	if cb, ok := sm.cbEvent[Before][e]; ok {
		if err := cb(sm.Current(), args...); err != nil {
			return err
		}
	}

	if cb, ok := sm.cbState[Enter][dst]; ok {
		if err := cb(sm.Current(), args...); err != nil {
			return err
		}
	}

	if cb, ok := sm.cbState[Leave][sm.Current()]; ok {
		if err := cb(sm.Current(), args...); err != nil {
			return err
		}
	}

	if dst == sm.Current() {
		return nil
	}

	sm.current = dst

	if cb, ok := sm.cbEvent[After][e]; ok {
		if err := cb(sm.Current(), args...); err != nil {
			return err
		}
	}

	return nil
}

// TODO: Callback, Condition 구분 필요
func (sm *StateMachine) Can(e Event, args ...interface{}) (bool, error) {
	dst, ok := sm.transitions[node{e, sm.Current()}]
	if !ok {
		return false, &InvalidTransitionError{Event: e, From: sm.Current()}
	}

	if cb, ok := sm.cbEvent[Before][e]; ok {
		if err := cb(sm.Current(), args...); err != nil {
			return false, err
		}
	}

	if cb, ok := sm.cbState[Enter][dst]; ok {
		if err := cb(sm.Current(), args...); err != nil {
			return false, err
		}
	}

	if cb, ok := sm.cbState[Leave][sm.Current()]; ok {
		if err := cb(sm.Current(), args...); err != nil {
			return false, err
		}
	}

	return true, nil
}
