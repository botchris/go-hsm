package hsm

// Effect definition of transition effect
type Effect struct {
	label  string
	method ActionFunc
}

// NewEffect returns a new effect builder
func NewEffect() EffectBuilder {
	return &effectBuilder{}
}
