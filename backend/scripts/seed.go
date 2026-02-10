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

	// Helper function to create a skill
	createSkill := func(name, category, description string, parentID *primitive.ObjectID) models.Skill {
		id := primitive.NewObjectID()
		ancestors := []primitive.ObjectID{}
		if parentID != nil {
			ancestors = append(ancestors, *parentID)
		}

		return models.Skill{
			ID:          id,
			Name:        name,
			Category:    category,
			Description: description,
			ParentID:    parentID,
			Ancestors:   ancestors,
			CreatedAt:   time.Now(),
		}
	}

	// Helper to create progress items
	createProgress := func(parentID primitive.ObjectID, name string, achieved bool, weightage int, comments string) models.ProgressItem {
		return models.ProgressItem{
			ID:            primitive.NewObjectID(),
			ParentSkillID: parentID,
			Name:          name,
			Achieved:      achieved,
			Weightage:     weightage,
			Comments:      comments,
		}
	}

	// --- 1. Frontend Mastery ---
	frontendRoot := createSkill("Frontend Mastery", "Career Path", "Mastering modern frontend development.", nil)
	skillCol.InsertOne(context.Background(), frontendRoot)

	// React
	react := createSkill("React", "Library", "Building user interfaces with React.", &frontendRoot.ID)
	skillCol.InsertOne(context.Background(), react)

	pReact := []models.ProgressItem{
		createProgress(react.ID, "Hooks (useState, useEffect)", true, 10, "Comfortable with basics"),
		createProgress(react.ID, "Context API", true, 8, "Used in last project"),
		createProgress(react.ID, "Performance Optimization", false, 12, "Need to learn memo and useCallback"),
		createProgress(react.ID, "Custom Hooks", true, 8, ""),
	}
	for _, p := range pReact {
		progressCol.InsertOne(context.Background(), p)
	}

	// CSS
	css := createSkill("CSS & Styling", "Language", "Styling web applications.", &frontendRoot.ID)
	skillCol.InsertOne(context.Background(), css)

	pCSS := []models.ProgressItem{
		createProgress(css.ID, "Flexbox", true, 8, ""),
		createProgress(css.ID, "Grid", true, 8, ""),
		createProgress(css.ID, "Animations & Keyframes", false, 5, "Tricky"),
		createProgress(css.ID, "Responsive Design", true, 10, ""),
	}
	for _, p := range pCSS {
		progressCol.InsertOne(context.Background(), p)
	}

	// TypeScript
	ts := createSkill("TypeScript", "Language", "Typed superset of JavaScript.", &frontendRoot.ID)
	skillCol.InsertOne(context.Background(), ts)

	pTS := []models.ProgressItem{
		createProgress(ts.ID, "Basic Types & Interfaces", true, 5, ""),
		createProgress(ts.ID, "Generics", false, 10, "Concept still fuzzy"),
		createProgress(ts.ID, "Utility Types", false, 5, ""),
	}
	for _, p := range pTS {
		progressCol.InsertOne(context.Background(), p)
	}

	// --- 2. Backend Mastery ---
	backendRoot := createSkill("Backend Mastery", "Career Path", "Server-side logic and databases.", nil)
	skillCol.InsertOne(context.Background(), backendRoot)

	// Node.js
	node := createSkill("Node.js", "Runtime", "JavaScript on the server.", &backendRoot.ID)
	skillCol.InsertOne(context.Background(), node)

	pNode := []models.ProgressItem{
		createProgress(node.ID, "Event Loop", true, 10, "Understood"),
		createProgress(node.ID, "Streams & Buffers", false, 8, "Need practice"),
		createProgress(node.ID, "File System API", true, 5, ""),
	}
	for _, p := range pNode {
		progressCol.InsertOne(context.Background(), p)
	}

	// Database
	dbSkill := createSkill("Databases", "Infrastructure", "Storing and retrieving data.", &backendRoot.ID)
	skillCol.InsertOne(context.Background(), dbSkill)

	pDB := []models.ProgressItem{
		createProgress(dbSkill.ID, "SQL vs NoSQL", true, 5, ""),
		createProgress(dbSkill.ID, "Indexing & Performance", false, 10, "Critical for interview"),
		createProgress(dbSkill.ID, "ACID Properties", true, 5, ""),
		createProgress(dbSkill.ID, "Normalization", true, 5, ""),
	}
	for _, p := range pDB {
		progressCol.InsertOne(context.Background(), p)
	}

	// Go
	goLang := createSkill("Go", "Language", "Efficient, compiled language.", &backendRoot.ID)
	skillCol.InsertOne(context.Background(), goLang)

	pGo := []models.ProgressItem{
		createProgress(goLang.ID, "Goroutines & Concurrency", true, 12, "Love this feature"),
		createProgress(goLang.ID, "Channels", true, 8, ""),
		createProgress(goLang.ID, "Interfaces", false, 8, "Different from Java/TS"),
		createProgress(goLang.ID, "Error Handling", true, 5, ""),
	}
	for _, p := range pGo {
		progressCol.InsertOne(context.Background(), p)
	}

	// --- 3. DevOps ---
	devOpsRoot := createSkill("DevOps", "Career Path", "Deployment and operations.", nil)
	skillCol.InsertOne(context.Background(), devOpsRoot)

	// Docker
	docker := createSkill("Docker", "Containerization", "Containerizing applications.", &devOpsRoot.ID)
	skillCol.InsertOne(context.Background(), docker)

	pDocker := []models.ProgressItem{
		createProgress(docker.ID, "Images & Containers", true, 5, ""),
		createProgress(docker.ID, "Docker Compose", true, 8, "Used frequently"),
		createProgress(docker.ID, "Networking", false, 5, ""),
	}
	for _, p := range pDocker {
		progressCol.InsertOne(context.Background(), p)
	}

	// Kubernetes
	k8s := createSkill("Kubernetes", "Orchestration", "Managing containerized applications.", &devOpsRoot.ID)
	skillCol.InsertOne(context.Background(), k8s)

	pK8s := []models.ProgressItem{
		createProgress(k8s.ID, "Pods & Services", false, 8, "Just started"),
		createProgress(k8s.ID, "Deployments", false, 8, ""),
		createProgress(k8s.ID, "Ingress", false, 8, ""),
	}
	for _, p := range pK8s {
		progressCol.InsertOne(context.Background(), p)
	}

	log.Println("Seeding complete with new comprehensive mock data!")
}
