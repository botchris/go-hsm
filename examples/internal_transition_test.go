package examples_test

import (
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternalTransition(t *testing.T) {
	context := &dummyCtx{}
	machine, err := prepareDummyMachine(context)

	//println(string(hsm.NewPlantUMLPrinter().Print(machine)))

	require.NoError(t, err)
	require.NotNil(t, machine)

	assert.NoError(t, machine.Signal(&dummySignalOne{}))
	assert.True(t, machine.At(dummyStateOne))
	assert.Equal(t, 1, context.fxCount)

	assert.NoError(t, machine.Signal(&dummySignalTwo{}))
	assert.True(t, machine.At(dummyStateTwo))
	assert.Equal(t, 1, context.fxCount)
	assert.NotEmpty(t, hsm.NewPlantUMLPrinter().Print(machine))
}

func prepareDummyMachine(context interface{}) (*hsm.HSM, error) {
	return hsm.NewBuilder().
		// meta
		WithName("dummy").
		WithContext(context).
		StartingAt(dummyStateOne).
		WithErrorState(hsm.NewErrorState().WithID("error").Build()).

		// states
		AddState(dummyStateOne).
		AddState(dummyStateTwo).

		// build
		Build()
}

// SIGNALS & CONTEXT
type (
	dummySignalOne struct{}
	dummySignalTwo struct{}
	dummyCtx       struct {
		fxCount int
	}
)

// MACHINE PARTS
var dummyStateOne = hsm.NewState().
	WithID("dummy1").
	AddTransitions(
		hsm.NewInternalTransition().
			When(&dummySignalOne{}).
			ApplyEffect(
				hsm.NewEffect().
					WithLabel("dummy()").
					WithMethod(func(ctx interface{}, signal hsm.Signal) error {
						ctx.(*dummyCtx).fxCount++

						return nil
					}).
					Build(),
			).
			Build(),
	).
	AddTransitions(
		hsm.NewTransition().
			When(&dummySignalTwo{}).
			GoTo("dummy2").
			Build(),
	).
	Build()

var dummyStateTwo = hsm.NewState().
	WithID("dummy2").
	Build()
