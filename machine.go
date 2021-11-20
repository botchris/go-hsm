package hsm

import (
	"fmt"
	"reflect"
	"sync"
)

// HSM represents a finite state machine
type HSM struct {
	// human-readable name of this machine
	name string

	// serves as machine's extended states
	// see: https://en.wikipedia.org/wiki/UML_state_machine#Extended_states
	context interface{}

	// list of states within this machine
	states map[string]*Vertex

	// pointer to the current state
	currentState *Vertex

	// pointer to a state that will be entered whenever an error occurs in the state machine
	errorState *Vertex

	// holds a history of (successfully) triggered signals in this HSM
	signalsHistory []string

	// holds a sequence history of states this HSM has been passing through
	statesHistory []string

	// guards access to HSM Signal() method
	signalMutex sync.RWMutex

	// guards to HSM current state
	currentMutex sync.RWMutex
}

// Snapshot provides a public snapshot
type Snapshot struct {
	// Current state ID
	StateID string

	// Whether current state is final or not
	Final bool

	// History of signals applied to this HSM
	SignalsHistory []string

	// History of states this HSM been at
	StatesHistory []string
}

// Current retrieves HSM`s current state
func (h *HSM) Current() *Vertex {
	h.currentMutex.RLock()
	defer h.currentMutex.RUnlock()

	return h.currentState
}

// At returns true if HSM is currently at the given state, false otherwise.
//
// This method returns true when asserting a parent state and HSM is currently at any of its children
// For instance, given the following hierarchy of states where current state is `F`:
//
// ```
//       A
//      /|\
//     B C D
//    /    /\
//  [F]   G  H
//  /
// I
// ```
//
// This method will return TRUE when checking for states {A, B, F} and FALSE for states {C, D, G, H, I}
func (h *HSM) At(vertex *Vertex) bool {
	h.currentMutex.RLock()
	defer h.currentMutex.RUnlock()

	if vertex.id == h.currentState.id {
		return true
	}

	// check hierarchy
	parent := h.currentState.parent
	for parent != nil {
		if vertex.id == parent.id {
			return true
		}

		parent = parent.parent
	}

	return false
}

// Finished whether HSM is at a final state
func (h *HSM) Finished() bool {
	h.currentMutex.RLock()
	defer h.currentMutex.RUnlock()

	return h.currentState.Final()
}

// Failed whether HSM is at error state
func (h *HSM) Failed() bool {
	h.currentMutex.RLock()
	defer h.currentMutex.RUnlock()

	return h.currentState == h.errorState
}

// Can checks whether the given trigger CAN be signaled, that is, it will produce a transition
func (h *HSM) Can(signal Signal) bool {
	h.currentMutex.RLock()
	defer h.currentMutex.RUnlock()

	for _, t := range h.AvailableSignals() {
		if h.kind(t) == h.kind(signal) {
			return true
		}
	}

	return false
}

// Signal applies the given trigger and fires corresponding transitions if available from current state, if something
// goes wrong this method may "panic"
func (h *HSM) Signal(signal Signal) error {
	h.signalMutex.Lock()
	defer h.signalMutex.Unlock()

	if err := h.tryProgress(); err != nil {
		return err
	}

	return h.apply(signal)
}

// Snapshot returns a serializable snapshot of this HSM.
func (h *HSM) Snapshot() Snapshot {
	h.currentMutex.RLock()
	defer h.currentMutex.RUnlock()

	return Snapshot{
		StateID:        h.currentState.id,
		Final:          h.currentState.Final(),
		SignalsHistory: h.signalsHistory,
		StatesHistory:  h.statesHistory,
	}
}

// AvailableSignals returns a set of events **susceptible** of producing a transition from the outside considering
// HSM`s current state; signals that could be used
func (h *HSM) AvailableSignals() []Signal {
	var signals = make(map[string]Signal)
	var results []Signal

	for _, t := range h.currentState.edges.list() {
		if t.guard == nil || t.guard.method(h.context) {
			signals[h.kind(t.signal)] = t.signal
		}
	}

	parent := h.currentState.parent
	for parent != nil {
		for _, t := range parent.edges.list() {
			if t.guard == nil || t.guard.method(h.context) {
				signals[h.kind(t.signal)] = t.signal
			}
		}

		parent = parent.parent
	}

	for _, evt := range signals {
		results = append(results, evt)
	}

	return results
}

