package services

import (
	"context"

	"masterypath/internal/apperror"
	"masterypath/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProgressService handles business logic for ProgressItem operations.
type ProgressService struct {
	collection *mongo.Collection
}

// NewProgressService creates a new ProgressService with the given collection.
func NewProgressService(col *mongo.Collection) *ProgressService {
	return &ProgressService{collection: col}
}

// Create inserts a new progress item after validating its weightage.
func (s *ProgressService) Create(ctx context.Context, item *models.ProgressItem) error {
	if err := item.ValidateWeightage(); err != nil {
		return apperror.New(400, err.Error())
	}

	item.ID = primitive.NewObjectID()

	_, err := s.collection.InsertOne(ctx, item)
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to insert progress item")
	}
	return nil
}

// GetBySkill retrieves all progress items for a skill with computed weight percentages.
// Uses a MongoDB aggregation pipeline with $setWindowFields for server-side percentage calculation.
func (s *ProgressService) GetBySkill(ctx context.Context, skillID primitive.ObjectID) ([]models.ProgressItem, error) {
	pipeline := mongo.Pipeline{
		// 1. Match items for this skill
		{{Key: "$match", Value: bson.D{{Key: "parent_skill_id", Value: skillID}}}},

		// 2. Calculate Total Weight using $setWindowFields
		{{Key: "$setWindowFields", Value: bson.D{
			{Key: "partitionBy", Value: nil},
			{Key: "output", Value: bson.D{
				{Key: "total_weight", Value: bson.D{{Key: "$sum", Value: "$weightage"}}},
			}},
		}}},

		// 3. Calculate Percentage
		{{Key: "$addFields", Value: bson.D{
			{Key: "weight_percent", Value: bson.D{
				{Key: "$multiply", Value: bson.A{
					bson.D{{Key: "$divide", Value: bson.A{"$weightage", "$total_weight"}}},
					100,
				}},
			}},
		}}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to fetch progress items")
	}
	defer cursor.Close(ctx)

	var items []models.ProgressItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to decode progress items")
	}

	if items == nil {
		items = []models.ProgressItem{}
	}
	return items, nil
}

// Update modifies an existing progress item identified by its ObjectID.
func (s *ProgressService) Update(ctx context.Context, id primitive.ObjectID, updateData *models.ProgressItem) error {
	// Validate weightage if being updated
	if updateData.Weightage != 0 {
		if err := updateData.ValidateWeightage(); err != nil {
			return apperror.New(400, err.Error())
		}
	}

	update := bson.M{
		"$set": bson.M{
			"name":      updateData.Name,
			"weightage": updateData.Weightage,
			"achieved":  updateData.Achieved,
			"comments":  updateData.Comments,
		},
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to update progress item")
	}
	return nil
}

// Delete removes a progress item from the database by its ObjectID.
func (s *ProgressService) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to delete progress item")
	}
	return nil
}
