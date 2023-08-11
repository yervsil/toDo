package service

import (
	"context"
	"time"

	"github.com/yervsil/toDo-microservice/internal/entity"
	"github.com/yervsil/toDo-microservice/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	active = "active"
)

type TaskService struct {
	repo *repository.Repository
}

func NewTaskService(repo *repository.Repository) *TaskService {
	return &TaskService{repo: repo}
}

func(t *TaskService) CreateTask(ctx context.Context, task entity.Task) (primitive.ObjectID, error){
	task.Status = active
	return t.repo.CreateTask(ctx, task)
}

func(t *TaskService) UpdateTask(ctx context.Context, task entity.Task, taskId primitive.ObjectID) error{
	task.Status = active
	return t.repo.UpdateTask(ctx, task, taskId)
}

func(t *TaskService) DeleteTask(ctx context.Context, taskId primitive.ObjectID) error{
	return t.repo.DeleteTask(ctx, taskId)
}

func(t *TaskService) StatusUpdate(ctx context.Context, taskId primitive.ObjectID) error{
	return t.repo.StatusUpdate(ctx, taskId)
}

func(t *TaskService) GetTasks(ctx context.Context, status string) ([]entity.Task, error){
	tasks, err := t.repo.GetTasks(ctx, status)
    if err != nil {
        return nil, err
    }

    for i, task := range tasks {
        activeDate, err := time.Parse("2006-01-02", task.ActiveAt)
        if err != nil {
            return nil, err
        }

        if activeDate.Weekday() == time.Saturday || activeDate.Weekday() == time.Sunday {
            tasks[i].Title = "ВЫХОДНОЙ - " + tasks[i].Title
        }
    }

    return tasks, nil
}