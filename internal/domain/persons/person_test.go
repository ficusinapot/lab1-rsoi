package persons_test

import (
	"context"
	"errors"
	"testing"

	domainpersons "github.com/ficusinapot/ds/internal/domain/persons"
	personcontracts "github.com/ficusinapot/ds/internal/models/coreifc/persons"
	dbifcmocks "github.com/ficusinapot/ds/internal/models/dbifc/mocks"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const testPersonName = "John"

func TestPersonUseCaseCreatePerson(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{Name: testPersonName, Age: 30, Address: "London", Work: "Engineer"}
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().Create(gomock.Any(), person).Return(42, nil)

	id, err := domainpersons.NewPersonUseCase(repository).CreatePerson(ctx, person)

	require.NoError(t, err)
	require.Equal(t, 42, id)
}

func TestPersonUseCaseCreatePersonWrapsRepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{Name: testPersonName}
	repositoryErr := errors.New("insert failed")
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().Create(gomock.Any(), person).Return(0, repositoryErr)

	id, err := domainpersons.NewPersonUseCase(repository).CreatePerson(ctx, person)

	require.Zero(t, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), repositoryErr.Error())
}

func TestPersonUseCaseGetPerson(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{ID: 42, Name: testPersonName}
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().GetByID(gomock.Any(), 42).Return(person, nil)

	actual, err := domainpersons.NewPersonUseCase(repository).GetPerson(ctx, 42)

	require.NoError(t, err)
	require.Same(t, person, actual)
}

func TestPersonUseCaseGetPersonPreservesNotFoundError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repositoryErr := personcontracts.ErrPersonNotFound.New("person not found")
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().GetByID(gomock.Any(), 42).Return((*entities.Person)(nil), repositoryErr)

	person, err := domainpersons.NewPersonUseCase(repository).GetPerson(ctx, 42)

	require.Nil(t, person)
	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)
}

func TestPersonUseCaseListPersons(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	persons := []entities.Person{{ID: 1, Name: testPersonName}, {ID: 2, Name: "Jane"}}
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().List(gomock.Any()).Return(persons, nil)

	actual, err := domainpersons.NewPersonUseCase(repository).ListPersons(ctx)

	require.NoError(t, err)
	require.Equal(t, persons, actual)
}

func TestPersonUseCaseListPersonsWrapsRepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repositoryErr := errors.New("select failed")
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().List(gomock.Any()).Return([]entities.Person(nil), repositoryErr)

	persons, err := domainpersons.NewPersonUseCase(repository).ListPersons(ctx)

	require.Nil(t, persons)
	require.Error(t, err)
	require.Contains(t, err.Error(), repositoryErr.Error())
}

func TestPersonUseCaseUpdatePerson(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{ID: 42, Name: testPersonName}
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().Update(gomock.Any(), person).Return(nil)

	err := domainpersons.NewPersonUseCase(repository).UpdatePerson(ctx, person)

	require.NoError(t, err)
}

func TestPersonUseCaseUpdatePersonPreservesNotFoundError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{ID: 42, Name: testPersonName}
	repositoryErr := personcontracts.ErrPersonNotFound.New("person not found")
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().Update(gomock.Any(), person).Return(repositoryErr)

	err := domainpersons.NewPersonUseCase(repository).UpdatePerson(ctx, person)

	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)
}

func TestPersonUseCaseDeletePerson(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().Delete(gomock.Any(), 42).Return(nil)

	err := domainpersons.NewPersonUseCase(repository).DeletePerson(ctx, 42)

	require.NoError(t, err)
}

func TestPersonUseCaseDeletePersonPreservesNotFoundError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repositoryErr := personcontracts.ErrPersonNotFound.New("person not found")
	repository := dbifcmocks.NewMockPersonRepository(gomock.NewController(t))
	repository.EXPECT().Delete(gomock.Any(), 42).Return(repositoryErr)

	err := domainpersons.NewPersonUseCase(repository).DeletePerson(ctx, 42)

	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)
}
