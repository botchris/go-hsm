package hsm

// ErrorVertexBuilder builder
type ErrorVertexBuilder interface {
	WithID(id string) ErrorVertexBuilder
	OnEntry(action *Action) ErrorVertexBuilder
	Build() *Vertex
}

type errorVertexBuilder struct {
	id      string
	onEntry *Action
}

// WithID defines vertex's identity, must be unique within the entire HSM
func (b *errorVertexBuilder) WithID(id string) ErrorVertexBuilder {
	b.id = id

	return b
}

// OnEntry defines vertex's entry action
func (b *errorVertexBuilder) OnEntry(action *Action) ErrorVertexBuilder {
	b.onEntry = action

	return b
}

// Build returns a vertex instance
func (b *errorVertexBuilder) Build() *Vertex {
	vertex := &Vertex{
		id:      b.id,
		kind:    vertexKindError,
		onEntry: b.onEntry,
		edges:   newEdgesCollection(),
	}

	return vertex
}
