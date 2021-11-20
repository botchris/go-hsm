package hsm

// GuardFunc public definition of guard functions
type GuardFunc func(ctx interface{}) bool

// Guard definition of transition guard, checks whether a transition can be performed or not based on given context;
// they MUST be side-effect free, at least none that would alter evaluation of other guards having the same trigger.
type Guard struct {
	label  string
	method GuardFunc
}

// NewGuard starts building a new guard condition
func NewGuard() GuardBuilder {
	return &guardBuilder{}
}
