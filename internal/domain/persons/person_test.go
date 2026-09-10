package persons_test

import (
	"context"
	"errors"
	"testing"

	domainpersons "github.com/ficusinapot/ds/internal/domain/persons"
	personcontracts "github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testPersonName = "John"

func TestPersonUseCaseCreatePerson(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{Name: testPersonName, Age: 30, Address: "London", Work: "Engineer"}
	repository := new(personRepositoryMock)
	repository.On("Create", mock.Anything, person).Return(42, nil).Once()

	id, err := domainpersons.NewPersonUseCase(repository).CreatePerson(ctx, person)

	require.NoError(t, err)
	require.Equal(t, 42, id)
	repository.AssertExpectations(t)
}

func TestPersonUseCaseCreatePersonWrapsRepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{Name: testPersonName}
	repositoryErr := errors.New("insert failed")
	repository := new(personRepositoryMock)
	repository.On("Create", mock.Anything, person).Return(0, repositoryErr).Once()

	id, err := domainpersons.NewPersonUseCase(repository).CreatePerson(ctx, person)

	require.Zero(t, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), repositoryErr.Error())
	repository.AssertExpectations(t)
}

func TestPersonUseCaseGetPerson(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{ID: 42, Name: testPersonName}
	repository := new(personRepositoryMock)
	repository.On("GetByID", mock.Anything, 42).Return(person, nil).Once()

	actual, err := domainpersons.NewPersonUseCase(repository).GetPerson(ctx, 42)

	require.NoError(t, err)
	require.Same(t, person, actual)
	repository.AssertExpectations(t)
}

func TestPersonUseCaseGetPersonPreservesNotFoundError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repositoryErr := personcontracts.ErrPersonNotFound.New("person not found")
	repository := new(personRepositoryMock)
	repository.On("GetByID", mock.Anything, 42).Return((*entities.Person)(nil), repositoryErr).Once()

	person, err := domainpersons.NewPersonUseCase(repository).GetPerson(ctx, 42)

	require.Nil(t, person)
	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)
	repository.AssertExpectations(t)
}

func TestPersonUseCaseListPersons(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	persons := []entities.Person{{ID: 1, Name: testPersonName}, {ID: 2, Name: "Jane"}}
	repository := new(personRepositoryMock)
	repository.On("List", mock.Anything).Return(persons, nil).Once()

	actual, err := domainpersons.NewPersonUseCase(repository).ListPersons(ctx)

	require.NoError(t, err)
	require.Equal(t, persons, actual)
	repository.AssertExpectations(t)
}

func TestPersonUseCaseListPersonsWrapsRepositoryError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repositoryErr := errors.New("select failed")
	repository := new(personRepositoryMock)
	repository.On("List", mock.Anything).Return([]entities.Person(nil), repositoryErr).Once()

	persons, err := domainpersons.NewPersonUseCase(repository).ListPersons(ctx)

	require.Nil(t, persons)
	require.Error(t, err)
	require.Contains(t, err.Error(), repositoryErr.Error())
	repository.AssertExpectations(t)
}

func TestPersonUseCaseUpdatePerson(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{ID: 42, Name: testPersonName}
	repository := new(personRepositoryMock)
	repository.On("Update", mock.Anything, person).Return(nil).Once()

	err := domainpersons.NewPersonUseCase(repository).UpdatePerson(ctx, person)

	require.NoError(t, err)
	repository.AssertExpectations(t)
}

func TestPersonUseCaseUpdatePersonPreservesNotFoundError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	person := &entities.Person{ID: 42, Name: testPersonName}
	repositoryErr := personcontracts.ErrPersonNotFound.New("person not found")
	repository := new(personRepositoryMock)
	repository.On("Update", mock.Anything, person).Return(repositoryErr).Once()

	err := domainpersons.NewPersonUseCase(repository).UpdatePerson(ctx, person)

	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)
	repository.AssertExpectations(t)
}

func TestPersonUseCaseDeletePerson(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repository := new(personRepositoryMock)
	repository.On("Delete", mock.Anything, 42).Return(nil).Once()

	err := domainpersons.NewPersonUseCase(repository).DeletePerson(ctx, 42)

	require.NoError(t, err)
	repository.AssertExpectations(t)
}

func TestPersonUseCaseDeletePersonPreservesNotFoundError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repositoryErr := personcontracts.ErrPersonNotFound.New("person not found")
	repository := new(personRepositoryMock)
	repository.On("Delete", mock.Anything, 42).Return(repositoryErr).Once()

	err := domainpersons.NewPersonUseCase(repository).DeletePerson(ctx, 42)

	require.True(t, errorx.IsOfType(err, personcontracts.ErrPersonNotFound), "unexpected error: %v", err)
	repository.AssertExpectations(t)
}

type personRepositoryMock struct {
	mock.Mock
}

func (m *personRepositoryMock) GetByID(ctx context.Context, id int) (*entities.Person, error) {
	args := m.Called(ctx, id)
	person, _ := args.Get(0).(*entities.Person)

	return person, args.Error(1) //nolint:wrapcheck // Mock returns the configured dependency error verbatim.
}

func (m *personRepositoryMock) List(ctx context.Context) ([]entities.Person, error) {
	args := m.Called(ctx)
	persons, _ := args.Get(0).([]entities.Person)

	return persons, args.Error(1) //nolint:wrapcheck // Mock returns the configured dependency error verbatim.
}

func (m *personRepositoryMock) Create(ctx context.Context, person *entities.Person) (int, error) {
	args := m.Called(ctx, person)

	return args.Int(0), args.Error(1) //nolint:wrapcheck // Mock returns the configured dependency error verbatim.
}

func (m *personRepositoryMock) Update(ctx context.Context, person *entities.Person) error {
	args := m.Called(ctx, person)

	return args.Error(0) //nolint:wrapcheck // Mock returns the configured dependency error verbatim.
}

func (m *personRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)

	return args.Error(0) //nolint:wrapcheck // Mock returns the configured dependency error verbatim.
}
