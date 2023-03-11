package hsm

import "reflect"

const (
	vertexKindState = iota
	vertexKindChoice
	vertexKindEntry
	vertexKindStart
	vertexKindFinal
	vertexKindError
)

// vertexKind private definition of vertex kind types.
type vertexKind int

// Vertex is an abstraction of a node in a state machine graph. In general, it can
// be the source or destination of any number of transitions.
//
// Subclasses of vertex are:
//
// - `state`
// - `pseudo-state`.
type Vertex[C any] struct {
	id         string
	kind       vertexKind
	parent     *Vertex[C]
	entryState *Vertex[C]
	onEntry    *Action[C]
	onExit     *Action[C]
	edges      *edgesCollection[C] // transitions indexed by signal type
}

// edgesCollection for handling transitions.
type edgesCollection[C any] struct {
	edges map[reflect.Type][]*Transition[C]
	count int
}

// newEdgesCollection builds a new empty collection.
func newEdgesCollection[C any]() *edgesCollection[C] {
	return &edgesCollection[C]{
		edges: make(map[reflect.Type][]*Transition[C]),
	}
}

// add registers a new transition in the collection.
func (c *edgesCollection[C]) add(t *Transition[C]) {
	if c.edges == nil {
		c.edges = make(map[reflect.Type][]*Transition[C])
	}

	if _, ok := c.edges[reflect.TypeOf(t.signal)]; !ok {
		c.edges[reflect.TypeOf(t.signal)] = make([]*Transition[C], 0)
	}

	if t.guard == nil {
		// APPEND
		c.edges[reflect.TypeOf(t.signal)] = append(c.edges[reflect.TypeOf(t.signal)], t)
	} else {
		// PREPEND
		c.edges[reflect.TypeOf(t.signal)] = append([]*Transition[C]{t}, c.edges[reflect.TypeOf(t.signal)]...)
	}

	c.count++
}

// list returns a plain list of transitions.
func (c *edgesCollection[C]) list() []*Transition[C] {
	var transitions []*Transition[C]
	for _, tg := range c.edges {
		transitions = append(transitions, tg...)
	}

	return transitions
}

// bySignal returns a plain list of transitions which signal matches the given one.
func (c *edgesCollection[C]) bySignal(s interface{}) []*Transition[C] {
	if transitions, ok := c.edges[reflect.TypeOf(s)]; ok {
		return transitions
	}

	return make([]*Transition[C], 0)
}

// size returns how many transitions this collection is holding.
func (c *edgesCollection[C]) size() int {
	return c.count
}

// ID returns vertex identity.
func (n *Vertex[C]) ID() string {
	return n.id
}

// Final indicates whether this vertex is a final state (has no outgoing transitions).
func (n *Vertex[C]) Final() bool {
	return n.edges.size() == 0
}

// NewStart starts building a staring vertex.
func NewStart[C any]() StartVertexBuilder[C] {
	return &startVertexBuilder[C]{
		edges: newEdgesCollection[C](),
	}
}

// NewState starts building a new state.
func NewState[C any]() StateVertexBuilder[C] {
	return &stateVertexBuilder[C]{
		edges: newEdgesCollection[C](),
	}
}

// NewChoice starts building a new choice pseudo-state.
func NewChoice[C any]() ChoiceVertexBuilder[C] {
	return &choiceVertexBuilder[C]{
		edges: newEdgesCollection[C](),
	}
}

// NewErrorState starts building a new error pseudo-state.
func NewErrorState[C any]() ErrorVertexBuilder[C] {
	return &errorVertexBuilder[C]{}
}

// NewEntryState starts building a new entry pseudo-state.
func NewEntryState[C any]() EntryVertexBuilder[C] {
	return &entryVertexBuilder[C]{
		edges: newEdgesCollection[C](),
	}
}

// NewFinalState starts building a new final pseudo-state.
func NewFinalState[C any]() FinalVertexBuilder[C] {
	return &finalVertexBuilder[C]{}
}
