package handler_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/yokeTH/gofiber-template/internal/core/domain"
	"github.com/yokeTH/gofiber-template/internal/core/service"
	"github.com/yokeTH/gofiber-template/internal/database"
	"github.com/yokeTH/gofiber-template/internal/handler"
	"github.com/yokeTH/gofiber-template/internal/repository"
	"github.com/yokeTH/gofiber-template/internal/server"
	"github.com/yokeTH/gofiber-template/pkg/mock"
)

func TestBook(t *testing.T) {
	tests := []struct {
		description  string
		method       string
		route        string
		body         any
		expectedCode int
		expectedBody map[string]any
	}{
		{
			description:  "get empty books",
			route:        "/books",
			method:       "GET",
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": []any{},
				"pagination": map[string]any{
					"current_page": float64(1),
					"last_page":    float64(0),
					"limit":        float64(10),
					"total":        float64(0),
				},
			},
		},
		{
			description:  "get not found book",
			route:        "/books/1",
			method:       "GET",
			expectedCode: 404,
			expectedBody: map[string]any{
				"error": "book not found",
			},
		},
		{
			description:  "create book",
			route:        "/books",
			method:       "POST",
			body:         domain.Book{Title: "test", Author: "test"},
			expectedCode: 201,
			expectedBody: map[string]any{
				"data": map[string]any{
					"id":     float64(1),
					"title":  "test",
					"author": "test",
				},
			},
		},
		{
			description:  "update book",
			route:        "/books/1",
			method:       "PATCH",
			body:         domain.Book{Title: "updated title", Author: "updated author"},
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": map[string]any{
					"id":     float64(1),
					"title":  "updated title",
					"author": "updated author",
				},
			},
		},
		{
			description:  "update book not found",
			route:        "/books/999",
			method:       "PATCH",
			body:         domain.Book{Title: "updated title", Author: "updated author"},
			expectedCode: 404,
			expectedBody: map[string]any{
				"error": "book not found",
			},
		},
		{
			description:  "delete book",
			route:        "/books/1",
			method:       "DELETE",
			expectedCode: 204,
			expectedBody: map[string]any(nil),
		},
	}

	db, err := mock.SetupMockDB()
	assert.Nil(t, err)

	err = db.AutoMigrate(&domain.Book{})
	assert.Nil(t, err)

	bookRepository := repository.NewBookRepository(&database.Database{DB: db})
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	s := server.New(
		server.WithName("MOCK SERVER"),
	)

	s.Get("/books", bookHandler.GetBooks)
	s.Get("/books/:id", bookHandler.GetBook)
	s.Post("/books", bookHandler.CreateBook)
	s.Patch("/books/:id", bookHandler.UpdateBook)
	s.Delete("/books/:id", bookHandler.DeleteBook)

	for _, test := range tests {
		var body io.Reader
		if test.body != nil {
			jsonBody, err := json.Marshal(test.body)
			assert.Nilf(t, err, test.description)
			body = bytes.NewReader(jsonBody)
		}

		req, _ := http.NewRequest(test.method, test.route, body)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.Test(req, -1)
		assert.Nilf(t, err, test.description)

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		resBody, err := io.ReadAll(res.Body)
		assert.Nilf(t, err, test.description)

		if test.expectedBody == nil {
			continue
		}

		var actual map[string]any
		err = json.Unmarshal(resBody, &actual)
		assert.Nilf(t, err, test.description)

		assert.Equalf(t, test.expectedBody, actual, test.description)
	}

	mock.CleanupMockDB()
}
