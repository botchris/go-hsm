package hsm

// StateVertexBuilder provides builder pattern interface for creating new HSM vertexes.
type StateVertexBuilder[C any] interface {
	WithID(id string) StateVertexBuilder[C]
	ParentOf(parent *Vertex[C]) StateVertexBuilder[C]
	WithEntryState(entry *Vertex[C]) StateVertexBuilder[C]
	OnEntry(action *Action[C]) StateVertexBuilder[C]
	OnExit(action *Action[C]) StateVertexBuilder[C]
	WithTransitions(transitions ...*Transition[C]) StateVertexBuilder[C]
	Build() *Vertex[C]
}

// stateVertexBuilder private vertex builder.
type stateVertexBuilder[C any] struct {
	id         string
	parent     *Vertex[C]
	entryState *Vertex[C]
	onEntry    *Action[C]
	onExit     *Action[C]
	edges      *edgesCollection[C]
}

// WithID defines vertex's identity, must be unique within the entire HSM.
func (b *stateVertexBuilder[C]) WithID(id string) StateVertexBuilder[C] {
	b.id = id

	return b
}

// ParentOf indicates vertex's parent.
func (b *stateVertexBuilder[C]) ParentOf(parent *Vertex[C]) StateVertexBuilder[C] {
	b.parent = parent

	return b
}

// WithEntryState defines an entry point for this vertex, used to manage composed states.
func (b *stateVertexBuilder[C]) WithEntryState(entry *Vertex[C]) StateVertexBuilder[C] {
	b.entryState = entry

	return b
}

// OnEntry defines vertex's entry action.
func (b *stateVertexBuilder[C]) OnEntry(action *Action[C]) StateVertexBuilder[C] {
	b.onEntry = action

	return b
}

// OnExit defines vertex's exit action.
func (b *stateVertexBuilder[C]) OnExit(action *Action[C]) StateVertexBuilder[C] {
	b.onExit = action

	return b
}

// WithTransitions registers the given transitions starting from this vertex.
func (b *stateVertexBuilder[C]) WithTransitions(transitions ...*Transition[C]) StateVertexBuilder[C] {
	for _, t := range transitions {
		b.edges.add(t)
	}

	return b
}

// Build returns a vertex instance.
func (b *stateVertexBuilder[C]) Build() *Vertex[C] {
	vertex := &Vertex[C]{
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
