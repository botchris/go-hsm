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

		//println(string(hsm.NewPlantUMLPrinter[*lobbyContext]().Print(machine)))

		require.NoError(t, err)
		require.NotNil(t, machine)

		assert.NoError(t, machine.Signal(&playerJoined{}))
		assert.NoError(t, machine.Signal(&playerJoined{}))
		assert.True(t, machine.At(playing))
		assert.True(t, machine.Finished())
		assert.False(t, machine.Failed())
		assert.NotEmpty(t, hsm.NewPlantUMLPrinter[*lobbyContext]().Print(machine))
	})
}

func prepareLobbyMachine(context *lobbyContext) (*hsm.HSM[*lobbyContext], error) {
	return hsm.NewBuilder[*lobbyContext]().
		// meta
		WithName("lobby").
		WithContext(context).
		StartingAt(start).
		WithErrorState(hsm.NewErrorState[*lobbyContext]().WithID("error").Build()).

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
var start = hsm.NewStart[*lobbyContext]().
	WithID(startID).
	AddTransitions(
		hsm.NewTransition[*lobbyContext]().
			GoTo(awaitingID).
			Build(),
	).
	Build()

var awaiting = hsm.NewState[*lobbyContext]().
	WithID(awaitingID).
	WithTransitions(
		// awaiting -join-> <<start choice>>
		hsm.NewTransition[*lobbyContext]().
			When(&playerJoined{}).
			ApplyEffect(joinEffect).
			GoTo(startChoiceID).
			Build()).
	Build()

var startChoice = hsm.NewChoice[*lobbyContext]().
	WithID(startChoiceID).
	AddTransitions(
		// <<start>> -[#players >= 2]-> playing
		hsm.NewTransition[*lobbyContext]().
			GuardedBy(enoughPlayersGuard).
			GoTo(playingID).
			Build(),
		// <<start>> -[else]-> awaiting
		hsm.NewTransition[*lobbyContext]().
			GoTo(awaitingID).
			Build(),
	).
	Build()

var playing = hsm.NewState[*lobbyContext]().
	WithID(playingID).
	Build()

var enoughPlayersGuard = hsm.NewGuard[*lobbyContext]().
	WithLabel("#players >= 2").
	WithMethod(func(ctx *lobbyContext) bool {
		return ctx.players >= 2
	}).
	Build()

var joinEffect = hsm.NewEffect[*lobbyContext]().
	WithLabel("lobby.players++").
	WithMethod(func(ctx *lobbyContext, trigger hsm.Signal) error {
		ctx.players++

		return nil
	}).
	Build()
