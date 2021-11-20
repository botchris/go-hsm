package examples_test

import (
	"testing"

	"github.com/botchris/go-hsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: this tests uses nil-context-machine example definition
func TestTakeSnapshot(t *testing.T) {
	machine, err := prepareNilMachine(nil)

	//println(string(hsm.NewPlantUMLPrinter().Print(machine)))

	require.NoError(t, err)
	require.NotNil(t, machine)

	assert.NoError(t, machine.Signal(&nSignal{}))
	snapshot := machine.Snapshot()
	assert.Equal(t, snapshot.StateID, "n2")
	assert.Equal(t, snapshot.Final, true)
	assert.Equal(t, snapshot.SignalsHistory, []string{"*nSignal"})
	assert.Equal(t, snapshot.StatesHistory, []string{"n1", "n2"})
	assert.False(t, machine.Failed())
	assert.NotEmpty(t, hsm.NewPlantUMLPrinter().Print(machine))
}

// NOTE: this tests uses nil-context-machine example definition
func TestRestoreFromSnapshot(t *testing.T) {
	snapshot := hsm.Snapshot{
		StateID:        "n1",
		Final:          false,
		SignalsHistory: []string{},
		StatesHistory:  []string{},
	}

	machine, err := hsm.NewBuilder().
		// meta
		WithName("nil").
		WithContext(nil).
		StartingAt(n1).
		WithErrorState(hsm.NewErrorState().WithID("error").Build()).

		// states
		AddState(n1).
		AddState(n2).

		// build
		Restore(snapshot)

	require.NoError(t, err)
	require.NotNil(t, machine)

	snapshotAfterBuild := machine.Snapshot()

	assert.Equal(t, snapshotAfterBuild.StateID, "n1")
	assert.Equal(t, snapshotAfterBuild.Final, false)
	assert.Equal(t, snapshotAfterBuild.SignalsHistory, []string{})
	assert.Equal(t, snapshotAfterBuild.StatesHistory, []string{})

	require.NoError(t, err)
	require.NotNil(t, machine)
}
