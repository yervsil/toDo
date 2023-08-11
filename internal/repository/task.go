package repository

import (
	"context"
	"errors"
	"time"

	"github.com/yervsil/toDo-microservice/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	done = "done"
	active = "active"
)

type taskRepository struct {
	db *mongo.Collection
}

func NewTaskRepoistory(db *mongo.Database) *taskRepository {
	return &taskRepository{db: db.Collection(tasksCollection)}
}

func (r *taskRepository) CreateTask(ctx context.Context, task entity.Task) (primitive.ObjectID, error){
	if isDuplicate(task, r.db){
		return primitive.ObjectID{}, errors.New("this document already exists")
	}
	res, err := r.db.InsertOne(ctx, task)
	
	return res.InsertedID.(primitive.ObjectID), err
}

func (r *taskRepository) UpdateTask(ctx context.Context, task entity.Task, taskId primitive.ObjectID) error{
	filter := bson.M{"_id": taskId}
	res, err := r.db.ReplaceOne(context.Background(), filter, task)
	if res.MatchedCount == 0 {
		return errors.New("no record found")
	}

	return err
}

func (r *taskRepository) DeleteTask(ctx context.Context, taskId primitive.ObjectID) error{
	filter := bson.M{"_id": taskId}
	res, err := r.db.DeleteOne(context.Background(), filter)
    if res.DeletedCount == 0 {
		return errors.New("no record found")
	}

	return err
}

func (r *taskRepository) StatusUpdate(ctx context.Context, taskId primitive.ObjectID) error{
	update := bson.M{"$set": bson.M{"status": done}}

	filter := bson.M{"_id": taskId}

	res, err := r.db.UpdateOne(ctx, filter, update)
	if res.MatchedCount == 0 {
		return errors.New("no record found")
	}

	return err
}

func (r *taskRepository) GetTasks(ctx context.Context, status string) ([]entity.Task, error){
	var filter primitive.M

	if status == active{
		filter = bson.M{
			"status":   status,
			"activeat": bson.M{"$lte": time.Now().Format("2006-01-02")},
		}
	}else if status == done{
		filter = bson.M{
			"status":   status,
		}
	}else{
		return nil, errors.New("incorrect url query")
	}

	projection := options.Find().SetProjection(bson.M{
		"status": 0,
	})

	sortOptions := options.Find().SetSort(bson.D{{"activeat", 1}})

	cursor, err := r.db.Find(ctx, filter, projection, sortOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []entity.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func isDuplicate(task entity.Task, collection *mongo.Collection) bool {
	filter := bson.M{
		"title":     task.Title,
		"activeat":  task.ActiveAt,
	}

	var existingTask entity.Task

	err := collection.FindOne(context.Background(), filter).Decode(&existingTask)
	if err == mongo.ErrNoDocuments {
		return false
	} else if err != nil {
		return true
	}
	
	return true
}