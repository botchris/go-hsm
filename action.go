package hsm

// ActionFunc public definition of an action method
type ActionFunc func(ctx interface{}, signal Signal) error

// Action definition of entry/exit logic
type Action struct {
	label  string
	method ActionFunc
}

// NewAction starts building a new vertex action instance
func NewAction() ActionBuilder {
	return &actionBuilder{}
}
