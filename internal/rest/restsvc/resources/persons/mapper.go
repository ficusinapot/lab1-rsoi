package persons

import "github.com/ficusinapot/ds/internal/models/entities"

func requestToModel(request Request) *entities.Person {
	return &entities.Person{
		ID:      0,
		Name:    request.Name,
		Age:     request.Age,
		Address: request.Address,
		Work:    request.Work,
	}
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
