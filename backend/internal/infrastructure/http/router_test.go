//go:build integration

package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	tests := []struct {
		name         string
		request      map[string]any
		responseCode int
	}{
		{
			name: "Success to create user",
			request: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": "password",
			},
			responseCode: http.StatusCreated,
		},
		{
			name: "Empty password",
			request: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": "",
			},
			responseCode: http.StatusBadRequest,
		},
		{
			name: "password length under 8 characters",
			request: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": "passwor",
			},
			responseCode: http.StatusBadRequest,
		},
		{
			name: "Illegal Request",
			request: map[string]any{
				"email":    "test@example.com",
				"password": "passwor",
			},
			responseCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.request)
			if err != nil {
				assert.FailNow(t, "fail to marshal json", err)
			}

			resp, err := http.Post(testServer.URL+"/signup", "application/json", bytes.NewBuffer(body))
			if err != nil {
				assert.FailNow(t, "fail to request", err)
			}

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					assert.Fail(t, "failed to close body", err)
				}
			}()

			assert.Equal(t, resp.StatusCode, tt.responseCode)

			err = testDb.Cleanup()
			if err != nil {
				assert.Fail(t, "fail to cleanup testDb")
			}
		})
	}
}

func TestSignup_DuplicateRequest(t *testing.T) {
	request := map[string]any{
		"name":     "Test",
		"email":    "test@example.com",
		"password": "password",
	}

	body, err := json.Marshal(request)
	if err != nil {
		assert.FailNow(t, "fail to marshal json", err)
	}

	resp, err := http.Post(testServer.URL+"/signup", "application/json", bytes.NewBuffer(body))
	if err != nil {
		assert.FailNow(t, "fail to request", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			assert.Fail(t, "failed to close body", err)
		}
	}()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp2, err := http.Post(testServer.URL+"/signup", "application/json", bytes.NewBuffer(body))
	if err != nil {
		assert.FailNow(t, "fail to request", err)
	}

	defer func() {
		err := resp2.Body.Close()
		if err != nil {
			assert.Fail(t, "failed to close body", err)
		}
	}()

	assert.Equal(t, http.StatusInternalServerError, resp2.StatusCode)

	err = testDb.Cleanup()
	if err != nil {
		assert.Fail(t, "fail to cleanup testDb")
	}
}
