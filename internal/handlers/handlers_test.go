package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"students-crud/internal/handlers"
	mock_handlers "students-crud/internal/handlers/mock"
	"students-crud/internal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandlers_CreateStudent(t *testing.T) {
	type mockBehavior func(s *mock_handlers.MockStorage, student *models.Student)

	testCases := []struct {
		name                string
		inputBody           string
		inputStudent        models.Student
		mockBehaviour       mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Student #1","email": "#1@mail.com"}`,
			inputStudent: models.Student{
				Name:  "Student #1",
				Email: "#1@mail.com",
			},
			mockBehaviour: func(s *mock_handlers.MockStorage, student *models.Student) {
				s.EXPECT().Create(gomock.Any(), student).Return(1, nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:      "Error on Create",
			inputBody: `{"name": "Student #2","email": "#2@mail.com"}`,
			inputStudent: models.Student{
				Name:  "Student #2",
				Email: "#2@mail.com",
			},
			mockBehaviour: func(s *mock_handlers.MockStorage, student *models.Student) {
				s.EXPECT().Create(gomock.Any(), student).Return(0, errors.New("failed to create student"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"failed to create student"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			storage := mock_handlers.NewMockStorage(c)
			testCase.mockBehaviour(storage, &testCase.inputStudent)

			handlers := handlers.NewHandlers(storage)

			r := gin.Default()
			r.POST("/students", handlers.CreateStudent)

			req, _ := http.NewRequest(http.MethodPost, "/students", bytes.NewBufferString(testCase.inputBody))
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.Equal(t, testCase.expectedRequestBody, rec.Body.String())
		})

	}
}

func TestHandlers_ReadStudent(t *testing.T) {
	type mockBehavior func(s *mock_handlers.MockStorage, id int)

	testCases := []struct {
		name                string
		inputID             int
		mockBehaviour       mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:    "OK",
			inputID: 1,
			mockBehaviour: func(s *mock_handlers.MockStorage, id int) {
				s.EXPECT().Read(gomock.Any(), id).Return(&models.Student{
					ID:    id,
					Name:  "Student #1",
					Email: "#1@mail.com",
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"name":"Student #1","email":"#1@mail.com"}`,
		},
		{
			name:    "Not Found",
			inputID: 2,
			mockBehaviour: func(s *mock_handlers.MockStorage, id int) {
				s.EXPECT().Read(gomock.Any(), id).Return(nil, errors.New("student not found")) // Возвращаем ошибку
			},
			expectedStatusCode:  404,
			expectedRequestBody: `{"error":"student not found"}`,
		},
		{
			name:    "Invalid ID",
			inputID: -1, // Неверный ID для проверки
			mockBehaviour: func(s *mock_handlers.MockStorage, id int) {
				// Никаких вызовов не требуется для Invalid ID
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid id"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			storage := mock_handlers.NewMockStorage(c)
			if testCase.name != "Invalid ID" {
				testCase.mockBehaviour(storage, testCase.inputID)
			}

			handlers := handlers.NewHandlers(storage)

			r := gin.Default()
			r.GET("/students/:id", handlers.ReadStudent)

			var req *http.Request
			if testCase.name == "Invalid ID" {
				req, _ = http.NewRequest(http.MethodGet, "/students/-1", nil) // Для проверки неверного ID
			} else {
				req, _ = http.NewRequest(http.MethodGet, "/students/"+strconv.Itoa(testCase.inputID), nil)
			}

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.Equal(t, testCase.expectedRequestBody, rec.Body.String())
		})
	}
}
