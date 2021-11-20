package examples_test

import (
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// state machine with no context associated
func TestNilContext(t *testing.T) {
	machine, err := prepareNilMachine(nil)

	//println(string(hsm.NewPlantUMLPrinter().Print(machine)))

	require.NoError(t, err)
	require.NotNil(t, machine)

	assert.NoError(t, machine.Signal(&nSignal{}))
	assert.False(t, machine.Failed())
	assert.NotEmpty(t, hsm.NewPlantUMLPrinter().Print(machine))
}

func prepareNilMachine(context interface{}) (*hsm.HSM, error) {
	return hsm.NewBuilder().
		// meta
		WithName("nil").
		WithContext(context).
		StartingAt(n1).
		WithErrorState(hsm.NewErrorState().WithID("error").Build()).

		// states
		AddState(n1).
		AddState(n2).

		// build
		Build()
}

type (
	nSignal struct{}
)

var (
	n1ID = "n1"
	n2ID = "n2"
)

var n1 = hsm.NewState().
	WithID(n1ID).
	AddTransitions(
		hsm.NewTransition().
			When(&nSignal{}).
			GoTo(n2ID).
			Build(),
	).
	Build()

var n2 = hsm.NewState().
	WithID(n2ID).
	Build()
