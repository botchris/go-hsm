package hsm

// StartVertexBuilder builder
type StartVertexBuilder interface {
	WithID(id string) StartVertexBuilder
	OnExit(action *Action) StartVertexBuilder
	AddTransitions(transitions ...*Transition) StartVertexBuilder
	Build() *Vertex
}

type startVertexBuilder struct {
	id     string
	onExit *Action
	edges  *edgesCollection
}

// WithID defines vertex's identity, must be unique within the entire HSM
func (b *startVertexBuilder) WithID(id string) StartVertexBuilder {
	b.id = id

	return b
}

// OnExit defines vertex's exit action
func (b *startVertexBuilder) OnExit(action *Action) StartVertexBuilder {
	b.onExit = action

	return b
}

// AddTransitions registers the given transitions starting from this vertex
func (b *startVertexBuilder) AddTransitions(transitions ...*Transition) StartVertexBuilder {
	for _, t := range transitions {
		b.edges.add(t)
	}

	return b
}

// Build returns a vertex instance
func (b *startVertexBuilder) Build() *Vertex {
	vertex := &Vertex{
		id:     b.id,
		kind:   vertexKindStart,
		onExit: b.onExit,
		edges:  b.edges,
	}

	return vertex
}
