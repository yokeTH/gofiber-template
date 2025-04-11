package server_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yokeTH/gofiber-template/internal/server"
)

func TestServer(t *testing.T) {
	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "health check",
			route:         "/health",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  `{"status":"ok"}`,
		},
		{
			description:   "not found",
			route:         "/not-found",
			expectedError: false,
			expectedCode:  404,
			expectedBody:  `{"error":"Cannot GET /not-found"}`,
		},
	}

	s := server.New(
		server.WithName("MOCK SERVER"),
	)

	for _, test := range tests {
		req, _ := http.NewRequest("GET", test.route, nil)
		res, err := s.Test(req, -1)
		assert.Equalf(t, test.expectedError, err != nil, test.description)
		if test.expectedError {
			continue
		}
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
		body, err := io.ReadAll(res.Body)
		assert.Nilf(t, err, test.description)
		assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}
