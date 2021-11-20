package examples_test

import (
	"strings"
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInvocationOrder based on example at: https://en.wikipedia.org/wiki/UML_state_machine#Transition_execution_sequence
func TestInvocationOrder(t *testing.T) {
	context := &orderContext{}
	machine, err := prepareOrderMachine(context)

	//println(string(hsm.NewPlantUMLPrinter().Print(machine)))

	require.NoError(t, err)
	require.NotNil(t, machine)

	assert.True(t, machine.At(s1))
	assert.True(t, machine.At(s))

	assert.NoError(t, machine.Signal(&tSignal{}))
	sequence := strings.Join(context.calls, "; ") + ";"
	assert.Equal(t, "a(); b(); t(); c(); d(); e();", sequence)
	assert.False(t, machine.Failed())
	assert.NotEmpty(t, hsm.NewPlantUMLPrinter().Print(machine))
	assert.True(t, machine.At(s21))
	assert.True(t, machine.Can(&wSignal{}))
}

func prepareOrderMachine(context interface{}) (*hsm.HSM, error) {
	return hsm.NewBuilder().
		// meta
		WithName("order").
		WithContext(context).
		StartingAt(s11).
		WithErrorState(hsm.NewErrorState().WithID("error").Build()).

		// states
		AddState(w).
		AddState(s).
		AddState(s1).
		AddState(s11).
		AddState(s2).
		AddState(s21).

		// build
		Build()
}

type (
	tSignal      struct{}
	wSignal      struct{}
	orderContext struct {
		calls []string
	}
)

var w = hsm.NewState().
	WithID("w").
	Build()

var s = hsm.NewState().
	WithID("s").
	OnExit(
		hsm.NewAction().
			WithLabel("f()").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				ctx.(*orderContext).calls = append(ctx.(*orderContext).calls, "f()")
				return nil
			}).
			Build(),
	).
	AddTransitions(
		hsm.NewTransition().
			When(&wSignal{}).
			GoTo("w").
			Build(),
	).
	Build()

var s1 = hsm.NewState().
	WithID("s1").
	ParentOf(s).
	OnExit(
		hsm.NewAction().
			WithLabel("b()").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				ctx.(*orderContext).calls = append(ctx.(*orderContext).calls, "b()")
				return nil
			}).
			Build(),
	).
	AddTransitions(
		hsm.NewTransition().
			When(&tSignal{}).
			GuardedBy(gGuard).
			ApplyEffect(tEffect).
			GoTo("s2").
			Build(),
	).
	Build()

var s11 = hsm.NewState().
	WithID("s11").
	ParentOf(s1).
	OnExit(
		hsm.NewAction().
			WithLabel("a()").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				ctx.(*orderContext).calls = append(ctx.(*orderContext).calls, "a()")
				return nil
			}).
			Build(),
	).
	Build()

var s2 = hsm.NewState().
	WithID("s2").
	ParentOf(s).
	WithEntryState(s2Entry).
	OnEntry(
		hsm.NewAction().
			WithLabel("c()").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				ctx.(*orderContext).calls = append(ctx.(*orderContext).calls, "c()")
				return nil
			}).
			Build(),
	).
	Build()

var s2Entry = hsm.NewEntryState().
	WithID("s2 entry").
	AddTransitions(
		hsm.NewTransition().
			ApplyEffect(dEffect).
			GoTo("s21").
			Build(),
	).
	Build()

var s21 = hsm.NewState().
	WithID("s21").
	ParentOf(s2).
	OnEntry(
		hsm.NewAction().
			WithLabel("e()").
			WithMethod(func(ctx interface{}, signal hsm.Signal) error {
				ctx.(*orderContext).calls = append(ctx.(*orderContext).calls, "e()")
				return nil
			}).
			Build(),
	).
	Build()

var tEffect = hsm.NewEffect().
	WithLabel("t()").
	WithMethod(func(ctx interface{}, signal hsm.Signal) error {
		ctx.(*orderContext).calls = append(ctx.(*orderContext).calls, "t()")
		return nil
	}).
	Build()

var dEffect = hsm.NewEffect().
	WithLabel("d()").
	WithMethod(func(ctx interface{}, signal hsm.Signal) error {
		ctx.(*orderContext).calls = append(ctx.(*orderContext).calls, "d()")
		return nil
	}).
	Build()

var gGuard = hsm.NewGuard().
	WithLabel("g()").
	WithMethod(func(ctx interface{}) bool {
		return true
	}).
	Build()
