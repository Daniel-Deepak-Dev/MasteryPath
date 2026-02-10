package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContributionDay struct {
	Date  string `json:"date" bson:"_id"`
	Count int    `json:"count" bson:"count"`
	Level int    `json:"level" bson:"level"` // 0-4 scale for coloring
}

type AnalyticsService interface {
	GetContributionData(ctx context.Context) ([]ContributionDay, error)
}

type analyticsService struct {
	progressCol *mongo.Collection
}

func NewAnalyticsService(progressCol *mongo.Collection) AnalyticsService {
	return &analyticsService{
		progressCol: progressCol,
	}
}

func (s *analyticsService) GetContributionData(ctx context.Context) ([]ContributionDay, error) {
	// 1. Calculate the date 365 days ago
	oneYearAgo := time.Now().AddDate(0, 0, -365)

	// 2. Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match items completed in the last year
		{{Key: "$match", Value: bson.D{
			{Key: "completed_at", Value: bson.D{{Key: "$gte", Value: oneYearAgo}}},
		}}},
		// Group by date (YYYY-MM-DD)
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "format", Value: "%Y-%m-%d"},
					{Key: "date", Value: "$completed_at"},
				}},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		// Sort by date ascending
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	// 3. Execute the pipeline
	cursor, err := s.progressCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 4. Decode results
	var results []ContributionDay
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// 5. Calculate levels (0-4) based on count
	// Simple logic for now: 0=0, 1=1-2, 2=3-4, 3=5-6, 4=7+
	for i := range results {
		count := results[i].Count
		if count == 0 {
			results[i].Level = 0
		} else if count <= 2 {
			results[i].Level = 1
		} else if count <= 4 {
			results[i].Level = 2
		} else if count <= 6 {
			results[i].Level = 3
		} else {
			results[i].Level = 4
		}
	}

	// Make sure we return an empty slice instead of nil if no data
	if results == nil {
		results = []ContributionDay{}
	}

	return results, nil
}
