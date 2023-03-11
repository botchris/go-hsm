package examples_test

import (
	"fmt"
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	t.Run("WHEN onEnter returns error THEN transition fails", func(t *testing.T) {
		context := &errorContext{
			onExitA: func() error {
				return nil
			},
			transitionFx: func() error {
				return nil
			},
			onEntryB: func() error {
				return fmt.Errorf("dummy error, enter B")
			},
		}
		machine, err := prepareErrorMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*errorContext]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.Error(t, machine.Signal(&dummySignal{}))
		assert.Equal(t, 1, context.errorsCount)
		assert.True(t, machine.Failed())
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*errorContext]().Print(machine))
	})

	t.Run("WHEN transition effect returns error THEN transition fails", func(t *testing.T) {
		context := &errorContext{
			onExitA: func() error {
				return nil
			},
			transitionFx: func() error {
				return fmt.Errorf("dummy error, transition FX")
			},
			onEntryB: func() error {
				return nil
			},
		}
		machine, err := prepareErrorMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*errorContext]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.Error(t, machine.Signal(&dummySignal{}))
		assert.Equal(t, 1, context.errorsCount)
		assert.True(t, machine.Failed())
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*errorContext]().Print(machine))
	})

	t.Run("WHEN onExit returns error THEN transition fails", func(t *testing.T) {
		context := &errorContext{
			onExitA: func() error {
				return fmt.Errorf("dummy error, exit A")
			},
			transitionFx: func() error {
				return nil
			},
			onEntryB: func() error {
				return nil
			},
		}
		machine, err := prepareErrorMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*errorContext]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.Error(t, machine.Signal(&dummySignal{}))
		assert.Equal(t, 1, context.errorsCount)
		assert.True(t, machine.Failed())
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*errorContext]().Print(machine))
	})
}

func prepareErrorMachine(context *errorContext) (*hsm.HSM[*errorContext], error) {
	return hsm.NewBuilder[*errorContext]().
		// meta
		WithName("error").
		WithContext(context).
		StartingAt(stateA).
		WithErrorState(
			hsm.NewErrorState[*errorContext]().
				WithID("error").
				OnEntry(
					hsm.NewAction[*errorContext]().
						WithLabel("errorsCount++").
						WithMethod(func(ctx *errorContext, signal hsm.Signal) error {
							ctx.errorsCount++

							return nil
						}).
						Build(),
				).
				Build()).

		// states
		AddState(stateA).
		AddState(stateB).

		// build
		Build()
}

// SIGNALS & CONTEXT
type (
	dummySignal  struct{}
	errorContext struct {
		errorsCount int

		onExitA      func() error
		transitionFx func() error
		onEntryB     func() error
	}
)

// STATE IDS
var (
	stateAID = "A"
	stateBID = "B"
)

// MACHINE PARTS
var stateA = hsm.NewState[*errorContext]().
	WithID(stateAID).
	OnExit(
		hsm.NewAction[*errorContext]().
			WithLabel("ctx.onExitA()").
			WithMethod(func(ctx *errorContext, signal hsm.Signal) error {
				return ctx.onExitA()
			}).
			Build(),
	).
	AddTransitions(
		// A -dummy_signal/dummy_fx-> B
		hsm.NewTransition[*errorContext]().
			When(&dummySignal{}).
			ApplyEffect(
				hsm.NewEffect[*errorContext]().
					WithLabel("ctx.transitionFx()").
					WithMethod(func(ctx *errorContext, signal hsm.Signal) error {
						return ctx.transitionFx()
					}).
					Build(),
			).
			GoTo("B").
			Build(),
	).
	Build()

var stateB = hsm.NewState[*errorContext]().
	WithID(stateBID).
	OnEntry(
		hsm.NewAction[*errorContext]().
			WithLabel("ctx.onEntryB()").
			WithMethod(func(ctx *errorContext, signal hsm.Signal) error {
				return ctx.onEntryB()
			}).
			Build(),
	).
	Build()
