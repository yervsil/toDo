package service

import (
	"context"

	"github.com/yervsil/toDo-microservice/internal/entity"
	"github.com/yervsil/toDo-microservice/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Task interface {
	CreateTask(ctx context.Context, input entity.Task) (primitive.ObjectID, error)
	UpdateTask(ctx context.Context, input entity.Task, taskId primitive.ObjectID) error
	DeleteTask(ctx context.Context, taskId primitive.ObjectID) error
	StatusUpdate(ctx context.Context, taskId primitive.ObjectID) error
	GetTasks(ctx context.Context, status string) ([]entity.Task, error)
}

type Service struct {
	Task
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Task: NewTaskService(repo)}
}