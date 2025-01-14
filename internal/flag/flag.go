package flag

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type Flag struct {
	Apply          bool
	Clients        []string
	ConfigFilePath string
	Debug          bool
	Help           bool
	Period         int
	Precision      int
	PrintDefaults  func()
	Version        bool
	Workspaces     []string
}

const (
	ErrFlagConvertConfigFilePathError = FlagErr("Cannot convert configuration relative filepath to absolute. Probably HOME environment variable is missing.")
)

type FlagErr string

func (e FlagErr) Error() string {
	return string(e)
}

func InitializeFlags(args []string) (Flag, error) {
	flagSet := pflag.NewFlagSet(args[0], pflag.ExitOnError)

	flag := Flag{
		PrintDefaults: flagSet.PrintDefaults,
	}

	flagSet.BoolVarP(&flag.Apply, "apply", "a", false, "Update jira tasks workload")
	flagSet.BoolVarP(&flag.Debug, "debug", "d", false, "Debug mode - Include already logged time entries")
	flagSet.BoolVarP(&flag.Help, "help", "h", false, "Display help")
	flagSet.BoolVarP(&flag.Version, "version", "v", false, "Show build detials")

	flagSet.StringVar(&flag.ConfigFilePath, "config", "~/.clockify-to-jira/config.yaml", "Config file path")

	flagSet.StringSliceVarP(&flag.Workspaces, "workspace", "w", []string{}, "Filter by workspaceId")
	flagSet.StringSliceVarP(&flag.Clients, "client", "c", []string{}, "Filter by clientId")

	flagSet.IntVarP(&flag.Period, "period", "p", 7, "Migrate time entries from last given days")
	flagSet.IntVarP(&flag.Precision, "tryb-niepokorny", "t", 15, "Rounding up the value of logged time up (minutes)")

	flagSet.Parse(args[1:])

	err := flag.convertConfigFilePathToAbsolute()

	if err != nil {
		return Flag{}, err
	}

	err = flag.validateFlags()

	if err != nil {
		return Flag{}, err
	}

	return flag, nil
}

func (f *Flag) convertConfigFilePathToAbsolute() error {
	dirname, err := os.UserHomeDir()

	if err != nil {
		return ErrFlagConvertConfigFilePathError
	}

	f.ConfigFilePath = strings.Replace(f.ConfigFilePath, "~", dirname, 1)

	return nil
}
