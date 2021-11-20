package hsm

// TransitionBuilder provides builder pattern interface for creating new HSM regular transitions
type TransitionBuilder interface {
	When(signal Signal) TransitionBuilder
	GuardedBy(guard *Guard) TransitionBuilder
	ApplyEffect(effect *Effect) TransitionBuilder
	GoTo(stateID string) TransitionBuilder
	Build() *Transition
}

// transitionBuilder private transition builder
type transitionBuilder struct {
	signal      Signal
	guard       *Guard
	effect      *Effect
	nextStateID string
}

// When indicates which signal activates this transition
func (b *transitionBuilder) When(signal Signal) TransitionBuilder {
	b.signal = signal

	return b
}

// GuardedBy indicates this transition is guarded by the given guard
func (b *transitionBuilder) GuardedBy(guard *Guard) TransitionBuilder {
	b.guard = guard

	return b
}

// ApplyEffect registers an effect for this transition
func (b *transitionBuilder) ApplyEffect(effect *Effect) TransitionBuilder {
	b.effect = effect

	return b
}

// GoTo defines the next state where to transition to
func (b *transitionBuilder) GoTo(stateID string) TransitionBuilder {
	b.nextStateID = stateID

	return b
}

// Build finalizes the building process of this transition
func (b *transitionBuilder) Build() *Transition {
	var signal = b.signal
	var guard = b.guard
	var effect = b.effect

	transition := &Transition{
		kind:        transitionKindNormal,
		signal:      signal,
		guard:       guard,
		effect:      effect,
		nextStateID: b.nextStateID,
	}

	return transition
}
