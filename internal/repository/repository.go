package repository

import (
	"context"

	"github.com/yervsil/toDo-microservice/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Task interface {
	CreateTask(ctx context.Context, task entity.Task) (primitive.ObjectID, error)
	UpdateTask(ctx context.Context, task entity.Task, taskId primitive.ObjectID) error
	DeleteTask(ctx context.Context, taskId primitive.ObjectID) error
	StatusUpdate(ctx context.Context, taskId primitive.ObjectID) error
	GetTasks(ctx context.Context, status string) ([]entity.Task, error)
}

type Repository struct {
	Task
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{Task: NewTaskRepoistory(db)}
}