// apply Applies the given signal on this HSM
func (h *HSM) apply(signal Signal) error {
	// do while
	var nextState = h.currentState
	for ok := true; ok; ok = nextState != nil {
		transition := h.getTransition(nextState, signal)

		// If there were no transitions for the given signal for the current
		// state, check if there are any transitions for any of the parent
		// states (if any):
		if transition == nil {
			nextState = nextState.parent
			continue
		}

		// A transition must have a next state defined. If the user has not
		// defined the next state, go to error state:
		if transition.nextStatePtr == nil {
			h.goToErrorState(signal)
			return fmt.Errorf("transition has no next state defined, hsm `%s`", h.name)
		}

		nextState = transition.nextStatePtr

		switch transition.kind {
		case transitionKindInternal:
			return h.doInternalTransition(nextState, transition, signal)
		case transitionKindNormal:
			return h.doNormalTransition(nextState, transition, signal)
		}
	}

	return fmt.Errorf("no transition was found from state `%s` and signal `%s`, hsm `%s`", h.currentState.id, h.kind(signal), h.name)
}

func (h *HSM) doInternalTransition(nextState *Vertex, transition *Transition, signal Signal) error {
	// Run transition effect (if any)
	if transition.effect != nil {
		if err := transition.effect.method(h.context, signal); err != nil {
			h.goToErrorState(signal)
			return err
		}
	}

	// Record in history this successfully applied signal
	h.signalsHistory = append(h.signalsHistory, h.kind(signal))

	// success
	return nil
}

func (h *HSM) doNormalTransition(nextState *Vertex, transition *Transition, signal Signal) error {
	// If the new state is a parent state, enter its entry state (if it has one).
	// Step down through the whole family tree until a state without an entry state is found:
	for nextState.entryState != nil {
		nextState = nextState.entryState
	}

	// Run exit actions only if the current state is left (only if it does not return to itself):
	if h.currentState.onExit != nil {
		if err := h.currentState.onExit.method(h.context, signal); err != nil {
			h.goToErrorState(signal)
			return err
		}
	}

	// Call the current state's parent state exit action if it has one
	// and if new parent state is different than the current state's parent
	if h.currentState.parent != nil &&
		h.currentState.parent.onExit != nil &&
		nextState.parent != h.currentState.parent {
		if err := h.currentState.parent.onExit.method(h.context, signal); err != nil {
			h.goToErrorState(signal)
			return err
		}
	}

	// Run transition effect (if any)
	if transition.effect != nil {
		if err := transition.effect.method(h.context, signal); err != nil {
			h.goToErrorState(signal)
			return err
		}
	}

	// Call the new state's parent state entry action if it has one
	// and if its parent state is different than the current states parent
	// state
	if nextState.parent != nil &&
		nextState.parent.onEntry != nil &&
		nextState.parent != h.currentState.parent {
		if err := nextState.parent.onEntry.method(h.context, signal); err != nil {
			h.goToErrorState(signal)
			return err
		}
	}

	// Call the new state's entry actions if it has any:
	if nextState.onEntry != nil {
		if err := nextState.onEntry.method(h.context, signal); err != nil {
			h.goToErrorState(signal)
			return err
		}
	}

	h.write(nextState, true)
	if h.currentState == h.errorState {
		return fmt.Errorf("error state reached, hsm `%s`", h.name)
	}

	// Record in history this successfully applied signal
	h.signalsHistory = append(h.signalsHistory, h.kind(signal))

	// If next state is a choice pseudo-state then evaluate its branches and transition accordingly
	if h.currentState.kind == vertexKindChoice {
		return h.apply(nil)
	}

	if unconditional := h.getTransition(nextState, nil); unconditional != nil {
		return h.apply(nil)
	}

	// success condition
	return nil
}

func (h *HSM) goToErrorState(signal Signal) {
	h.write(h.errorState, true)

	if s := h.currentState; s != nil && s.onEntry != nil {
		_ = s.onEntry.method(h.context, signal)
	}
}

func (h *HSM) getTransition(from *Vertex, signal Signal) *Transition {
	for _, t := range from.edges.bySignal(signal) {
		if t.guard == nil {
			return t
		}

		if t.guard.method(h.context) {
			return t
		}
	}

	return nil
}

// tryProgress forces hsm to progress if nil signal can be triggered
func (h *HSM) tryProgress() error {
	transitions := h.currentState.edges.bySignal(nil)
	if len(transitions) > 0 {
		return h.apply(nil)
	}

	return nil
}

// kind returns the name of type for the given element
func (h *HSM) kind(i interface{}) string {
	t := reflect.TypeOf(i)
	if t == nil {
		return "nil"
	}

	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}

	return t.Name()
}

// write changes machine state, it is the only point in the code where this occur
func (h *HSM) write(vertex *Vertex, log bool) {
	if log {
		h.statesHistory = append(h.statesHistory, vertex.id)
	}

	h.currentState = vertex
}
