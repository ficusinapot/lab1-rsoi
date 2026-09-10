//go:build integration

package dbrepo_test

import (
	"context"
	"testing"

	"github.com/ficusinapot/ds/internal/db/repos"
	personcontracts "github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/entities"
	"github.com/ficusinapot/ds/tests/infrastructure"

	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/require"
)

func TestPersonRepositoryCRUD(t *testing.T) {
	ctx := context.Background()
	pg := infrastructure.NewPostgres(t)
	repo := repos.NewPersonRepository(pg.DB)

	id, err := repo.Create(ctx, &entities.Person{
		Name:    "Ivan",
		Age:     21,
		Address: "Moscow",
		Work:    "Engineer",
	})
	require.NoError(t, err)
	require.Positive(t, id)

	person, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, &entities.Person{
		ID:      id,
		Name:    "Ivan",
		Age:     21,
		Address: "Moscow",
		Work:    "Engineer",
	}, person)

	person.Name = "Petr"
	person.Age = 0
	person.Address = "Kazan"
	person.Work = "Manager"
	require.NoError(t, repo.Update(ctx, person))

	updated, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, person, updated)

	persons, err := repo.List(ctx)
	require.NoError(t, err)
	require.Contains(t, persons, *person)

	require.NoError(t, repo.Delete(ctx, id))

	_, err = repo.GetByID(ctx, id)
	require.Error(t, err)
	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)
}

func TestPersonRepositoryReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	pg := infrastructure.NewPostgres(t)
	repo := repos.NewPersonRepository(pg.DB)

	_, err := repo.GetByID(ctx, 404)
	require.Error(t, err)
	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)

	err = repo.Update(ctx, &entities.Person{ID: 404, Name: "No", Age: 1, Address: "Nowhere", Work: "None"})
	require.Error(t, err)
	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)

	err = repo.Delete(ctx, 404)
	require.Error(t, err)
	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)
}
