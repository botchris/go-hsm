package hsm

// EffectBuilder provides builder pattern interface for creating new HSM transition effects
type EffectBuilder interface {
	WithLabel(label string) EffectBuilder
	WithMethod(method ActionFunc) EffectBuilder
	Build() *Effect
}

// effectBuilder private effect builder
type effectBuilder struct {
	label  string
	method ActionFunc
}

// WithLabel defines effect's label
func (b *effectBuilder) WithLabel(label string) EffectBuilder {
	b.label = label

	return b
}

// WithMethod defines effect's action
func (b *effectBuilder) WithMethod(method ActionFunc) EffectBuilder {
	b.method = method

	return b
}

// Build builds and returns the effect
func (b *effectBuilder) Build() *Effect {
	return &Effect{
		label:  b.label,
		method: b.method,
	}
}
