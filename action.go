package hsm

// ActionFunc public definition of an action method.
type ActionFunc[C any] func(ctx C, signal Signal) error

// Action definition of entry/exit logic.
type Action[C any] struct {
	label  string
	method ActionFunc[C]
}

// String returns a string representation of the action.
func (a *Action[C]) String() string {
	if a.label != "" {
		return a.label
	}

	return fnSignatureString(a.method)
}

// NewAction starts building a new vertex action instance.
func NewAction[C any]() ActionBuilder[C] {
	return &actionBuilder[C]{}
}
