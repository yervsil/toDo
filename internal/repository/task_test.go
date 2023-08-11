package repository

import (

	"context"
	"errors"

	"testing"


	"github.com/stretchr/testify/assert"
	"github.com/yervsil/toDo-microservice/internal/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	newTask := entity.Task{
		Title:    "New Task",
		ActiveAt: "2023-08-15",
	}

	mt.Run("success", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "test.task", mtest.FirstBatch)
		killCursors := mtest.CreateCursorResponse(0, "test.task", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		repo := &taskRepository{
			db: mt.Coll,
		}

		insertedID, err := repo.CreateTask(context.Background(), newTask)
		assert.Nil(t, err)
		assert.NotEqual(t, primitive.ObjectID{}, insertedID)
	})

	mt.Run("duplicate_document", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse())
		repo := &taskRepository{
			db: mt.Coll,
		}

		_, err := repo.CreateTask(context.Background(), newTask)
		assert.Equal(t, "this document already exists", err.Error())
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse())
		repo := &taskRepository{
			db: mt.Coll,
		}

		_, err := repo.CreateTask(context.Background(), newTask)
		assert.NotNil(t, err)
	})
}


func TestUpdateTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	taskID := primitive.NewObjectID()
	taskToUpdate := entity.Task{
		Title:    "Updated Title",
		ActiveAt: "2023-08-10",
	}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		repo := &taskRepository{
			db: mt.Coll,
		}

		err := repo.UpdateTask(context.Background(), taskToUpdate, taskID)
		assert.Nil(t, err)
	})

	mt.Run("no_record_found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		repo := &taskRepository{
			db: mt.Coll,
		}

		err := repo.UpdateTask(context.Background(), taskToUpdate, taskID)
		assert.Equal(t, "no record found", err.Error())
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse())
		repo := &taskRepository{
			db: mt.Coll,
		}

		err := repo.UpdateTask(context.Background(), taskToUpdate, taskID)
		assert.NotNil(t, err)
	})
}

func TestDeleteTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	taskID := primitive.NewObjectID()

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		repo := &taskRepository{
			db: mt.Coll,
		}

		err := repo.DeleteTask(context.Background(), taskID)
		assert.Nil(t, err)
	})

	mt.Run("no_record_found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		repo := &taskRepository{
			db: mt.Coll,
		}

		err := repo.DeleteTask(context.Background(), taskID)
		assert.Equal(t, "no record found", err.Error())
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse())
		repo := &taskRepository{
			db: mt.Coll,
		}

		err := repo.DeleteTask(context.Background(), taskID)
		assert.NotNil(t, err)
	})
}

func TestStatusUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()


	taskID := primitive.NewObjectID()

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
			
		repo := &taskRepository{
			db: mt.Coll,
		}
	
		err := repo.StatusUpdate(context.Background(), taskID)
	
		assert.Equal(t, nil, err)
	})

	mt.Run("no_record_found", func(mt *mtest.T) {
		

		mt.AddMockResponses(mtest.CreateSuccessResponse())
			
		repo := &taskRepository{
			db: mt.Coll,
		}
	
		err := repo.StatusUpdate(context.Background(), taskID)

		assert.Equal(t, "no record found", err.Error())
	})
}

func TestGetTasks(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		want := []entity.Task{
			{Title: "купить telephone", ActiveAt: "2022-07-30"},
			{Title: "купить iphone", ActiveAt: "2023-07-30"},
		}

		tr := &taskRepository{
			db: mt.Coll,
		}
	
		first := mtest.CreateCursorResponse(1, "test.task", mtest.FirstBatch, bson.D{
			{Key: "title", Value: "купить telephone"},
			{Key: "activeat", Value: "2022-07-30"},
		})

		second := mtest.CreateCursorResponse(1, "test.task", mtest.NextBatch, bson.D{
			{Key: "title", Value: "купить iphone"},
			{Key: "activeat", Value: "2023-07-30"},
		})
		
		killCursors := mtest.CreateCursorResponse(0, "test.task", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)
		status := "done"

		got, err := tr.GetTasks(context.Background(), status)
		if err != nil {
			t.Fatalf("expected: no error, got: %v", err)
		}
		
		if !assert.Equal(t, got, want){
			t.Fatalf("expected: %v, got: %v", want, got)
		}
	})

	mt.Run("error", func(mt *mtest.T) {
		want := errors.New("incorrect url query")
		
		tr := &taskRepository{
			db: mt.Coll,
		}

		status := "Active"

		_, got := tr.GetTasks(context.Background(), status)
		
		if !assert.Equal(t, got, want){
			t.Fatalf("expected: %v, got: %v", want, got)
		}
	})
}