package restsvc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitErrorMessage(t *testing.T) {
	t.Parallel()

	title, detail, ok := splitErrorMessage("Invalid request: person id is required")

	require.True(t, ok)
	require.Equal(t, "Invalid request", title)
	require.Equal(t, "person id is required", detail)
}

func TestSplitErrorMessageRejectsPlainMessage(t *testing.T) {
	t.Parallel()

	_, _, ok := splitErrorMessage("person id is required")

	require.False(t, ok)
}
