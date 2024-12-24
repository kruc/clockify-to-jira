package flag

import (
	"os"
	"testing"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

func initFlagTestsHomeEnvVariable(t *testing.T) {
	err := os.Setenv("HOME", "/home/user")

	if err != nil {
		t.Log(err)
	}
}

func TestInitializeFlagsSuccess(t *testing.T) {

	t.Run("Init default flags", func(t *testing.T) {

		initFlagTestsHomeEnvVariable(t)

		args := []string{os.Args[0]}
		flag, err := InitializeFlags(args)

		assert.Errors(t, err, nil)
		assert.Bools(t, flag.Help, false)
		assert.Bools(t, flag.Debug, false)
		assert.Bools(t, flag.Version, false)
		assert.Bools(t, flag.Apply, false)
		assert.Ints(t, flag.Precision, 15)
		assert.Ints(t, flag.Period, 7)
		assert.Strings(t, flag.ConfigFilePath, "/home/user/.clockify-to-jira/config.yaml")
		assert.StringSlices(t, flag.Workspaces, []string{})
		assert.StringSlices(t, flag.Clients, []string{})
	})

	t.Run("Return error on convert filepath fail", func(t *testing.T) {
		err := os.Unsetenv("HOME")

		if err != nil {
			t.Log(err)
		}

		args := []string{os.Args[0]}

		_, err = InitializeFlags(args)

		assert.Errors(t, err, ErrFlagConvertConfigFilePathError)
	})

	t.Run("Init custom flags - shorthands", func(t *testing.T) {

		initFlagTestsHomeEnvVariable(t)

		args := []string{
			os.Args[0],
			"-c",
			"clientId",
			"-d",
			"-h",
			"-p",
			"31",
			"-t",
			"30",
			"-v",
			"-w",
			"workspaceId",
		}

		flag, err := InitializeFlags(args)

		assert.Errors(t, err, nil)
		assert.Bools(t, flag.Apply, false)
		assert.Bools(t, flag.Debug, true)
		assert.Bools(t, flag.Help, true)
		assert.Ints(t, flag.Period, 31)
		assert.Ints(t, flag.Precision, 30)
		assert.Bools(t, flag.Version, true)
		assert.StringSlices(t, flag.Workspaces, []string{"workspaceId"})
		assert.StringSlices(t, flag.Clients, []string{"clientId"})
	})

	t.Run("Init custom flags - full names", func(t *testing.T) {
		initFlagTestsHomeEnvVariable(t)

		args := []string{
			os.Args[0],
			"--apply",
			"--client",
			"clientId1,clientId2",
			"--config",
			"~/custom-path/.clockify-to-jira/config.yaml",
			"--help",
			"--period",
			"31",
			"--tryb-niepokorny",
			"30",
			"--version",
			"--workspace",
			"workspaceId1,workspaceId2",
		}

		flag, err := InitializeFlags(args)

		assert.Errors(t, err, nil)
		assert.Bools(t, flag.Apply, true)
		assert.Strings(t, flag.ConfigFilePath, "/home/user/custom-path/.clockify-to-jira/config.yaml")
		assert.Bools(t, flag.Debug, false)
		assert.Bools(t, flag.Help, true)
		assert.Ints(t, flag.Period, 31)
		assert.Ints(t, flag.Precision, 30)
		assert.Bools(t, flag.Version, true)
		assert.StringSlices(t, flag.Workspaces, []string{"workspaceId1", "workspaceId2"})
		assert.StringSlices(t, flag.Clients, []string{"clientId1", "clientId2"})
	})
}

func TestValidateFlags(t *testing.T) {
	t.Run("Return error if debug and apply flag are true", func(t *testing.T) {
		initFlagTestsHomeEnvVariable(t)

		args := []string{
			os.Args[0],
			"-a",
			"-d",
		}

		_, err := InitializeFlags(args)

		assert.Errors(t, err, ErrFlagApplyDebugConflict)
	})

	t.Run("Return error if period flag value is less than 1", func(t *testing.T) {
		initFlagTestsHomeEnvVariable(t)

		args := []string{
			os.Args[0],
			"-p",
			"-2",
		}

		_, err := InitializeFlags(args)

		assert.Errors(t, err, ErrFlagPeriodLessThanOne)
	})
}

func TestError(t *testing.T) {
	t.Run("ErrNotFound", func(t *testing.T) {
		got := FlagErr("Error message").Error()
		want := "Error message"

		assert.Strings(t, got, want)
	})
}
