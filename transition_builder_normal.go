package hsm

// TransitionBuilder provides builder pattern interface for creating new HSM regular transitions.
type TransitionBuilder[C any] interface {
	When(signal Signal) TransitionBuilder[C]
	GuardedBy(guard *Guard[C]) TransitionBuilder[C]
	ApplyEffect(effect *Effect[C]) TransitionBuilder[C]
	GoTo(stateID string) TransitionBuilder[C]
	Build() *Transition[C]
}

// transitionBuilder private transition builder.
type transitionBuilder[C any] struct {
	signal      Signal
	guard       *Guard[C]
	effect      *Effect[C]
	nextStateID string
}

// When indicates which signal activates this transition.
func (b *transitionBuilder[C]) When(signal Signal) TransitionBuilder[C] {
	b.signal = signal

	return b
}

// GuardedBy indicates this transition is guarded by the given guard.
func (b *transitionBuilder[C]) GuardedBy(guard *Guard[C]) TransitionBuilder[C] {
	b.guard = guard

	return b
}

// ApplyEffect registers an effect for this transition.
func (b *transitionBuilder[C]) ApplyEffect(effect *Effect[C]) TransitionBuilder[C] {
	b.effect = effect

	return b
}

// GoTo defines the next state where to transition to.
func (b *transitionBuilder[C]) GoTo(stateID string) TransitionBuilder[C] {
	b.nextStateID = stateID

	return b
}

// Build finalizes the building process of this transition.
func (b *transitionBuilder[C]) Build() *Transition[C] {
	var (
		signal = b.signal
		guard  = b.guard
		effect = b.effect
	)

	transition := &Transition[C]{
		kind:        transitionKindNormal,
		signal:      signal,
		guard:       guard,
		effect:      effect,
		nextStateID: b.nextStateID,
	}

	return transition
}
