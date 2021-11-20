package hsm

const (
	transitionKindNormal = iota
	transitionKindInternal
)

// transitionKind private definition of transition types
type transitionKind int

// Transition represents a transition between two states within a HSM
type Transition struct {
	kind         transitionKind
	signal       Signal
	guard        *Guard
	effect       *Effect
	nextStateID  string
	nextStatePtr *Vertex
}

// NewTransition returns a new transition builder
func NewTransition() TransitionBuilder {
	return &transitionBuilder{}
}

// NewInternalTransition returns a new internal transition builder.
func NewInternalTransition() InternalTransitionBuilder {
	return &internalTransitionBuilder{}
}
