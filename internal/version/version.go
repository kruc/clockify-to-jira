package version

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

var (
	BuildVersion string
	BuildDate    string
	GitCommit    string
	VersionFlag  bool
)

func init() {
	flag.BoolVarP(&VersionFlag, "version", "v", false, "Display version")
}

func DisplayVersion() {
	fmt.Printf("BuildVersion: %s\tBuildDate: %s\tGitCommit: %s\n", BuildVersion, BuildDate, GitCommit)
}
