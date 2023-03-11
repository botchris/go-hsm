package hsm

// FinalVertexBuilder builder.
type FinalVertexBuilder[C any] interface {
	WithID(id string) FinalVertexBuilder[C]
	OnEntry(action *Action[C]) FinalVertexBuilder[C]
	Build() *Vertex[C]
}

type finalVertexBuilder[C any] struct {
	id      string
	onEntry *Action[C]
}

// WithID defines vertex's identity, must be unique within the entire HSM.
func (b *finalVertexBuilder[C]) WithID(id string) FinalVertexBuilder[C] {
	b.id = id

	return b
}

// OnEntry defines vertex's entry action.
func (b *finalVertexBuilder[C]) OnEntry(action *Action[C]) FinalVertexBuilder[C] {
	b.onEntry = action

	return b
}

// Build returns a vertex instance.
func (b *finalVertexBuilder[C]) Build() *Vertex[C] {
	vertex := &Vertex[C]{
		id:      b.id,
		kind:    vertexKindFinal,
		onEntry: b.onEntry,
		edges:   newEdgesCollection[C](),
	}

	return vertex
}
