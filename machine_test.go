package hsm_test

import (
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/require"
)

func TestHSM_Mutex(t *testing.T) {
	ctx := &mutexContext{}
	machine, err := prepareMutexMachine(ctx)

	require.NoError(t, err)
	require.NotNil(t, machine)

	ctx.hsm = machine
	require.NoError(t, machine.Signal(&tSignal{}))
}

func prepareMutexMachine(context *mutexContext) (*hsm.HSM[*mutexContext], error) {
	return hsm.NewBuilder[*mutexContext]().
		// meta
		WithName("mutex").
		WithContext(context).
		StartingAt(aState).
		WithErrorState(hsm.NewErrorState[*mutexContext]().WithID("error").Build()).

		// states
		AddState(aState).
		AddState(bState).

		// build
		Build()
}

// SIGNALS & CONTEXT.
type (
	tSignal      struct{}
	mutexContext struct {
		hsm *hsm.HSM[*mutexContext]
	}
)

// STATE IDS.
var (
	aID = "A"
	bID = "B"
)

// MACHINE PARTS.
var aState = hsm.NewState[*mutexContext]().
	WithID(aID).
	WithTransitions(
		// open -handle-> closed
		hsm.NewTransition[*mutexContext]().
			When(&tSignal{}).
			GoTo(bID).
			Build(),
	).
	Build()

var bState = hsm.NewState[*mutexContext]().
	WithID(bID).
	OnEntry(
		hsm.NewAction[*mutexContext]().
			WithLabel("read(hsm.currentState)").
			WithFunc(func(ctx *mutexContext, signal hsm.Signal) error {
				ctx.hsm.Current()
				ctx.hsm.At(aState)
				ctx.hsm.Finished()
				ctx.hsm.Failed()
				ctx.hsm.Can(&tSignal{})

				return nil
			}).
			Build(),
	).
	Build()
