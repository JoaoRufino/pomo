package mongodb

import (
	"context"
	"fmt"

	"github.com/joaorufino/pomo/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDBStore struct {
	client   *mongo.Client
	database *mongo.Database
	logger   *zap.SugaredLogger
}

func NewStore(uri, databaseName string, logger *zap.SugaredLogger) (*MongoDBStore, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Error("failed to connect to MongoDB", zap.Error(err))
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.Error("failed to ping MongoDB", zap.Error(err))
		return nil, err
	}

	database := client.Database(databaseName)

	return &MongoDBStore{
		client:   client,
		database: database,
		logger:   logger,
	}, nil
}

func (s *MongoDBStore) InitDB() error {
	// Ensure indexes or initial setup if necessary
	return nil
}

func (s *MongoDBStore) TaskSave(ctx context.Context, task *models.Task) (int, error) {
	collection := s.database.Collection("tasks")
	filter := bson.M{"id": task.ID}

	update := bson.M{
		"$set": task,
	}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		s.logger.Error("failed to save task", zap.Error(err))
		return 0, err
	}
	s.logger.Info("task saved successfully", zap.Int("taskID", task.ID))
	return task.ID, nil
}

func (s *MongoDBStore) GetAllTasks(ctx context.Context) (models.List, error) {
	collection := s.database.Collection("tasks")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		s.logger.Error("failed to get all tasks", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks models.List
	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			s.logger.Error("failed to decode task", zap.Error(err))
			continue
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("error occurred during cursor iteration", zap.Error(err))
		return nil, err
	}

	return tasks, nil
}

func (s *MongoDBStore) TaskDeleteByID(ctx context.Context, taskID int) error {
	collection := s.database.Collection("tasks")
	filter := bson.M{"id": taskID}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		s.logger.Error("failed to delete task by ID", zap.Int("taskID", taskID), zap.Error(err))
		return err
	}
	s.logger.Info("task deleted successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *MongoDBStore) TaskGetByID(ctx context.Context, taskID int) (*models.Task, error) {
	collection := s.database.Collection("tasks")
	filter := bson.M{"id": taskID}

	var task models.Task
	err := collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrNotFound
		}
		s.logger.Error("failed to get task by ID", zap.Int("taskID", taskID), zap.Error(err))
		return nil, err
	}

	return &task, nil
}

func (s *MongoDBStore) PomodoroSave(ctx context.Context, taskID int, pomodoro *models.Pomodoro) error {
	collection := s.database.Collection("pomodoros")
	filter := bson.M{"id": pomodoro.ID}

	update := bson.M{
		"$set": pomodoro,
	}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		s.logger.Error("failed to save pomodoro", zap.Error(err))
		return err
	}
	s.logger.Info("pomodoro saved successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *MongoDBStore) PomodoroGetByTaskID(ctx context.Context, taskID int) ([]models.Pomodoro, error) {
	collection := s.database.Collection("pomodoros")
	filter := bson.M{"task_id": taskID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get pomodoros by task ID", zap.Int("taskID", taskID), zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var pomodoros []models.Pomodoro
	for cursor.Next(ctx) {
		var pomodoro models.Pomodoro
		if err := cursor.Decode(&pomodoro); err != nil {
			s.logger.Error("failed to decode pomodoro", zap.Error(err))
			continue
		}
		pomodoros = append(pomodoros, pomodoro)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("error occurred during cursor iteration", zap.Error(err))
		return nil, err
	}

	return pomodoros, nil
}

func (s *MongoDBStore) PomodoroDeleteByTaskID(ctx context.Context, taskID int) error {
	collection := s.database.Collection("pomodoros")
	filter := bson.M{"task_id": taskID}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		s.logger.Error("failed to delete pomodoros by task ID", zap.Int("taskID", taskID), zap.Error(err))
		return err
	}
	s.logger.Info("pomodoros deleted successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *MongoDBStore) Close() error {
	return s.client.Disconnect(context.Background())
}

func getTaskKey(taskID int) string {
	return fmt.Sprintf("task:%d", taskID)
}

func getPomodoroKey(taskID, pomodoroID int) string {
	return fmt.Sprintf("pomodoro:%d:%d", taskID, pomodoroID)
}
