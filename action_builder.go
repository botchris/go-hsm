package hsm

// ActionBuilder provides builder pattern interface for creating new action methods.
type ActionBuilder[C any] interface {
	WithLabel(label string) ActionBuilder[C]
	WithMethod(method ActionFunc[C]) ActionBuilder[C]
	Build() *Action[C]
}

// actionBuilder private action builder.
type actionBuilder[C any] struct {
	label  string
	method ActionFunc[C]
}

// WithLabel defines action's label.
func (b *actionBuilder[C]) WithLabel(label string) ActionBuilder[C] {
	b.label = label

	return b
}

// WithMethod defines action's method.
func (b *actionBuilder[C]) WithMethod(method ActionFunc[C]) ActionBuilder[C] {
	b.method = method

	return b
}

// Build returns a new action instance.
func (b *actionBuilder[C]) Build() *Action[C] {
	return &Action[C]{
		label:  b.label,
		method: b.method,
	}
}
