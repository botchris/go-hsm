package hsm

// ChoiceVertexBuilder builder
type ChoiceVertexBuilder interface {
	WithID(id string) ChoiceVertexBuilder
	ParentOf(parent *Vertex) ChoiceVertexBuilder
	AddTransitions(transitions ...*Transition) ChoiceVertexBuilder
	Build() *Vertex
}

type choiceVertexBuilder struct {
	id     string
	parent *Vertex
	edges  *edgesCollection
}

// WithID defines vertex's identity, must be unique within the entire HSM
func (b *choiceVertexBuilder) WithID(id string) ChoiceVertexBuilder {
	b.id = id

	return b
}

// ParentOf indicates vertex's parent
func (b *choiceVertexBuilder) ParentOf(parent *Vertex) ChoiceVertexBuilder {
	b.parent = parent

	return b
}

// AddTransitions registers the given transitions starting from this vertex
func (b *choiceVertexBuilder) AddTransitions(transitions ...*Transition) ChoiceVertexBuilder {
	for _, t := range transitions {
		b.edges.add(t)
	}

	return b
}

// Build returns a vertex instance
func (b *choiceVertexBuilder) Build() *Vertex {
	vertex := &Vertex{
		id:     b.id,
		kind:   vertexKindChoice,
		parent: b.parent,
		edges:  b.edges,
	}

	return vertex
}
