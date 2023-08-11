package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yervsil/toDo-microservice/internal/entity"
	"github.com/yervsil/toDo-microservice/internal/service"
	service_mocks "github.com/yervsil/toDo-microservice/internal/service/mocks"
	"github.com/yervsil/toDo-microservice/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_createTask(t *testing.T) {
	taskID := primitive.NewObjectID()

	type mockBehavior func(r *service_mocks.MockTask, ctx context.Context, task entity.Task)

	tests := []struct {
		name                 string
		inputBody            string
		inputTask            entity.Task
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"title":"Купить книгу", "activeAt":"2023-08-04"}`,
			inputTask: entity.Task{
				Title:    "Купить книгу",
				ActiveAt: "2023-08-04",
			},
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, task entity.Task) {
				r.EXPECT().CreateTask(ctx, task).Return(taskID, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`{"id":"%s"}`, taskID.Hex()),
		},

		{
			name:                 "InvalidInput",
			inputBody:            `{"invalid_field":"value"}`, 
			mockBehavior:         func(r *service_mocks.MockTask, ctx context.Context, task entity.Task) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},

		{
			name:                 "InvalidDateFormat",
			inputBody:            `{"title":"Купить книгу", "activeAt":"invalid_date"}`, // Некорректный формат даты
			mockBehavior:         func(r *service_mocks.MockTask, ctx context.Context, task entity.Task) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect date format"}`,
		},

		{
			name:      "ServiceError",
			inputBody: `{"title":"Купить книгу", "activeAt":"2023-08-04"}`,
			inputTask: entity.Task{
				Title:    "Купить книгу",
				ActiveAt: "2023-08-04",
			},
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, task entity.Task) {
				r.EXPECT().CreateTask(ctx, task).Return(primitive.ObjectID{}, errors.New("service error"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockTask(c)
			test.mockBehavior(repo, context.Background(), test.inputTask)

			services := &service.Service{Task: repo}
			handler := Handler{services, logger.New("local")}

			// Init Endpoint
			r := gin.New()
			r.POST("/tasks", handler.createTask)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/tasks",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_updateTask(t *testing.T) {

	type mockBehavior func(r *service_mocks.MockTask, ctx context.Context, task entity.Task, taskID primitive.ObjectID)

	tests := []struct {
		name                 string
		inputBody            string
		inputTask            entity.Task
		taskID               string
		ctx                  *gin.Context
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"title":"Купить книгу - Высоконагруженные приложения", "activeAt":"2023-08-05"}`,
			inputTask: entity.Task{
				Title:    "Купить книгу - Высоконагруженные приложения",
				ActiveAt: "2023-08-05",
			},
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, task entity.Task, taskID primitive.ObjectID) {
				r.EXPECT().UpdateTask(ctx, task, taskID).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `"successfully updated"`,
		},
		{
			name:      "InvalidInput",
			inputBody: `{}`, // Invalid input body
			inputTask: entity.Task{},
			taskID:    "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, task entity.Task, taskID primitive.ObjectID) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "incorrectIDParam",
			inputBody: `{"title":"Купить книгу - Высоконагруженные приложения", "activeAt":"2023-08-05"}`,
			inputTask: entity.Task{},
			taskID:    "64d1c8747124f",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, task entity.Task, taskID primitive.ObjectID) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id param"}`,
		},
		{
			name:      "InvalidDateFormat",
			inputBody: `{"title":"Купить книгу - Высоконагруженные приложения", "activeAt":"2023/08/05"}`,
			inputTask: entity.Task{
				Title:    "Купить книгу - Высоконагруженные приложения",
				ActiveAt: "2023/08/05",
			},
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, task entity.Task, taskID primitive.ObjectID) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"incorrect date format"}`,
		},
		{
			name:      "ServiceError",
			inputBody: `{"title":"Купить книгу - Высоконагруженные приложения", "activeAt":"2023-08-05"}`,
			inputTask: entity.Task{
				Title:    "Купить книгу - Высоконагруженные приложения",
				ActiveAt: "2023-08-05",
			},
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, task entity.Task, taskID primitive.ObjectID) {
				r.EXPECT().UpdateTask(ctx, task, taskID).Return(errors.New("update error"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"error":"update error"}`,
		},
	}

	for _, test := range tests {
		id, _ := primitive.ObjectIDFromHex(test.taskID)
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockTask(c)
			test.mockBehavior(repo, context.Background(), test.inputTask, id)

			services := &service.Service{Task: repo}
			handler := Handler{services, logger.New("local")}

			// Init Endpoint
			r := gin.New()
			r.PUT("/tasks/:id", handler.updateTask)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/tasks/%s", test.taskID),
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteTask(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID)

	tests := []struct {
		name                 string
		taskID               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok",
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID) {
				r.EXPECT().DeleteTask(ctx, taskID).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `"successfully deleted"`,
		},
		{
			name:                 "InvalidIDParam",
			taskID:               "64d1c8747124f40af8030b", // Invalid taskID
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id param"}`,
		},
		{
			name:   "NotFound",
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID) {
				r.EXPECT().DeleteTask(ctx, taskID).Return(errors.New("task not found"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"error":"task not found"}`,
		},
		{
			name:   "InternalServerError",
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID) {
				r.EXPECT().DeleteTask(ctx, taskID).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		id, _ := primitive.ObjectIDFromHex(test.taskID)
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockTask(c)
			test.mockBehavior(repo, context.Background(), id)

			services := &service.Service{Task: repo}
			handler := Handler{services, logger.New("local")}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/tasks/:id", handler.deleteTask)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/tasks/%s", test.taskID), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_statusUpdate(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID)

	tests := []struct {
		name                 string
		taskID               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok",
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID) {
				r.EXPECT().StatusUpdate(ctx, taskID).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `"status has been changed"`,
		},
		{
			name:                 "InvalidIDParam",
			taskID:               "64d1c8747124f40af8030b", // Invalid taskID
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid id param"}`,
		},
		{
			name:   "NotFound",
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID) {
				r.EXPECT().StatusUpdate(ctx, taskID).Return(errors.New("task not found"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"error":"task not found"}`,
		},
		{
			name:   "InternalServerError",
			taskID: "64d1c8747124f40af803840b",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, taskID primitive.ObjectID) {
				r.EXPECT().StatusUpdate(ctx, taskID).Return(errors.New("internal server error"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockTask(c)
			ctx := context.Background()
			taskID, _ := primitive.ObjectIDFromHex(test.taskID)
			
			test.mockBehavior(repo, ctx, taskID)

			services := &service.Service{Task: repo}
			handler := Handler{services, logger.New("local")}

			// Init Endpoint
			r := gin.New()
			r.PUT("/tasks/:id/done", handler.statusUpdate)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/tasks/%s/done", test.taskID), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}


func TestHandler_getTasks(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockTask, ctx context.Context, status string)

	tests := []struct {
		name                 string
		queryStatus          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "ActiveStatus_NoTasks",
			queryStatus:  "active",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, status string) {
				r.EXPECT().GetTasks(ctx, status).Return([]entity.Task{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name:         "CompletedStatus_NoTasks",
			queryStatus:  "done",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, status string) {
				r.EXPECT().GetTasks(ctx, status).Return([]entity.Task{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name:         "ActiveStatus_WithTasks",
			queryStatus:  "active",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, status string) {
				tasks := []entity.Task{
					{Title: "Task 1", ActiveAt: "2023-08-10"},
					{Title: "Task 2", ActiveAt: "2023-08-11"},
				}
				r.EXPECT().GetTasks(ctx, status).Return(tasks, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"title":"Task 1","activeAt":"2023-08-10"},{"title":"Task 2","activeAt":"2023-08-11"}]`,
		},
		{
			name:         "InvalidStatus",
			queryStatus:  "invalid_status",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, status string) {
				r.EXPECT().GetTasks(ctx, status).Return([]entity.Task{}, errors.New("invalid status parameter"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"error":"invalid status parameter"}`,
		},
		{
			name:         "InternalServerError",
			queryStatus:  "active",
			mockBehavior: func(r *service_mocks.MockTask, ctx context.Context, status string) {
				r.EXPECT().GetTasks(ctx, status).Return(nil, errors.New("internal server error"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockTask(c)
			ctx := context.Background()
			test.mockBehavior(repo, ctx, test.queryStatus)

			services := &service.Service{Task: repo}
			handler := Handler{services, logger.New("local")}

			// Init Endpoint
			r := gin.New()
			r.GET("/tasks", handler.getTasks)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/tasks?status=%s", test.queryStatus), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}