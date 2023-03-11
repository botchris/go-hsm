package examples_test

import (
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChoicePseudoState(t *testing.T) {
	t.Run("should go to c3", func(t *testing.T) {
		context := &choiceCtx{
			g3: true,
		}
		machine, err := prepareChoiceMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.NoError(t, machine.Signal(&choiceSignal{}))
		assert.True(t, machine.At(c3))
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine))
	})

	t.Run("should go to c4", func(t *testing.T) {
		context := &choiceCtx{
			g4: true,
		}
		machine, err := prepareChoiceMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.NoError(t, machine.Signal(&choiceSignal{}))
		assert.True(t, machine.At(c4))
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine))
	})

	t.Run("should go to c5", func(t *testing.T) {
		context := &choiceCtx{
			g5: true,
		}
		machine, err := prepareChoiceMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.NoError(t, machine.Signal(&choiceSignal{}))
		assert.True(t, machine.At(c5))
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine))
	})

	t.Run("should go else branch", func(t *testing.T) {
		context := &choiceCtx{}
		machine, err := prepareChoiceMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.NoError(t, machine.Signal(&choiceSignal{}))
		assert.True(t, machine.At(final))
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine))
	})

	t.Run("multiple valid edges should go trough any one", func(t *testing.T) {
		context := &choiceCtx{
			g3: true,
			g4: false,
			g5: true,
		}
		machine, err := prepareChoiceMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.NoError(t, machine.Signal(&choiceSignal{}))
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*choiceCtx]().Print(machine))
	})
}

func prepareChoiceMachine(context *choiceCtx) (*hsm.HSM[*choiceCtx], error) {
	return hsm.NewBuilder[*choiceCtx]().
		// meta
		WithName("choice").
		WithContext(context).
		StartingAt(c0).
		WithErrorState(hsm.NewErrorState[*choiceCtx]().WithID("error").Build()).

		// states
		AddState(c0).
		AddState(c1).
		AddState(c2).
		AddState(c3).
		AddState(c4).
		AddState(c5).
		AddState(final).

		// build
		Build()
}

// SIGNALS & CONTEXT
type (
	choiceCtx struct {
		g3 bool
		g4 bool
		g5 bool
	}
	choiceSignal struct{}
)

// STATE IDS
var (
	c0ID    = "c0"
	c1ID    = "c1"
	c2ID    = "c2"
	c3ID    = "c3"
	c4ID    = "c4"
	c5ID    = "c5"
	finalID = "end"
)

// MACHINE PARTS
var c0 = hsm.NewStart[*choiceCtx]().
	WithID(c0ID).
	AddTransitions(
		hsm.NewTransition[*choiceCtx]().
			GoTo(c1ID).
			Build(),
	).
	Build()

var c1 = hsm.NewState[*choiceCtx]().
	WithID(c1ID).
	WithTransitions(
		// c1 -choiceSignal-> c2
		hsm.NewTransition[*choiceCtx]().
			When(&choiceSignal{}).
			GoTo(c2ID).
			Build()).
	Build()

var c2 = hsm.NewChoice[*choiceCtx]().
	WithID(c2ID).
	AddTransitions(
		// c2 -> c3
		hsm.NewTransition[*choiceCtx]().
			GoTo(c3ID).
			GuardedBy(
				hsm.NewGuard[*choiceCtx]().
					WithLabel("g3").
					WithMethod(func(ctx *choiceCtx) bool {
						return ctx.g3
					}).
					Build(),
			).
			Build(),
		// c2 -> c4
		hsm.NewTransition[*choiceCtx]().
			GoTo(c4ID).
			GuardedBy(
				hsm.NewGuard[*choiceCtx]().
					WithLabel("g4").
					WithMethod(func(ctx *choiceCtx) bool {
						return ctx.g4
					}).
					Build(),
			).
			Build(),
		// c2 -> c5
		hsm.NewTransition[*choiceCtx]().
			GoTo(c5ID).
			GuardedBy(
				hsm.NewGuard[*choiceCtx]().
					WithLabel("g5").
					WithMethod(func(ctx *choiceCtx) bool {
						return ctx.g5
					}).
					Build(),
			).
			Build(),
		// c2 -[else]-> end
		hsm.NewTransition[*choiceCtx]().
			GoTo(finalID).
			Build(),
	).
	Build()

var c3 = hsm.NewState[*choiceCtx]().
	WithID(c3ID).
	Build()

var c4 = hsm.NewState[*choiceCtx]().
	WithID(c4ID).
	Build()

var c5 = hsm.NewState[*choiceCtx]().
	WithID(c5ID).
	Build()

var final = hsm.NewFinalState[*choiceCtx]().
	WithID(finalID).
	Build()
