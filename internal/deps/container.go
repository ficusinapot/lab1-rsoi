package deps

import (
	"database/sql"

	"github.com/ficusinapot/ds/internal/models/coreifc/persons"
	"github.com/ficusinapot/ds/internal/models/coreifc/status"
	"github.com/ficusinapot/ds/internal/models/dbifc"
	"github.com/joomcode/errorx"
)

type Key[T any] struct {
	name string
}

var (
	DBKey               = Key[*sql.DB]{name: "database"}
	PersonRepositoryKey = Key[dbifc.PersonRepository]{name: "person repository"}
	StatusProviderKey   = Key[dbifc.StatusProvider]{name: "status provider"}
	PersonUseCaseKey    = Key[persons.PersonUseCase]{name: "person use case"}
	StatusUseCaseKey    = Key[status.StatusUseCase]{name: "status use case"}
)

type Container struct {
	values map[string]any
}

func New() Container {
	return Container{
		values: make(map[string]any),
	}
}

func Add[T any](container *Container, key Key[T], value T) {
	if container.values == nil {
		container.values = make(map[string]any)
	}

	container.values[key.name] = value
}

func Resolve[T any](container Container, key Key[T]) (T, error) {
	var zero T
	if container.values == nil {
		return zero, errorx.IllegalState.New("dependency container is empty")
	}

	raw, ok := container.values[key.name]
	if !ok || raw == nil {
		return zero, errorx.IllegalState.New("dependency %s is not registered", key.name)
	}

	value, ok := raw.(T)
	if !ok {
		return zero, errorx.IllegalState.New("dependency %s has invalid type", key.name)
	}

	return value, nil
}

func (c Container) Merge(other Container) Container {
	merged := New()
	for key, value := range c.values {
		merged.values[key] = value
	}
	for key, value := range other.values {
		merged.values[key] = value
	}

	return merged
}
