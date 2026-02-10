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

	// 1. Root Skill: Salesforce Development
	sfDevID := primitive.NewObjectID()
	sfDev := models.Skill{
		ID:          sfDevID,
		Name:        "Salesforce Development",
		Category:    "Technology",
		Description: "Mastering the Salesforce platform development ecosystem.",
		Ancestors:   []primitive.ObjectID{},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), sfDev)

	// 2. Child Skill: Apex
	apexID := primitive.NewObjectID()
	apex := models.Skill{
		ID:          apexID,
		Name:        "Apex",
		Category:    "Programming Language",
		Description: "Proprietary language for Salesforce backend logic.",
		ParentID:    &sfDevID,
		Ancestors:   []primitive.ObjectID{sfDevID},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), apex)

	// 3. Child Skill: Triggers
	triggersID := primitive.NewObjectID()
	triggers := models.Skill{
		ID:          triggersID,
		Name:        "Triggers",
		Category:    "Backend Logic",
		Description: "Automating logic on database events.",
		ParentID:    &sfDevID,
		Ancestors:   []primitive.ObjectID{sfDevID},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), triggers)

	// 4. Child Skill: LWC
	lwcID := primitive.NewObjectID()
	lwc := models.Skill{
		ID:          lwcID,
		Name:        "LWC",
		Category:    "Frontend Framework",
		Description: "Modern web components for Salesforce UI.",
		ParentID:    &sfDevID,
		Ancestors:   []primitive.ObjectID{sfDevID},
		CreatedAt:   time.Now(),
	}
	skillCol.InsertOne(context.Background(), lwc)

	log.Println("Hierarchical skills populated")

	// 5. Progress Items for Apex (Aiming for ~65%)
	// Total Weight: 13 (Syntax) + 8 (Collections) + 5 (SOQL) = 26
	// Achieved: 13 + 3 (partial credit? no, boolean) -> 13 + 5 = 18 / 26 = 69%
	apexProgress := []models.ProgressItem{
		{ID: primitive.NewObjectID(), ParentSkillID: apexID, Name: "Syntax & Basics", Achieved: true, Weightage: 13, Comments: "Solid understanding"},
		{ID: primitive.NewObjectID(), ParentSkillID: apexID, Name: "Collections (List, Set, Map)", Achieved: false, Weightage: 8, Comments: "Need more practice with Maps"},
		{ID: primitive.NewObjectID(), ParentSkillID: apexID, Name: "SOQL Queries", Achieved: true, Weightage: 5, Comments: "Good with basic queries"},
	}
	for _, p := range apexProgress {
		progressCol.InsertOne(context.Background(), p)
	}

	// 6. Progress Items for Triggers (Aiming for ~40%)
	// Total: 8 + 5 + 3 = 16
	// Achieved: 5 = 5/16 = 31%
	triggersProgress := []models.ProgressItem{
		{ID: primitive.NewObjectID(), ParentSkillID: triggersID, Name: "Context Variables", Achieved: true, Weightage: 5},
		{ID: primitive.NewObjectID(), ParentSkillID: triggersID, Name: "Bulkification Patterns", Achieved: false, Weightage: 8},
		{ID: primitive.NewObjectID(), ParentSkillID: triggersID, Name: "Order of Execution", Achieved: false, Weightage: 3},
	}
	for _, p := range triggersProgress {
		progressCol.InsertOne(context.Background(), p)
	}

	// 7. Progress Items for LWC (Aiming for ~80%)
	// Total: 5 + 5 + 5 + 5 = 20
	// Achieved: 5+5+5 = 15/20 = 75%
	lwcProgress := []models.ProgressItem{
		{ID: primitive.NewObjectID(), ParentSkillID: lwcID, Name: "Component Structure", Achieved: true, Weightage: 5},
		{ID: primitive.NewObjectID(), ParentSkillID: lwcID, Name: "Wire Service", Achieved: true, Weightage: 5},
		{ID: primitive.NewObjectID(), ParentSkillID: lwcID, Name: "Events & Communication", Achieved: false, Weightage: 5},
		{ID: primitive.NewObjectID(), ParentSkillID: lwcID, Name: "Lightning Data Service", Achieved: true, Weightage: 5},
	}
	for _, p := range lwcProgress {
		progressCol.InsertOne(context.Background(), p)
	}

	log.Println("Progress items populated")
	log.Println("Seeding complete!")
}
