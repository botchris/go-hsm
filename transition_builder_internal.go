package hsm

// InternalTransitionBuilder provides builder pattern interface for creating new HSM internal transitions
type InternalTransitionBuilder interface {
	When(signal Signal) InternalTransitionBuilder
	GuardedBy(guard *Guard) InternalTransitionBuilder
	ApplyEffect(effect *Effect) InternalTransitionBuilder
	Build() *Transition
}

// internalTransitionBuilder private transition builder
type internalTransitionBuilder struct {
	signal      Signal
	guard       *Guard
	effect      *Effect
	nextStateID string
}

// When indicates which signal activates this transition
func (b *internalTransitionBuilder) When(signal Signal) InternalTransitionBuilder {
	b.signal = signal

	return b
}

// GuardedBy indicates this transition is guarded by the given guard
func (b *internalTransitionBuilder) GuardedBy(guard *Guard) InternalTransitionBuilder {
	b.guard = guard

	return b
}

// ApplyEffect registers an effect for this transition
func (b *internalTransitionBuilder) ApplyEffect(effect *Effect) InternalTransitionBuilder {
	b.effect = effect

	return b
}

// Build finalizes the building process of this transition
func (b *internalTransitionBuilder) Build() *Transition {
	var signal = b.signal
	var guard = b.guard
	var effect = b.effect

	transition := &Transition{
		kind:        transitionKindInternal,
		signal:      signal,
		guard:       guard,
		effect:      effect,
		nextStateID: b.nextStateID,
	}

	return transition
}
