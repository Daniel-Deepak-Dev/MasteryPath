package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// DashboardSubSkill represents a sub-skill with its mastery percentage for the dashboard view.
type DashboardSubSkill struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Name    string             `json:"name" bson:"name"`
	Mastery float64            `json:"mastery" bson:"mastery"`
}

// DashboardSkill represents a root-level (master) skill with its sub-skills for the dashboard.
type DashboardSkill struct {
	ID             primitive.ObjectID  `json:"id" bson:"_id"`
	Name           string              `json:"name" bson:"name"`
	Category       string              `json:"category" bson:"category"`
	Description    string              `json:"description" bson:"description"`
	SubSkills      []DashboardSubSkill `json:"sub_skills"`
	OverallMastery float64             `json:"overall_mastery"`
}
