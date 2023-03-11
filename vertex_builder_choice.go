//nolint:dupl
package hsm

// ChoiceVertexBuilder builder.
type ChoiceVertexBuilder[C any] interface {
	WithID(id string) ChoiceVertexBuilder[C]
	ParentOf(parent *Vertex[C]) ChoiceVertexBuilder[C]
	AddTransitions(transitions ...*Transition[C]) ChoiceVertexBuilder[C]
	Build() *Vertex[C]
}

type choiceVertexBuilder[C any] struct {
	id     string
	parent *Vertex[C]
	edges  *edgesCollection[C]
}

// WithID defines vertex's identity, must be unique within the entire HSM.
func (b *choiceVertexBuilder[C]) WithID(id string) ChoiceVertexBuilder[C] {
	b.id = id

	return b
}

// ParentOf indicates vertex's parent.
func (b *choiceVertexBuilder[C]) ParentOf(parent *Vertex[C]) ChoiceVertexBuilder[C] {
	b.parent = parent

	return b
}

// AddTransitions registers the given transitions starting from this vertex.
func (b *choiceVertexBuilder[C]) AddTransitions(transitions ...*Transition[C]) ChoiceVertexBuilder[C] {
	for _, t := range transitions {
		b.edges.add(t)
	}

	return b
}

// Build returns a vertex instance.
func (b *choiceVertexBuilder[C]) Build() *Vertex[C] {
	vertex := &Vertex[C]{
		id:     b.id,
		kind:   vertexKindChoice,
		parent: b.parent,
		edges:  b.edges,
	}

	return vertex
}
