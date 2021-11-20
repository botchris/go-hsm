package examples_test

import (
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChoice(t *testing.T) {
	t.Run("WHEN two players join THEN lobby must start playing", func(t *testing.T) {
		machine, err := prepareLobbyMachine(&lobbyContext{})

		//println(string(hsm.NewPlantUMLPrinter().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.NoError(t, machine.Signal(&playerJoined{}))
		assert.NoError(t, machine.Signal(&playerJoined{}))
		assert.True(t, machine.At(playing))
		assert.True(t, machine.Finished())
		assert.False(t, machine.Failed())
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter().Print(machine))
	})
}

func prepareLobbyMachine(context *lobbyContext) (*hsm.HSM, error) {
	return hsm.NewBuilder().
		// meta
		WithName("lobby").
		WithContext(context).
		StartingAt(start).
		WithErrorState(hsm.NewErrorState().WithID("error").Build()).

		// states
		AddState(start).
		AddState(awaiting).
		AddState(startChoice).
		AddState(playing).

		// build
		Build()
}

// SIGNALS & CONTEXT
type (
	playerJoined struct{}
	lobbyContext struct {
		players int
	}
)

// STATE IDS
var (
	startID       = "start"
	awaitingID    = "awaiting"
	startChoiceID = "start choice"
	playingID     = "playing"
)

// MACHINE PARTS
var start = hsm.NewStart().
	WithID(startID).
	AddTransitions(
		hsm.NewTransition().
			GoTo(awaitingID).
			Build(),
	).
	Build()

var awaiting = hsm.NewState().
	WithID(awaitingID).
	AddTransitions(
		// awaiting -join-> <<start choice>>
		hsm.NewTransition().
			When(&playerJoined{}).
			ApplyEffect(joinEffect).
			GoTo(startChoiceID).
			Build()).
	Build()

var startChoice = hsm.NewChoice().
	WithID(startChoiceID).
	AddTransitions(
		// <<start>> -[#players >= 2]-> playing
		hsm.NewTransition().
			GuardedBy(enoughPlayersGuard).
			GoTo(playingID).
			Build(),
		// <<start>> -[else]-> awaiting
		hsm.NewTransition().
			GoTo(awaitingID).
			Build()).
	Build()

var playing = hsm.NewState().
	WithID(playingID).
	Build()

var enoughPlayersGuard = hsm.NewGuard().
	WithLabel("#players >= 2").
	WithMethod(func(ctx interface{}) bool {
		lobby := ctx.(*lobbyContext)

		return lobby.players >= 2
	}).
	Build()

var joinEffect = hsm.NewEffect().
	WithLabel("lobby.players++").
	WithMethod(func(ctx interface{}, trigger hsm.Signal) error {
		ctx.(*lobbyContext).players++
		return nil
	}).
	Build()
