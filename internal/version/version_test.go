package version

import (
	"bytes"
	"testing"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

func TestDisplayVersion(t *testing.T) {
	// Provide as ldflags
	// var (
	// 	BuildVersion = "v1.2.3"
	// 	BuildDate    = "08.10.2022 09:19:08"
	// 	GitCommit    = "3dfbf35"
	// )
	//
	// For VSCode users

	// settings.json

	// "go.testFlags": [
	//   "-ldflags",
	//   "-X 'github.com/kruc/clockify-to-jira/internal/version.BuildVersion=v1.2.3' -X 'github.com/kruc/clockify-to-jira/internal/version.GitCommit=3dfbf35' -X 'github.com/kruc/clockify-to-jira/internal/version.BuildDate=08.10.2022 09:19:08'",
	// ]

	t.Run("Dispaly version details", func(t *testing.T) {
		var buffer bytes.Buffer

		ShowBuildDetails(&buffer)

		assert.Strings(t, buffer.String(),
			"BuildVersion: v1.2.3\tBuildDate: 08.10.2022 09:19:08\tGitCommit: 3dfbf35\n")
	})
}
