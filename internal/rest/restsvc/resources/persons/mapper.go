package persons

import "github.com/ficusinapot/ds/internal/models/entities"

func requestToModel(request Request) *entities.Person {
	return &entities.Person{
		ID:      0,
		Name:    valueOrZero(request.Name),
		Age:     valueOrZero(request.Age),
		Address: valueOrZero(request.Address),
		Work:    valueOrZero(request.Work),
	}
}

func mergeRequest(person *entities.Person, request Request) *entities.Person {
	if request.Name != nil {
		person.Name = *request.Name
	}
	if request.Age != nil {
		person.Age = *request.Age
	}
	if request.Address != nil {
		person.Address = *request.Address
	}
	if request.Work != nil {
		person.Work = *request.Work
	}

	return person
}

func valueOrZero[T any](value *T) T {
	if value == nil {
		var zero T
		return zero
	}

	return *value
}

func modelToResponse(person *entities.Person) Response {
	return Response{
		ID:      person.ID,
		Name:    person.Name,
		Age:     person.Age,
		Address: person.Address,
		Work:    person.Work,
	}
}

func modelsToResponses(persons []entities.Person) []Response {
	responses := make([]Response, 0, len(persons))
	for i := range persons {
		responses = append(responses, modelToResponse(&persons[i]))
	}

	return responses
}
