package repos

import (
	dbent "github.com/ficusinapot/ds/internal/db/ent"
	"github.com/ficusinapot/ds/internal/models/entities"
)

func personFromEnt(person *dbent.Person) *entities.Person {
	return &entities.Person{
		ID:      person.ID,
		Name:    person.Name,
		Age:     person.Age,
		Address: person.Address,
		Work:    person.Work,
	}
}

func personsFromEnt(persons []*dbent.Person) []entities.Person {
	result := make([]entities.Person, 0, len(persons))
	for i := range persons {
		result = append(result, *personFromEnt(persons[i]))
	}

	return result
}
