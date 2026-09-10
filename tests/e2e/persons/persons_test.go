//go:build e2e

package persons_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ficusinapot/ds/tests/e2e/infrastructure"

	"github.com/stretchr/testify/require"
)

type person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type personRequest struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

func TestPersonsAPI(t *testing.T) {
	service := infrastructure.NewService(t)
	client := http.Client{Timeout: 5 * time.Second}

	createResponse := infrastructure.Do(t, client, http.MethodPost, service.BaseURL+"/api/v1/persons", personRequest{
		Name:    "Ivan",
		Age:     21,
		Address: "Moscow",
		Work:    "Engineer",
	}, http.StatusCreated)
	defer createResponse.Body.Close()

	location := createResponse.Header.Get("Location")
	require.NotEmpty(t, location)

	personID, err := strconv.Atoi(strings.TrimPrefix(location, "/api/v1/persons/"))
	require.NoError(t, err)

	var created person
	require.NoError(t, json.NewDecoder(createResponse.Body).Decode(&created))
	require.Positive(t, created.ID)
	require.Equal(t, personID, created.ID)

	got := infrastructure.DoJSON[person](t, client, http.MethodGet, fmt.Sprintf("%s/api/v1/persons/%d", service.BaseURL, created.ID), nil, http.StatusOK)
	require.Equal(t, created, got)

	list := infrastructure.DoJSON[[]person](t, client, http.MethodGet, service.BaseURL+"/api/v1/persons", nil, http.StatusOK)
	require.Contains(t, list, created)

	updated := infrastructure.DoJSON[person](t, client, http.MethodPatch, fmt.Sprintf("%s/api/v1/persons/%d", service.BaseURL, created.ID), map[string]string{
		"name":    "Petr",
		"address": "Kazan",
	}, http.StatusOK)
	require.Equal(t, person{
		ID:      created.ID,
		Name:    "Petr",
		Age:     21,
		Address: "Kazan",
		Work:    "Engineer",
	}, updated)

	infrastructure.DoNoBody(t, client, http.MethodDelete, fmt.Sprintf("%s/api/v1/persons/%d", service.BaseURL, created.ID), nil, http.StatusNoContent)
	infrastructure.DoNoBody(t, client, http.MethodGet, fmt.Sprintf("%s/api/v1/persons/%d", service.BaseURL, created.ID), nil, http.StatusNotFound)
}
