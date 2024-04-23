package cli

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testAppID                = "test-app-id"
	testAppDesc              = "test app desc"
	testCmdName              = "test-cmd"
	testCmdDesc              = "test cmd desc"
	testCommandTickerOutput  = "ticker-ticked"
	testCommandCtxDoneOutput = "ctx-done"
	testCommandOutput        string
)

var testCommandHandler = func(parentCtx context.Context) {
	ticker := time.NewTicker(time.Millisecond)
	select {
	case <-ticker.C:
		if testCommandOutput == "" {
			testCommandOutput = testCommandTickerOutput
		}
	case <-parentCtx.Done():
		if testCommandOutput == "" {
			testCommandOutput = testCommandCtxDoneOutput
		}
	}
}

func reset() {
	cmd = nil
	testCommandOutput = ""
}

func TestInit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer reset()
		assert.Nil(t, cmd)

		expectedCmdName := testAppID
		expectedCmdDesc := testAppDesc
		Init(expectedCmdName, testAppDesc)

		assert.Equal(t, expectedCmdName, cmd.Name())
		assert.Equal(t, expectedCmdDesc, cmd.Short)
		assert.Equal(t, expectedCmdDesc, cmd.Long)
	})
}

func TestAddCommand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer reset()
		Init(testAppID, testAppDesc)

		expectedCmdName := testCmdName
		expectedCmdDesc := testCmdDesc
		err := AddCommand(expectedCmdName, expectedCmdDesc, func(_ context.Context) {})
		assert.Nil(t, err)

		foundCmd, _, err := cmd.Root().Traverse([]string{testCmdName})
		assert.Nil(t, err)

		assert.Equal(t, expectedCmdName, foundCmd.Name())
		assert.Equal(t, expectedCmdDesc, foundCmd.Short)
		assert.Equal(t, expectedCmdDesc, foundCmd.Long)
	})

	t.Run("returnsErrWhenNotInitialized", func(t *testing.T) {
		defer reset()
		expected := ErrNotInitialized
		actual := AddCommand(testCmdName, testCmdDesc, testCommandHandler)
		assert.Equal(t, expected, actual)
	})
}

func TestAssignStringFlag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer reset()
		Init(testAppID, testAppDesc)
		cmd.SetOut(bytes.NewBuffer(make([]byte, 0)))

		var actual string
		expected := "test-val"
		flagName := "test-flag"
		flagDefaultValue := "default-test-val"
		flagDescription := "test-flag-description"

		AssignStringFlag(&actual, flagName, flagDefaultValue, flagDescription)

		err := cmd.Execute()
		assert.Nil(t, err)
		assert.Equal(t, flagDefaultValue, actual)

		cmd.SetArgs([]string{"--" + flagName, expected})
		err = cmd.Execute()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestRun(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer reset()
		Init(testAppID, testAppDesc)

		_ = AddCommand(testCmdName, testCmdDesc, testCommandHandler)
		cmd.SetArgs([]string{testCmdName})

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Nanosecond)
		err := Run(ctx)
		assert.Nil(t, err)
		cancelFunc()
		assert.Equal(t, testCommandCtxDoneOutput, testCommandOutput)

		reset()
		Init(testAppID, testAppDesc)

		_ = AddCommand(testCmdName, testCmdDesc, testCommandHandler)
		cmd.SetArgs([]string{testCmdName})

		ctx, cancelFunc = context.WithTimeout(context.Background(), time.Millisecond*10)
		err = Run(ctx)
		assert.Nil(t, err)
		cancelFunc()
		assert.Equal(t, testCommandTickerOutput, testCommandOutput)
	})
}
