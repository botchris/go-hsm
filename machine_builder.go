package hsm

import (
	"fmt"
)

// Builder defines a builder pattern for creating new FSMs.
type Builder[C any] struct {
	hsm   *HSM[C]
	start *Vertex[C]
}

// NewBuilder returns a new builder instance.
func NewBuilder[C any]() *Builder[C] {
	builder := &Builder[C]{
		hsm: &HSM[C]{
			signalsHistory: make([]string, 0),
			statesHistory:  make([]string, 0),
			states:         make(map[string]*Vertex[C]),
		},
	}

	return builder
}

// WithName defines a name for this HSM instance, used for visual representations.
func (b *Builder[C]) WithName(name string) *Builder[C] {
	b.hsm.name = name

	return b
}

// WithErrorState registers an error state.
func (b *Builder[C]) WithErrorState(state *Vertex[C]) *Builder[C] {
	b.hsm.errorState = state

	return b
}

// WithContext sets HSM`s context.
func (b *Builder[C]) WithContext(ctx C) *Builder[C] {
	b.hsm.context = ctx

	return b
}

// StartingAt sets HSM`s starting state.
func (b *Builder[C]) StartingAt(state *Vertex[C]) *Builder[C] {
	b.start = state

	return b
}

// AddState registers a single state into this machine.
func (b *Builder[C]) AddState(state *Vertex[C]) *Builder[C] {
	return b.AddStates(state)
}

// AddStates registers multiple states at once.
func (b *Builder[C]) AddStates(states ...*Vertex[C]) *Builder[C] {
	for _, s := range states {
		b.hsm.states[s.id] = s
	}

	return b
}

// Restore builds a new machine instance and restores from the given snapshot.
// no guards are checked nor entry/exit logic will be executed.
func (b *Builder[C]) Restore(snapshot Snapshot) (*HSM[C], error) {
	machine, err := b.Build()
	if err != nil {
		return nil, err
	}

	if _, ok := machine.states[snapshot.StateID]; !ok {
		return nil, fmt.Errorf("starting state `%s` does not exists", snapshot.StateID)
	}

	machine.signalsHistory = snapshot.SignalsHistory
	machine.statesHistory = snapshot.StatesHistory
	machine.write(machine.states[snapshot.StateID], false)

	// force hsm to progress if nil signal can be triggered
	if err := machine.tryProgress(); err != nil {
		return nil, err
	}

	return machine, nil
}

// Build builds the HSM.
func (b *Builder[C]) Build() (*HSM[C], error) {
	if b.hsm.name == "" {
		return nil, fmt.Errorf("no name was provided fot his HSM")
	}

	if b.start == nil {
		return nil, fmt.Errorf("no starting state was provided")
	}

	if b.hsm.errorState == nil {
		return nil, fmt.Errorf("no error state was defined")
	}

	if err := b.validateVertex(b.hsm.errorState); err != nil {
		return nil, err
	}

	// compute state transitions
	for _, s := range b.hsm.states {
		var transitions []*Transition[C]
		transitions = append(transitions, s.edges.list()...)

		if s.entryState != nil && s.entryState.edges.size() > 0 {
			transitions = append(transitions, s.entryState.edges.list()...)
		}

		for _, t := range transitions {
			if t.kind == transitionKindInternal {
				t.nextStateID = s.id
			}

			if _, ok := b.hsm.states[t.nextStateID]; !ok {
				return nil, fmt.Errorf("state `%s` not found for transition", t.nextStateID)
			}

			t.nextStatePtr = b.hsm.states[t.nextStateID]
		}

		if err := b.validateVertex(s); err != nil {
			return nil, err
		}
	}

	b.hsm.write(b.start, true)

	return b.hsm, nil
}

// validate ensures the integrity of the given vertex and all its parts
//
//nolint:gocyclo
func (b *Builder[C]) validateVertex(v *Vertex[C]) error {
	if v.id == "" {
		return fmt.Errorf("invalid state identity, cannot be empty")
	}

	if v.parent != nil {
		if _, ok := b.hsm.states[v.parent.id]; !ok {
			return fmt.Errorf("invalid state parent, parent state `%s` was not found in this machine", v.parent.id)
		}
	}

	if v.onEntry != nil {
		if v.onEntry.label == "" {
			return fmt.Errorf("invalid state entry logic, no action label was provided")
		}

		if v.onEntry.method == nil {
			return fmt.Errorf("invalid state entry logic, no method was defined")
		}
	}

	if v.onExit != nil {
		if v.onExit.label == "" {
			return fmt.Errorf("invalid state exit logic, no label was provided")
		}

		if v.onExit.method == nil {
			return fmt.Errorf("invalid state exit logic, no method was defined")
		}
	}

	for _, t := range v.edges.list() {
		if t.nextStateID == "" {
			return fmt.Errorf("invalid transition, no next state was provided")
		}

		if _, ok := b.hsm.states[t.nextStateID]; !ok {
			return fmt.Errorf("invalid transition, no next state `%s` does not exists", t.nextStateID)
		}

		if v.kind == vertexKindFinal {
			return fmt.Errorf("invalid transition, final states cannot have outgoing transitions")
		}

		if v.kind == vertexKindError {
			return fmt.Errorf("invalid transition, error states cannot have outgoing transitions")
		}

		if t.guard != nil && t.guard.label == "" {
			return fmt.Errorf("invalid transition, nameless guard provided")
		}

		if t.effect != nil && t.effect.label == "" {
			return fmt.Errorf("invalid transition, effects must provide a valid human-readable representation")
		}
	}

	return nil
}
