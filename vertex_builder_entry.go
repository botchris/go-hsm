package hsm

// EntryVertexBuilder builder
type EntryVertexBuilder interface {
	WithID(id string) EntryVertexBuilder
	ParentOf(parent *Vertex) EntryVertexBuilder
	OnEntry(action *Action) EntryVertexBuilder
	OnExit(action *Action) EntryVertexBuilder
	AddTransitions(transitions ...*Transition) EntryVertexBuilder
	Build() *Vertex
}

type entryVertexBuilder struct {
	id      string
	parent  *Vertex
	onEntry *Action
	onExit  *Action
	edges   *edgesCollection
}

// WithID defines vertex's identity, must be unique within the entire HSM
func (b *entryVertexBuilder) WithID(id string) EntryVertexBuilder {
	b.id = id

	return b
}

// ParentOf indicates vertex's parent
func (b *entryVertexBuilder) ParentOf(parent *Vertex) EntryVertexBuilder {
	b.parent = parent

	return b
}

// OnEntry defines vertex's entry action
func (b *entryVertexBuilder) OnEntry(action *Action) EntryVertexBuilder {
	b.onEntry = action

	return b
}

// OnExit defines vertex's exit action
func (b *entryVertexBuilder) OnExit(action *Action) EntryVertexBuilder {
	b.onExit = action

	return b
}

// AddTransitions registers the given transitions starting from this vertex
func (b *entryVertexBuilder) AddTransitions(transitions ...*Transition) EntryVertexBuilder {
	for _, t := range transitions {
		b.edges.add(t)
	}

	return b
}

// Build returns a vertex instance
func (b *entryVertexBuilder) Build() *Vertex {
	vertex := &Vertex{
		id:      b.id,
		kind:    vertexKindEntry,
		parent:  b.parent,
		onEntry: b.onEntry,
		onExit:  b.onExit,
		edges:   b.edges,
	}

	return vertex
}
