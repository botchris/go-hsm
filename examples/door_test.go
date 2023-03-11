package examples_test

import (
	"fmt"
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoor(t *testing.T) {
	t.Run("WHEN door is closed THEN can be locked", func(t *testing.T) {
		context := &doorContext{}
		machine, err := prepareDoorMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*doorContext]().Print(machine)))

		require.NoError(t, err)
		require.True(t, machine.Can(&handleSignal{}))
		require.NoError(t, machine.Signal(&handleSignal{}))
		require.True(t, machine.Can(&keysSignal{}))
		require.NoError(t, machine.Signal(&keysSignal{}))

		require.Equal(t, []string{
			"exiting OPEN state",
			"entering CLOSED state",
			"exiting CLOSED state",
			"entering LOCKED state",
		}, context.logs)

		assert.Equal(t, 1, context.closedCounter)
		assert.False(t, machine.Failed())
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*doorContext]().Print(machine))
	})

	t.Run("WHEN door is opened THEN cannot be locked using keys", func(t *testing.T) {
		context := &doorContext{}
		machine, err := prepareDoorMachine(context)

		//println(string(hsm.NewPlantUMLPrinter[*doorContext]().Print(machine)))

		require.NoError(t, err)
		assert.Error(t, machine.Signal(&keysSignal{}))
		assert.Empty(t, context.logs)
		assert.False(t, machine.Failed())
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*doorContext]().Print(machine))
	})
}

func prepareDoorMachine(context *doorContext) (*hsm.HSM[*doorContext], error) {
	return hsm.NewBuilder[*doorContext]().
		// meta
		WithName("door").
		WithContext(context).
		StartingAt(openState).
		WithErrorState(hsm.NewErrorState[*doorContext]().WithID("error").Build()).

		// states
		AddState(openState).
		AddState(closedState).
		AddState(lockedState).

		// build
		Build()
}

// SIGNALS & CONTEXT
type (
	handleSignal struct{}
	keysSignal   struct{}
	doorContext  struct {
		logs          []string
		closedCounter int
	}
)

// STATE IDS
var (
	openID   = "open"
	closedID = "closed"
	lockedID = "locked"
)

// MACHINE PARTS
var openState = hsm.NewState[*doorContext]().
	WithID(openID).
	OnEntry(
		hsm.NewAction[*doorContext]().
			WithLabel("log(entering open)").
			WithFunc(func(ctx *doorContext, signal hsm.Signal) error {
				ctx.logs = append(ctx.logs, fmt.Sprintf("entering OPEN state"))

				return nil
			}).
			Build()).
	OnExit(
		hsm.NewAction[*doorContext]().
			WithLabel("log(exiting open)").
			WithFunc(func(ctx *doorContext, signal hsm.Signal) error {
				ctx.logs = append(ctx.logs, fmt.Sprint("exiting OPEN state"))

				return nil
			}).
			Build()).
	WithTransitions(
		// open -handle-> closed
		hsm.NewTransition[*doorContext]().
			When(&handleSignal{}).
			ApplyEffect(countClosedEffect).
			GoTo(closedID).
			Build(),
	).
	Build()

var closedState = hsm.NewState[*doorContext]().
	WithID(closedID).
	OnEntry(
		hsm.NewAction[*doorContext]().
			WithLabel("log(entering close)").
			WithFunc(func(ctx *doorContext, signal hsm.Signal) error {
				ctx.logs = append(ctx.logs, fmt.Sprint("entering CLOSED state"))

				return nil
			}).
			Build(),
	).
	OnExit(
		hsm.NewAction[*doorContext]().
			WithLabel("log(exiting close)").
			WithFunc(func(ctx *doorContext, signal hsm.Signal) error {
				ctx.logs = append(ctx.logs, fmt.Sprint("exiting CLOSED state"))

				return nil
			}).
			Build(),
	).
	WithTransitions(
		// closed -handle-> open
		hsm.NewTransition[*doorContext]().
			When(&handleSignal{}).
			GoTo(openID).
			Build(),
		// closed -keys-> locked
		hsm.NewTransition[*doorContext]().
			When(&keysSignal{}).
			GoTo(lockedID).
			Build(),
	).
	Build()

var lockedState = hsm.NewState[*doorContext]().
	WithID(lockedID).
	OnEntry(hsm.NewAction[*doorContext]().
		WithLabel("log(entering locked)").
		WithFunc(func(ctx *doorContext, signal hsm.Signal) error {
			ctx.logs = append(ctx.logs, fmt.Sprint("entering LOCKED state"))

			return nil
		}).
		Build(),
	).
	OnExit(
		hsm.NewAction[*doorContext]().
			WithLabel("log(exiting locked)").
			WithFunc(func(ctx *doorContext, signal hsm.Signal) error {
				ctx.logs = append(ctx.logs, fmt.Sprint("exiting LOCKED state"))

				return nil
			}).
			Build(),
	).
	WithTransitions(
		// locked -handle-> closed
		hsm.NewTransition[*doorContext]().
			When(&keysSignal{}).
			GoTo(closedID).
			Build(),
	).
	Build()

var countClosedEffect = hsm.NewEffect[*doorContext]().
	WithLabel("closed++").
	WithMethod(func(ctx *doorContext, trigger hsm.Signal) error {
		ctx.closedCounter++

		return nil
	}).
	Build()
