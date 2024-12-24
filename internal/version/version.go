package version

import (
	"fmt"
	"io"
)

var (
	BuildVersion string
	BuildDate    string
	GitCommit    string
	VersionFlag  bool
)

func ShowBuildDetails(w io.Writer) {
	fmt.Fprintf(w, "BuildVersion: %s\tBuildDate: %s\tGitCommit: %s\n", BuildVersion, BuildDate, GitCommit)
}
