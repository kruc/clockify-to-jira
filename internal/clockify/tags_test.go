package clockify

import (
	"testing"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

func TestGetWorkspaceTags(t *testing.T) {

	t.Run("Get workspace tags", func(t *testing.T) {
		fakeClient := &fakeClient{}
		fakeClient.getTagsSuccessResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")

		tags, err := apiClient.GetWorkspaceTags("workspaceId1")

		assert.Errors(t, err, nil)
		assert.Ints(t, len(tags), 2)
		assert.Strings(t, tags["tag1"].ID, "tagId1")
		assert.Strings(t, tags["tag2"].ID, "tagId2")
	})

	t.Run("Get error on get workspace", func(t *testing.T) {
		fakeClient := &fakeClient{}
		fakeClient.getTagsErrorResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")

		_, err := apiClient.GetWorkspaceTags("workspaceId1")

		assert.Errors(t, err, ErrClockifyFailToFetchWorkspaceTags)
	})
}
