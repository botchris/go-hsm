package hsm

// GuardBuilder provides builder pattern interface for creating new guard conditions
type GuardBuilder interface {
	WithLabel(label string) GuardBuilder
	WithMethod(method GuardFunc) GuardBuilder
	Build() *Guard
}

// guardBuilder private guard builder
type guardBuilder struct {
	label  string
	method GuardFunc
}

// WithLabel defines guard's label
func (b *guardBuilder) WithLabel(label string) GuardBuilder {
	b.label = label

	return b
}

// WithMethod defines guard's method
func (b *guardBuilder) WithMethod(method GuardFunc) GuardBuilder {
	b.method = method

	return b
}

// Build finalizes the building process of this guard
func (b *guardBuilder) Build() *Guard {
	return &Guard{
		label:  b.label,
		method: b.method,
	}
}
