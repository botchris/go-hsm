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

	//println(string(hsm.NewPlantUMLPrinter[*orderContext]().Print(machine)))

	require.NoError(t, err)
	require.NotNil(t, machine)

	assert.True(t, machine.At(s1))
	assert.True(t, machine.At(s))

	assert.NoError(t, machine.Signal(&tSignal{}))
	sequence := strings.Join(context.calls, "; ") + ";"
	assert.Equal(t, "a(); b(); t(); c(); d(); e();", sequence)
	assert.False(t, machine.Failed())
	assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*orderContext]().Print(machine))
	assert.True(t, machine.At(s21))
	assert.True(t, machine.Can(&wSignal{}))
}

func prepareOrderMachine(context *orderContext) (*hsm.HSM[*orderContext], error) {
	return hsm.NewBuilder[*orderContext]().
		// meta
		WithName("order").
		WithContext(context).
		StartingAt(s11).
		WithErrorState(hsm.NewErrorState[*orderContext]().WithID("error").Build()).

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

var w = hsm.NewState[*orderContext]().
	WithID("w").
	Build()

var s = hsm.NewState[*orderContext]().
	WithID("s").
	OnExit(
		hsm.NewAction[*orderContext]().
			WithLabel("f()").
			WithFunc(func(ctx *orderContext, signal hsm.Signal) error {
				ctx.calls = append(ctx.calls, "f()")

				return nil
			}).
			Build(),
	).
	WithTransitions(
		hsm.NewTransition[*orderContext]().
			When(&wSignal{}).
			GoTo("w").
			Build(),
	).
	Build()

var s1 = hsm.NewState[*orderContext]().
	WithID("s1").
	ParentOf(s).
	OnExit(
		hsm.NewAction[*orderContext]().
			WithLabel("b()").
			WithFunc(func(ctx *orderContext, signal hsm.Signal) error {
				ctx.calls = append(ctx.calls, "b()")
				return nil
			}).
			Build(),
	).
	WithTransitions(
		hsm.NewTransition[*orderContext]().
			When(&tSignal{}).
			GuardedBy(gGuard).
			ApplyEffect(tEffect).
			GoTo("s2").
			Build(),
	).
	Build()

var s11 = hsm.NewState[*orderContext]().
	WithID("s11").
	ParentOf(s1).
	OnExit(
		hsm.NewAction[*orderContext]().
			WithLabel("a()").
			WithFunc(func(ctx *orderContext, signal hsm.Signal) error {
				ctx.calls = append(ctx.calls, "a()")

				return nil
			}).
			Build(),
	).
	Build()

var s2 = hsm.NewState[*orderContext]().
	WithID("s2").
	ParentOf(s).
	WithEntryState(s2Entry).
	OnEntry(
		hsm.NewAction[*orderContext]().
			WithLabel("c()").
			WithFunc(func(ctx *orderContext, signal hsm.Signal) error {
				ctx.calls = append(ctx.calls, "c()")

				return nil
			}).
			Build(),
	).
	Build()

var s2Entry = hsm.NewEntryState[*orderContext]().
	WithID("s2 entry").
	WithTransitions(
		hsm.NewTransition[*orderContext]().
			ApplyEffect(dEffect).
			GoTo("s21").
			Build(),
	).
	Build()

var s21 = hsm.NewState[*orderContext]().
	WithID("s21").
	ParentOf(s2).
	OnEntry(
		hsm.NewAction[*orderContext]().
			WithLabel("e()").
			WithFunc(func(ctx *orderContext, signal hsm.Signal) error {
				ctx.calls = append(ctx.calls, "e()")

				return nil
			}).
			Build(),
	).
	Build()

var tEffect = hsm.NewEffect[*orderContext]().
	WithLabel("t()").
	WithMethod(func(ctx *orderContext, signal hsm.Signal) error {
		ctx.calls = append(ctx.calls, "t()")

		return nil
	}).
	Build()

var dEffect = hsm.NewEffect[*orderContext]().
	WithLabel("d()").
	WithMethod(func(ctx *orderContext, signal hsm.Signal) error {
		ctx.calls = append(ctx.calls, "d()")

		return nil
	}).
	Build()

var gGuard = hsm.NewGuard[*orderContext]().
	WithLabel("g()").
	WithMethod(func(ctx *orderContext) bool {
		return true
	}).
	Build()
