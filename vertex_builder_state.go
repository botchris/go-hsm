package hsm

// StateVertexBuilder provides builder pattern interface for creating new HSM vertexes
type StateVertexBuilder interface {
	WithID(id string) StateVertexBuilder
	ParentOf(parent *Vertex) StateVertexBuilder
	WithEntryState(entry *Vertex) StateVertexBuilder
	OnEntry(action *Action) StateVertexBuilder
	OnExit(action *Action) StateVertexBuilder
	AddTransitions(transitions ...*Transition) StateVertexBuilder
	Build() *Vertex
}

// stateVertexBuilder private vertex builder
type stateVertexBuilder struct {
	id         string
	parent     *Vertex
	entryState *Vertex
	onEntry    *Action
	onExit     *Action
	edges      *edgesCollection
}

// WithID defines vertex's identity, must be unique within the entire HSM
func (b *stateVertexBuilder) WithID(id string) StateVertexBuilder {
	b.id = id

	return b
}

// ParentOf indicates vertex's parent
func (b *stateVertexBuilder) ParentOf(parent *Vertex) StateVertexBuilder {
	b.parent = parent

	return b
}

// WithEntryState defines an entry point for this vertex, used to manage composed states
func (b *stateVertexBuilder) WithEntryState(entry *Vertex) StateVertexBuilder {
	b.entryState = entry

	return b
}

// OnEntry defines vertex's entry action
func (b *stateVertexBuilder) OnEntry(action *Action) StateVertexBuilder {
	b.onEntry = action

	return b
}

// OnExit defines vertex's exit action
func (b *stateVertexBuilder) OnExit(action *Action) StateVertexBuilder {
	b.onExit = action

	return b
}

// AddTransitions registers the given transitions starting from this vertex
func (b *stateVertexBuilder) AddTransitions(transitions ...*Transition) StateVertexBuilder {
	for _, t := range transitions {
		b.edges.add(t)
	}

	return b
}

// Build returns a vertex instance
func (b *stateVertexBuilder) Build() *Vertex {
	vertex := &Vertex{
		id:         b.id,
		kind:       vertexKindState,
		parent:     b.parent,
		entryState: b.entryState,
		onEntry:    b.onEntry,
		onExit:     b.onExit,
		edges:      b.edges,
	}

	if vertex.entryState != nil {
		vertex.entryState.parent = vertex
	}

	return vertex
}
