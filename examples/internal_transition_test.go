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

	//println(string(hsm.NewPlantUMLPrinter[*dummyCtx]().Print(machine)))

	require.NoError(t, err)
	require.NotNil(t, machine)

	assert.NoError(t, machine.Signal(&dummySignalOne{}))
	assert.True(t, machine.At(dummyStateOne))
	assert.Equal(t, 1, context.fxCount)

	assert.NoError(t, machine.Signal(&dummySignalTwo{}))
	assert.True(t, machine.At(dummyStateTwo))
	assert.Equal(t, 1, context.fxCount)
	assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*dummyCtx]().Print(machine))
}

func prepareDummyMachine(context *dummyCtx) (*hsm.HSM[*dummyCtx], error) {
	return hsm.NewBuilder[*dummyCtx]().
		// meta
		WithName("dummy").
		WithContext(context).
		StartingAt(dummyStateOne).
		WithErrorState(hsm.NewErrorState[*dummyCtx]().WithID("error").Build()).

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
var dummyStateOne = hsm.NewState[*dummyCtx]().
	WithID("dummy1").
	AddTransitions(
		hsm.NewInternalTransition[*dummyCtx]().
			When(&dummySignalOne{}).
			ApplyEffect(
				hsm.NewEffect[*dummyCtx]().
					WithLabel("dummy()").
					WithMethod(func(ctx *dummyCtx, signal hsm.Signal) error {
						ctx.fxCount++

						return nil
					}).
					Build(),
			).
			Build(),
	).
	AddTransitions(
		hsm.NewTransition[*dummyCtx]().
			When(&dummySignalTwo{}).
			GoTo("dummy2").
			Build(),
	).
	Build()

var dummyStateTwo = hsm.NewState[*dummyCtx]().
	WithID("dummy2").
	Build()
