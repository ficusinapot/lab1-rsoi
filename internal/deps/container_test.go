package deps

import (
	"testing"

	"github.com/ficusinapot/ds/internal/models/dbifc"
	"github.com/stretchr/testify/require"
)

func TestContainerResolve(t *testing.T) {
	t.Parallel()

	var repository dbifc.PersonRepository = &fakePersonRepository{}
	container := New()
	Add(&container, PersonRepositoryKey, repository)

	resolved, err := Resolve(container, PersonRepositoryKey)

	require.NoError(t, err)
	require.Same(t, repository, resolved)
}

func TestContainerMerge(t *testing.T) {
	t.Parallel()

	var repository dbifc.PersonRepository = &fakePersonRepository{}
	left := New()
	right := New()
	Add(&right, PersonRepositoryKey, repository)

	merged := left.Merge(right)
	resolved, err := Resolve(merged, PersonRepositoryKey)

	require.NoError(t, err)
	require.Same(t, repository, resolved)
}

func TestContainerResolveMissing(t *testing.T) {
	t.Parallel()

	_, err := Resolve(New(), PersonRepositoryKey)

	require.Error(t, err)
}

type fakePersonRepository struct {
	dbifc.PersonRepository
}
