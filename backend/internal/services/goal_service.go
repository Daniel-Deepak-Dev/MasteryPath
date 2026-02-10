package services

import (
	"context"
	"time"

	"masterypath/internal/apperror"
	"masterypath/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GoalService handles business logic for Goal operations.
type GoalService struct {
	collection *mongo.Collection
}

// NewGoalService creates a new GoalService with the given collection.
func NewGoalService(col *mongo.Collection) *GoalService {
	return &GoalService{collection: col}
}

// GetAll retrieves all goals from the database.
func (s *GoalService) GetAll(ctx context.Context) ([]models.Goal, error) {
	var goals []models.Goal
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to fetch goals")
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &goals); err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to decode goals")
	}

	// Return empty slice, not nil
	if goals == nil {
		goals = []models.Goal{}
	}
	return goals, nil
}

// GetByID retrieves a single goal by its ObjectID.
func (s *GoalService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Goal, error) {
	var goal models.Goal
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&goal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to fetch goal")
	}
	return &goal, nil
}

// Create inserts a new goal into the database.
func (s *GoalService) Create(ctx context.Context, goal *models.Goal) error {
	goal.ID = primitive.NewObjectID()
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()

	_, err := s.collection.InsertOne(ctx, goal)
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to insert goal")
	}
	return nil
}

// Update modifies an existing goal identified by its ObjectID.
func (s *GoalService) Update(ctx context.Context, id primitive.ObjectID, goal *models.Goal) error {
	update := bson.M{
		"$set": bson.M{
			"title":       goal.Title,
			"description": goal.Description,
			"completed":   goal.Completed,
			"updated_at":  time.Now(),
		},
	}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to update goal")
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	goal.ID = id
	return nil
}

// Delete removes a goal from the database by its ObjectID.
func (s *GoalService) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to delete goal")
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
