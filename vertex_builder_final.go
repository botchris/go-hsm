package hsm

// FinalVertexBuilder builder
type FinalVertexBuilder interface {
	WithID(id string) FinalVertexBuilder
	OnEntry(action *Action) FinalVertexBuilder
	Build() *Vertex
}

type finalVertexBuilder struct {
	id      string
	onEntry *Action
}

// WithID defines vertex's identity, must be unique within the entire HSM
func (b *finalVertexBuilder) WithID(id string) FinalVertexBuilder {
	b.id = id

	return b
}

// OnEntry defines vertex's entry action
func (b *finalVertexBuilder) OnEntry(action *Action) FinalVertexBuilder {
	b.onEntry = action

	return b
}

// Build returns a vertex instance
func (b *finalVertexBuilder) Build() *Vertex {
	vertex := &Vertex{
		id:      b.id,
		kind:    vertexKindFinal,
		onEntry: b.onEntry,
		edges:   newEdgesCollection(),
	}

	return vertex
}
