package hsm

// InternalTransitionBuilder provides builder pattern interface for creating new HSM internal transitions.
type InternalTransitionBuilder[C any] interface {
	When(signal Signal) InternalTransitionBuilder[C]
	GuardedBy(guard *Guard[C]) InternalTransitionBuilder[C]
	ApplyEffect(effect *Effect[C]) InternalTransitionBuilder[C]
	Build() *Transition[C]
}

// internalTransitionBuilder private transition builder.
type internalTransitionBuilder[C any] struct {
	signal      Signal
	guard       *Guard[C]
	effect      *Effect[C]
	nextStateID string
}

// When indicates which signal activates this transition.
func (b *internalTransitionBuilder[C]) When(signal Signal) InternalTransitionBuilder[C] {
	b.signal = signal

	return b
}

// GuardedBy indicates this transition is guarded by the given guard.
func (b *internalTransitionBuilder[C]) GuardedBy(guard *Guard[C]) InternalTransitionBuilder[C] {
	b.guard = guard

	return b
}

// ApplyEffect registers an effect for this transition.
func (b *internalTransitionBuilder[C]) ApplyEffect(effect *Effect[C]) InternalTransitionBuilder[C] {
	b.effect = effect

	return b
}

// Build finalizes the building process of this transition.
func (b *internalTransitionBuilder[C]) Build() *Transition[C] {
	var (
		signal = b.signal
		guard  = b.guard
		effect = b.effect
	)

	transition := &Transition[C]{
		kind:        transitionKindInternal,
		signal:      signal,
		guard:       guard,
		effect:      effect,
		nextStateID: b.nextStateID,
	}

	return transition
}
