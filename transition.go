package hsm

const (
	transitionKindNormal = iota
	transitionKindInternal
)

// transitionKind private definition of transition types.
type transitionKind int

// Transition represents a transition between two states within a HSM.
type Transition[C any] struct {
	kind         transitionKind
	signal       Signal
	guard        *Guard[C]
	effect       *Effect[C]
	nextStateID  string
	nextStatePtr *Vertex[C]
}

// NewTransition returns a new transition builder.
func NewTransition[C any]() TransitionBuilder[C] {
	return &transitionBuilder[C]{}
}

// NewInternalTransition returns a new internal transition builder.
func NewInternalTransition[C any]() InternalTransitionBuilder[C] {
	return &internalTransitionBuilder[C]{}
}
