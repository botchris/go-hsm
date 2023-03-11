package hsm

// EffectBuilder provides builder pattern interface for creating new HSM transition effects.
type EffectBuilder[C any] interface {
	WithLabel(label string) EffectBuilder[C]
	WithMethod(method ActionFunc[C]) EffectBuilder[C]
	Build() *Effect[C]
}

// effectBuilder private effect builder.
type effectBuilder[C any] struct {
	label  string
	method ActionFunc[C]
}

// WithLabel defines effect's label.
func (b *effectBuilder[C]) WithLabel(label string) EffectBuilder[C] {
	b.label = label

	return b
}

// WithMethod defines effect's action.
func (b *effectBuilder[C]) WithMethod(method ActionFunc[C]) EffectBuilder[C] {
	b.method = method

	return b
}

// Build builds and returns the effect.
func (b *effectBuilder[C]) Build() *Effect[C] {
	return &Effect[C]{
		label:  b.label,
		method: b.method,
	}
}
