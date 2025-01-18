package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-starter/cmd/server/router"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserHandler_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc         string
		mockResponse http.HandlerFunc
		wantStatus   int
		wantBody     string
	}{
		{
			desc: "success",
			mockResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`[{
					"id": 1,
					"name": "Test User",
					"username": "testuser",
					"email": "test@example.com",
					"address": {
						"street": "Test St",
						"suite": "Apt 1",
						"city": "Test City",
						"zipcode": "12345"
					},
					"phone": "1-234-567-8900",
					"website": "test.com"
				}]`))
				assert.NoError(t, err)
			},
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":1,"name":"Test User","username":"testuser","email":"test@example.com","address":{"street":"Test St","suite":"Apt 1","city":"Test City","zipcode":"12345"},"phone":"1-234-567-8900","website":"test.com"}]`,
		},
		{
			desc: "external api returns error",
			mockResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"error": "origin server error"}`))
				assert.NoError(t, err)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"failed to fetch data from external API"}`,
		},
		{
			desc: "external api returns invalid json",
			mockResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`invalid json`))
				assert.NoError(t, err)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"failed to decode response"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			// Create a mock external API server
			mockServer := httptest.NewServer(tc.mockResponse)
			defer mockServer.Close()

			h := router.NewUserHandler(&http.Client{}, mockServer.URL)
			r := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			w := httptest.NewRecorder()

			// When:
			h.Get(r).Respond(w, r)

			got := w.Result()
			defer got.Body.Close()
			gotBodyBytes, err := io.ReadAll(got.Body)
			require.NoError(t, err)

			// Then:
			assert.Equal(t, tc.wantStatus, w.Code)
			require.JSONEq(t, tc.wantBody, string(gotBodyBytes))
		})
	}
}
