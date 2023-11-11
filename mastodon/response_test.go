package mastodon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalMastoddon(t *testing.T) {
	response := &MastodonData{}

	err := json.Unmarshal([]byte(exampleResponse), response)
	require.NoError(t, err)

}

const exampleResponse string = `
{
    "id": "111375651691866891",
    "created_at": "2023-11-08T15:32:55.000Z",
    "content": "<p>I like to micro blog</p>",
    "account": {
      "username": "sylvia"
    }
}`
