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

		//println(string(hsm.NewPlantUMLPrinter().Print(machine)))

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
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter().Print(machine))
	})

	t.Run("WHEN door is opened THEN cannot be locked using keys", func(t *testing.T) {
		context := &doorContext{}
		machine, err := prepareDoorMachine(context)

		//println(string(hsm.NewPlantUMLPrinter().Print(machine)))

		require.NoError(t, err)
		assert.Error(t, machine.Signal(&keysSignal{}))
		assert.Empty(t, context.logs)
		assert.False(t, machine.Failed())
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter().Print(machine))
	})
}

func prepareDoorMachine(context interface{}) (*hsm.HSM, error) {
	return hsm.NewBuilder().
		// meta
		WithName("door").
		WithContext(context).
		StartingAt(openState).
		WithErrorState(hsm.NewErrorState().WithID("error").Build()).

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
var openState = hsm.NewState().
	WithID(openID).
	OnEntry(
		hsm.NewAction().
			WithLabel("log(entering open)").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				context := ctx.(*doorContext)
				context.logs = append(context.logs, fmt.Sprintf("entering OPEN state"))

				return nil
			}).
			Build()).
	OnExit(
		hsm.NewAction().
			WithLabel("log(exiting open)").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				context := ctx.(*doorContext)
				context.logs = append(context.logs, fmt.Sprint("exiting OPEN state"))

				return nil
			}).
			Build()).
	AddTransitions(
		// open -handle-> closed
		hsm.NewTransition().
			When(&handleSignal{}).
			ApplyEffect(countClosedEffect).
			GoTo(closedID).
			Build()).
	Build()

var closedState = hsm.NewState().
	WithID(closedID).
	OnEntry(
		hsm.NewAction().
			WithLabel("log(entering close)").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				context := ctx.(*doorContext)
				context.logs = append(context.logs, fmt.Sprint("entering CLOSED state"))

				return nil
			}).
			Build()).
	OnExit(
		hsm.NewAction().
			WithLabel("log(exiting close)").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				context := ctx.(*doorContext)
				context.logs = append(context.logs, fmt.Sprint("exiting CLOSED state"))

				return nil
			}).
			Build()).
	AddTransitions(
		// closed -handle-> open
		hsm.NewTransition().
			When(&handleSignal{}).
			GoTo(openID).
			Build(),
		// closed -keys-> locked
		hsm.NewTransition().
			When(&keysSignal{}).
			GoTo(lockedID).
			Build()).
	Build()

var lockedState = hsm.NewState().
	WithID(lockedID).
	OnEntry(hsm.NewAction().
		WithLabel("log(entering locked)").
		WithMethod(func(ctx interface{}, signal hsm.Signal) error {
			context := ctx.(*doorContext)
			context.logs = append(context.logs, fmt.Sprint("entering LOCKED state"))

			return nil
		}).
		Build()).
	OnExit(
		hsm.NewAction().
			WithLabel("log(exiting locked)").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				context := ctx.(*doorContext)
				context.logs = append(context.logs, fmt.Sprint("exiting LOCKED state"))

				return nil
			}).
			Build()).
	AddTransitions(
		// locked -handle-> closed
		hsm.NewTransition().
			When(&keysSignal{}).
			GoTo(closedID).
			Build()).
	Build()

var countClosedEffect = hsm.NewEffect().
	WithLabel("closed++").
	WithMethod(func(ctx interface{}, trigger hsm.Signal) error {
		ctx.(*doorContext).closedCounter++
		return nil
	}).
	Build()
