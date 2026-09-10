//go:build e2e

package infrastructure

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func DoJSON[T any](t *testing.T, client http.Client, method, url string, body any, expectedStatus int) T {
	t.Helper()

	response := Do(t, client, method, url, body, expectedStatus)
	defer response.Body.Close()

	var value T
	if err := json.NewDecoder(response.Body).Decode(&value); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	return value
}

func DoNoBody(t *testing.T, client http.Client, method, url string, body any, expectedStatus int) {
	t.Helper()

	response := Do(t, client, method, url, body, expectedStatus)
	defer response.Body.Close()
}

func Do(t *testing.T, client http.Client, method, url string, body any, expectedStatus int) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(payload)
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("send request: %v", err)
	}

	if response.StatusCode != expectedStatus {
		payload, _ := io.ReadAll(response.Body)
		_ = response.Body.Close()
		t.Fatalf("unexpected status: got %d, want %d, body: %s", response.StatusCode, expectedStatus, strings.TrimSpace(string(payload)))
	}

	return response
}
