package hsm

// Effect definition of transition effect.
type Effect[C any] struct {
	label  string
	method ActionFunc[C]
}

// NewEffect returns a new effect builder.
func NewEffect[C any]() EffectBuilder[C] {
	return &effectBuilder[C]{}
}
