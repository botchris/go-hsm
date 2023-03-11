package hsm

// ErrorVertexBuilder builder.
type ErrorVertexBuilder[C any] interface {
	WithID(id string) ErrorVertexBuilder[C]
	OnEntry(action *Action[C]) ErrorVertexBuilder[C]
	Build() *Vertex[C]
}

type errorVertexBuilder[C any] struct {
	id      string
	onEntry *Action[C]
}

// WithID defines vertex's identity, must be unique within the entire HSM.
func (b *errorVertexBuilder[C]) WithID(id string) ErrorVertexBuilder[C] {
	b.id = id

	return b
}

// OnEntry defines vertex's entry action.
func (b *errorVertexBuilder[C]) OnEntry(action *Action[C]) ErrorVertexBuilder[C] {
	b.onEntry = action

	return b
}

// Build returns a vertex instance.
func (b *errorVertexBuilder[C]) Build() *Vertex[C] {
	vertex := &Vertex[C]{
		id:      b.id,
		kind:    vertexKindError,
		onEntry: b.onEntry,
		edges:   newEdgesCollection[C](),
	}

	return vertex
}
