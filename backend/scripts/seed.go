package main

import (
	"context"
	"log"
	"masterypath/internal/config"
	"masterypath/internal/database"
	"masterypath/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	cfg := config.LoadConfig()
	client := database.ConnectDB(cfg.MongoURI)
	defer client.Disconnect(context.Background())

	db := client.Database("masterypath")
	skillCol := db.Collection("skills")
	progressCol := db.Collection("progress")

	// Clear existing data
	_, _ = skillCol.DeleteMany(context.Background(), bson.M{})
	_, _ = progressCol.DeleteMany(context.Background(), bson.M{})

	log.Println("Cleared existing skills and progress data")

	// 1. Root Skill: Web Development
	webDevID := primitive.NewObjectID()
	webDev := models.Skill{
		ID:          webDevID,
		Name:        "Web Development",
		Category:    "Technical",
		Description: "The art of building websites and web applications.",
		Ancestors:   []primitive.ObjectID{},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), webDev)

	// 2. Child Skill: Frontend
	frontendID := primitive.NewObjectID()
	frontend := models.Skill{
		ID:          frontendID,
		Name:        "Frontend",
		Category:    "Technical",
		Description: "User interface and client-side logic.",
		ParentID:    &webDevID,
		Ancestors:   []primitive.ObjectID{webDevID},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), frontend)

	// 3. Grandchild: React
	reactID := primitive.NewObjectID()
	react := models.Skill{
		ID:          reactID,
		Name:        "React",
		Category:    "Technical",
		Description: "A JavaScript library for building user interfaces.",
		ParentID:    &frontendID,
		Ancestors:   []primitive.ObjectID{webDevID, frontendID},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), react)

	// 4. Child Skill: Backend
	backendID := primitive.NewObjectID()
	backend := models.Skill{
		ID:          backendID,
		Name:        "Backend",
		Category:    "Technical",
		Description: "Server-side logic and database management.",
		ParentID:    &webDevID,
		Ancestors:   []primitive.ObjectID{webDevID},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), backend)

	// 5. Grandchild: Go
	goID := primitive.NewObjectID()
	goSkill := models.Skill{
		ID:          goID,
		Name:        "Go",
		Category:    "Technical",
		Description: "A statically typed, compiled programming language designed at Google.",
		ParentID:    &backendID,
		Ancestors:   []primitive.ObjectID{webDevID, backendID},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), goSkill)

	log.Println("Hierarchical skills populated")

	// 6. Progress Items for React
	reactProgress := []models.ProgressItem{
		{
			ID:            primitive.NewObjectID(),
			ParentSkillID: reactID,
			Name:          "JSX and Components",
			Achieved:      true,
			Weightage:     3,
			Comments:      "Mastered basics of JSX and functional components.",
		},
		{
			ID:            primitive.NewObjectID(),
			ParentSkillID: reactID,
			Name:          "Hooks (useState, useEffect)",
			Achieved:      true,
			Weightage:     5,
			Comments:      "Understands core hooks and side effects.",
		},
		{
			ID:            primitive.NewObjectID(),
			ParentSkillID: reactID,
			Name:          "TanStack Query Implementation",
			Achieved:      false,
			Weightage:     8,
			Comments:      "Need to implement state management transition.",
		},
	}

	for _, p := range reactProgress {
		progressCol.InsertOne(context.Background(), p)
	}

	// 7. Progress Items for Go
	goProgress := []models.ProgressItem{
		{
			ID:            primitive.NewObjectID(),
			ParentSkillID: goID,
			Name:          "Syntax and Types",
			Achieved:      true,
			Weightage:     3,
			Comments:      "Comfortable with basic syntax.",
		},
		{
			ID:            primitive.NewObjectID(),
			ParentSkillID: goID,
			Name:          "Goroutines and Channels",
			Achieved:      false,
			Weightage:     13,
			Comments:      "Need to deep dive into concurrency.",
		},
	}

	for _, p := range goProgress {
		progressCol.InsertOne(context.Background(), p)
	}

	log.Println("Progress items populated")
	log.Println("Seeding complete!")
}
