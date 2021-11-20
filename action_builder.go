package hsm

// ActionBuilder provides builder pattern interface for creating new action methods
type ActionBuilder interface {
	WithLabel(label string) ActionBuilder
	WithMethod(method ActionFunc) ActionBuilder
	Build() *Action
}

// actionBuilder private action builder
type actionBuilder struct {
	label  string
	method ActionFunc
}

// WithLabel defines action's label
func (b *actionBuilder) WithLabel(label string) ActionBuilder {
	b.label = label

	return b
}

// WithMethod defines action's method
func (b *actionBuilder) WithMethod(method ActionFunc) ActionBuilder {
	b.method = method

	return b
}

// Build returns a new action instance
func (b *actionBuilder) Build() *Action {
	return &Action{
		label:  b.label,
		method: b.method,
	}
}
