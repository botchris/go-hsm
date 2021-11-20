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

	//println(string(hsm.NewPlantUMLPrinter().Print(machine)))
}

func prepareMutexMachine(context interface{}) (*hsm.HSM, error) {
	return hsm.NewBuilder().
		// meta
		WithName("mutex").
		WithContext(context).
		StartingAt(aState).
		WithErrorState(hsm.NewErrorState().WithID("error").Build()).

		// states
		AddState(aState).
		AddState(bState).

		// build
		Build()
}

// SIGNALS & CONTEXT
type (
	tSignal      struct{}
	mutexContext struct {
		hsm *hsm.HSM
	}
)

// STATE IDS
var (
	aID = "A"
	bID = "B"
)

// MACHINE PARTS
var aState = hsm.NewState().
	WithID(aID).
	AddTransitions(
		// open -handle-> closed
		hsm.NewTransition().
			When(&tSignal{}).
			GoTo(bID).
			Build(),
	).
	Build()

var bState = hsm.NewState().
	WithID(bID).
	OnEntry(
		hsm.NewAction().
			WithLabel("read(hsm.currentState)").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				mutexContext := ctx.(*mutexContext)
				mutexContext.hsm.Current()
				mutexContext.hsm.At(aState)
				mutexContext.hsm.Finished()
				mutexContext.hsm.Failed()
				mutexContext.hsm.Can(&tSignal{})

				return nil
			}).
			Build(),
	).
	Build()
