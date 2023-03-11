//nolint:dupl
package hsm

// StartVertexBuilder builder.
type StartVertexBuilder[C any] interface {
	WithID(id string) StartVertexBuilder[C]
	OnExit(action *Action[C]) StartVertexBuilder[C]
	AddTransitions(transitions ...*Transition[C]) StartVertexBuilder[C]
	Build() *Vertex[C]
}

type startVertexBuilder[C any] struct {
	id     string
	onExit *Action[C]
	edges  *edgesCollection[C]
}

// WithID defines vertex's identity, must be unique within the entire HSM.
func (b *startVertexBuilder[C]) WithID(id string) StartVertexBuilder[C] {
	b.id = id

	return b
}

// OnExit defines vertex's exit action.
func (b *startVertexBuilder[C]) OnExit(action *Action[C]) StartVertexBuilder[C] {
	b.onExit = action

	return b
}

// AddTransitions registers the given transitions starting from this vertex.
func (b *startVertexBuilder[C]) AddTransitions(transitions ...*Transition[C]) StartVertexBuilder[C] {
	for _, t := range transitions {
		b.edges.add(t)
	}

	return b
}

// Build returns a vertex instance.
func (b *startVertexBuilder[C]) Build() *Vertex[C] {
	vertex := &Vertex[C]{
		id:     b.id,
		kind:   vertexKindStart,
		onExit: b.onExit,
		edges:  b.edges,
	}

	return vertex
}
