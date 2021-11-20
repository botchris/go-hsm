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

// vertexKind private definition of vertex kind types
type vertexKind int

// Vertex is named element which is an abstraction of a node in a state machine graph. In general, it can
// be the source or destination of any number of transitions.
//
// Subclasses of vertex are:
//
// - `state`
// - `pseudo-state`
type Vertex struct {
	id         string
	kind       vertexKind
	parent     *Vertex
	entryState *Vertex
	onEntry    *Action
	onExit     *Action
	edges      *edgesCollection // transitions indexed by signal type
}

// edgesCollection for handling transitions
type edgesCollection struct {
	edges map[reflect.Type][]*Transition
	count int
}

// newEdgesCollection builds a new empty collection
func newEdgesCollection() *edgesCollection {
	return &edgesCollection{
		edges: make(map[reflect.Type][]*Transition),
	}
}

// add registers a new transition in the collection
func (c *edgesCollection) add(t *Transition) {
	if c.edges == nil {
		c.edges = make(map[reflect.Type][]*Transition)
	}

	if _, ok := c.edges[reflect.TypeOf(t.signal)]; !ok {
		c.edges[reflect.TypeOf(t.signal)] = make([]*Transition, 0)
	}

	if t.guard == nil {
		// APPEND
		c.edges[reflect.TypeOf(t.signal)] = append(c.edges[reflect.TypeOf(t.signal)], t)
	} else {
		// PREPEND
		c.edges[reflect.TypeOf(t.signal)] = append([]*Transition{t}, c.edges[reflect.TypeOf(t.signal)]...)
	}

	c.count++
}

// list returns a plain list of transitions
func (c *edgesCollection) list() []*Transition {
	var transitions []*Transition
	for _, tg := range c.edges {
		transitions = append(transitions, tg...)
	}

	return transitions
}

// bySignal returns a plain list of transitions which signal matches the given one
func (c *edgesCollection) bySignal(s interface{}) []*Transition {
	if transitions, ok := c.edges[reflect.TypeOf(s)]; ok {
		return transitions
	}

	return make([]*Transition, 0)
}

// size returns how many transitions this collections is holding
func (c *edgesCollection) size() int {
	return c.count
}

// ID returns vertex identity
func (n *Vertex) ID() string {
	return n.id
}

// Final indicates whether this vertex is a final state (has no outgoing transitions)
func (n *Vertex) Final() bool {
	return n.edges.size() == 0
}

// NewStart starts building a staring vertex
func NewStart() StartVertexBuilder {
	return &startVertexBuilder{
		edges: newEdgesCollection(),
	}
}

// NewState starts building a new state
func NewState() StateVertexBuilder {
	return &stateVertexBuilder{
		edges: newEdgesCollection(),
	}
}

// NewChoice starts building a new choice pseudo-state
func NewChoice() ChoiceVertexBuilder {
	return &choiceVertexBuilder{
		edges: newEdgesCollection(),
	}
}

// NewErrorState starts building a new error pseudo-state
func NewErrorState() ErrorVertexBuilder {
	return &errorVertexBuilder{}
}

// NewEntryState starts building a new entry pseudo-state
func NewEntryState() EntryVertexBuilder {
	return &entryVertexBuilder{
		edges: newEdgesCollection(),
	}
}

// NewFinalState starts building a new final pseudo-state
func NewFinalState() FinalVertexBuilder {
	return &finalVertexBuilder{}
}
