package hsm

// GuardBuilder provides builder pattern interface for creating new guard conditions.
type GuardBuilder[C any] interface {
	WithLabel(label string) GuardBuilder[C]
	WithMethod(method GuardFunc[C]) GuardBuilder[C]
	Build() *Guard[C]
}

// guardBuilder private guard builder.
type guardBuilder[C any] struct {
	label  string
	method GuardFunc[C]
}

// WithLabel defines guard's label.
func (b *guardBuilder[C]) WithLabel(label string) GuardBuilder[C] {
	b.label = label

	return b
}

// WithMethod defines guard's method.
func (b *guardBuilder[C]) WithMethod(method GuardFunc[C]) GuardBuilder[C] {
	b.method = method

	return b
}

// Build finalizes the building process of this guard.
func (b *guardBuilder[C]) Build() *Guard[C] {
	return &Guard[C]{
		label:  b.label,
		method: b.method,
	}
}
