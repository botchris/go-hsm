package hsm

// EntryVertexBuilder builder.
type EntryVertexBuilder[C any] interface {
	WithID(id string) EntryVertexBuilder[C]
	ParentOf(parent *Vertex[C]) EntryVertexBuilder[C]
	OnEntry(action *Action[C]) EntryVertexBuilder[C]
	OnExit(action *Action[C]) EntryVertexBuilder[C]
	WithTransitions(transitions ...*Transition[C]) EntryVertexBuilder[C]
	Build() *Vertex[C]
}

type entryVertexBuilder[C any] struct {
	id      string
	parent  *Vertex[C]
	onEntry *Action[C]
	onExit  *Action[C]
	edges   *edgesCollection[C]
}

// WithID defines vertex's identity, must be unique within the entire HSM.
func (b *entryVertexBuilder[C]) WithID(id string) EntryVertexBuilder[C] {
	b.id = id

	return b
}

// ParentOf indicates vertex's parent.
func (b *entryVertexBuilder[C]) ParentOf(parent *Vertex[C]) EntryVertexBuilder[C] {
	b.parent = parent

	return b
}

// OnEntry defines vertex's entry action.
func (b *entryVertexBuilder[C]) OnEntry(action *Action[C]) EntryVertexBuilder[C] {
	b.onEntry = action

	return b
}

// OnExit defines vertex's exit action.
func (b *entryVertexBuilder[C]) OnExit(action *Action[C]) EntryVertexBuilder[C] {
	b.onExit = action

	return b
}

// WithTransitions registers the given transitions starting from this vertex.
func (b *entryVertexBuilder[C]) WithTransitions(transitions ...*Transition[C]) EntryVertexBuilder[C] {
	for _, t := range transitions {
		b.edges.add(t)
	}

	return b
}

// Build returns a vertex instance.
func (b *entryVertexBuilder[C]) Build() *Vertex[C] {
	vertex := &Vertex[C]{
		id:      b.id,
		kind:    vertexKindEntry,
		parent:  b.parent,
		onEntry: b.onEntry,
		onExit:  b.onExit,
		edges:   b.edges,
	}

	return vertex
}
