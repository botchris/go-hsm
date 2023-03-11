package examples_test

import (
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// based on the following example: https://en.wikipedia.org/wiki/UML_state_machine#Entry_and_exit_actions
func TestOven(t *testing.T) {
	context := &ovenContext{}
	machine, err := prepareOvenMachine(context)

	//println(string(hsm.NewPlantUMLPrinter[*ovenContext]().Print(machine)))

	require.NoError(t, err)
	require.NotNil(t, machine)

	assert.NoError(t, machine.Signal(&doorClosed{}))
	assert.NotNil(t, context.lamp, "lamp should be off while heating")
	assert.NoError(t, machine.Signal(&doBaking{temp: 120}))
	assert.Equal(t, 120, context.temperature, "baking temp should be 120")
	assert.False(t, machine.Failed())
	assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*ovenContext]().Print(machine))
}

func prepareOvenMachine(context *ovenContext) (*hsm.HSM[*ovenContext], error) {
	return hsm.NewBuilder[*ovenContext]().
		// meta
		WithName("oven").
		WithContext(context).
		StartingAt(doorOpen).
		WithErrorState(
			hsm.NewErrorState[*ovenContext]().
				WithID("error").
				OnEntry(
					hsm.NewAction[*ovenContext]().
						WithLabel("log()").
						WithMethod(func(ctx *ovenContext, signal hsm.Signal) error {
							spew.Dump(signal)

							return nil
						}).
						Build(),
				).
				Build(),
		).

		// states
		AddState(heatingEntry).
		AddState(heating).
		AddState(toasting).
		AddState(baking).
		AddState(doorOpen).

		// build
		Build()
}

// SIGNALS & CONTEXT
type (
	doorOpened struct{}
	doorClosed struct{}
	doToasting struct{}
	doBaking   struct {
		temp int
	}
	ovenContext struct {
		lamp        bool
		temperature int
	}
)

// STATE IDS
var (
	heatingEntryID = "heating entry"
	heatingID      = "heating"
	toastingID     = "toasting"
	bakingID       = "baking"
	doorOpenID     = "door open"
)

// MACHINE PARTS
var heatingEntry = hsm.NewEntryState[*ovenContext]().
	WithID(heatingEntryID).
	AddTransitions(
		hsm.NewTransition[*ovenContext]().
			GoTo(toastingID).
			Build(),
	).
	Build()

var heating = hsm.NewState[*ovenContext]().
	WithID(heatingID).
	WithEntryState(heatingEntry).
	AddTransitions(
		// heating -doorOpened-> door_open
		hsm.NewTransition[*ovenContext]().
			When(&doorOpened{}).
			GoTo(doorOpenID).
			Build(),
		// heating -doToasting-> toasting
		hsm.NewTransition[*ovenContext]().
			When(&doToasting{}).
			GoTo(toastingID).
			Build(),
		// heating -doBaking-> baking
		hsm.NewTransition[*ovenContext]().
			When(&doBaking{}).
			GoTo(bakingID).
			Build(),
	).
	Build()

var toasting = hsm.NewState[*ovenContext]().
	WithID(toastingID).
	ParentOf(heating).
	OnEntry(
		hsm.NewAction[*ovenContext]().
			WithLabel("arm_time_event(oven.toastColor)").
			WithMethod(func(ctx *ovenContext, signal hsm.Signal) error {
				return nil
			}).
			Build(),
	).
	OnExit(
		hsm.NewAction[*ovenContext]().
			WithLabel("disarm_time_event()").
			WithMethod(func(ctx *ovenContext, signal hsm.Signal) error {
				return nil
			}).
			Build()).
	Build()

var baking = hsm.NewState[*ovenContext]().
	WithID(bakingID).
	ParentOf(heating).
	OnEntry(
		hsm.NewAction[*ovenContext]().
			WithLabel("set_temperature(signal.temp)").
			WithMethod(func(ctx *ovenContext, signal hsm.Signal) error {
				if s, ok := signal.(*doBaking); ok {
					ctx.temperature = s.temp
				}

				return nil
			}).
			Build()).
	OnExit(
		hsm.NewAction[*ovenContext]().
			WithLabel("set_temperature(0)").
			WithMethod(func(ctx *ovenContext, signal hsm.Signal) error {
				return nil
			}).
			Build(),
	).
	Build()

var doorOpen = hsm.NewState[*ovenContext]().
	WithID(doorOpenID).
	OnEntry(
		hsm.NewAction[*ovenContext]().
			WithLabel("lamp_on()").
			WithMethod(func(ctx *ovenContext, signal hsm.Signal) error {
				ctx.lamp = true

				return nil
			}).
			Build(),
	).
	OnExit(
		hsm.NewAction[*ovenContext]().
			WithLabel("lamp_off()").
			WithMethod(func(ctx *ovenContext, signal hsm.Signal) error {
				ctx.lamp = false

				return nil
			}).
			Build(),
	).
	AddTransitions(
		// door_open -doorClosed-> heating
		hsm.NewTransition[*ovenContext]().
			When(&doorClosed{}).
			GoTo(heatingID).
			Build(),
	).
	Build()
