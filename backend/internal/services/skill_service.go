package services

import (
	"context"
	"math"
	"time"

	"masterypath/internal/apperror"
	"masterypath/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SkillService handles business logic for Skill operations.
type SkillService struct {
	skillCollection    *mongo.Collection
	metadataCollection *mongo.Collection
	progressCollection *mongo.Collection
}

// NewSkillService creates a new SkillService with the required collections.
func NewSkillService(skillCol, metadataCol, progressCol *mongo.Collection) *SkillService {
	return &SkillService{
		skillCollection:    skillCol,
		metadataCollection: metadataCol,
		progressCollection: progressCol,
	}
}

// isValidCategory checks if a category is allowed by querying the metadata collection.
func (s *SkillService) isValidCategory(ctx context.Context, category string) bool {
	filter := bson.M{
		"type":      "SKILL_CATEGORY",
		"value":     category,
		"is_active": true,
	}
	count, err := s.metadataCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}
	return count > 0
}

// GetAll retrieves all skills from the database.
func (s *SkillService) GetAll(ctx context.Context) ([]models.Skill, error) {
	var skills []models.Skill
	cursor, err := s.skillCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to fetch skills")
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &skills); err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to decode skills")
	}

	if skills == nil {
		skills = []models.Skill{}
	}
	return skills, nil
}

// GetAllPopulated retrieves all skills with parent data populated.
func (s *SkillService) GetAllPopulated(ctx context.Context) ([]models.SkillPopulated, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "skills"},
			{Key: "localField", Value: "parent_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "parent_data"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$parent_data"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "parent_name", Value: bson.D{
				{Key: "$ifNull", Value: bson.A{"$parent_data.name", nil}},
			}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "parent_data", Value: 0},
		}}},
	}

	cursor, err := s.skillCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to fetch skills")
	}
	defer cursor.Close(ctx)

	var skills []models.SkillPopulated
	if err := cursor.All(ctx, &skills); err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to decode populated skills")
	}

	if skills == nil {
		skills = []models.SkillPopulated{}
	}
	return skills, nil
}

// GetTree retrieves a skill and all its descendants by root ID.
func (s *SkillService) GetTree(ctx context.Context, rootID primitive.ObjectID) ([]models.Skill, error) {
	filter := bson.M{
		"$or": bson.A{
			bson.M{"_id": rootID},
			bson.M{"ancestors": rootID},
		},
	}

	cursor, err := s.skillCollection.Find(ctx, filter)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to fetch skill tree")
	}
	defer cursor.Close(ctx)

	var skills []models.Skill
	if err := cursor.All(ctx, &skills); err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to decode skill tree")
	}

	if skills == nil {
		skills = []models.Skill{}
	}
	return skills, nil
}

// GetByID retrieves a single skill by its ObjectID with parent name populated.
func (s *SkillService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.SkillPopulated, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "skills"},
			{Key: "localField", Value: "parent_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "parent_data"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$parent_data"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "parent_name", Value: bson.D{
				{Key: "$ifNull", Value: bson.A{"$parent_data.name", nil}},
			}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "parent_data", Value: 0},
		}}},
	}

	cursor, err := s.skillCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to fetch skill")
	}
	defer cursor.Close(ctx)

	var results []models.SkillPopulated
	if err := cursor.All(ctx, &results); err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to decode skill")
	}

	if len(results) == 0 {
		return nil, apperror.ErrNotFound
	}
	return &results[0], nil
}

// Create inserts a new skill, computing the materialized path ancestors.
func (s *SkillService) Create(ctx context.Context, skill *models.Skill) error {
	// Validate category
	if skill.Category != "" && !s.isValidCategory(ctx, skill.Category) {
		return apperror.New(400, "Invalid category: "+skill.Category)
	}

	skill.ID = primitive.NewObjectID()
	skill.CreatedAt = time.Now()
	skill.UpdatedAt = time.Now()

	// Build Ancestors array (Materialized Path)
	if skill.ParentID != nil {
		var parent models.Skill
		err := s.skillCollection.FindOne(ctx, bson.M{"_id": *skill.ParentID}).Decode(&parent)
		if err != nil {
			return apperror.New(400, "Parent skill not found")
		}
		skill.Ancestors = append(parent.Ancestors, parent.ID)
	} else {
		skill.Ancestors = []primitive.ObjectID{} // Root skill
	}

	_, err := s.skillCollection.InsertOne(ctx, skill)
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to insert skill")
	}
	return nil
}

// Update modifies an existing skill and cascades ancestor updates to descendants.
func (s *SkillService) Update(ctx context.Context, id primitive.ObjectID, updateData *models.Skill) error {
	// Validate category
	if updateData.Category != "" && !s.isValidCategory(ctx, updateData.Category) {
		return apperror.New(400, "Invalid category: "+updateData.Category)
	}

	// Recompute ancestors
	var newAncestors []primitive.ObjectID
	if updateData.ParentID != nil {
		var parent models.Skill
		err := s.skillCollection.FindOne(ctx, bson.M{"_id": *updateData.ParentID}).Decode(&parent)
		if err != nil {
			return apperror.New(400, "Parent skill not found")
		}
		newAncestors = append(parent.Ancestors, parent.ID)
	} else {
		newAncestors = []primitive.ObjectID{}
	}

	update := bson.M{
		"$set": bson.M{
			"name":        updateData.Name,
			"category":    updateData.Category,
			"description": updateData.Description,
			"parent_id":   updateData.ParentID,
			"ancestors":   newAncestors,
			"updated_at":  time.Now(),
		},
	}

	_, err := s.skillCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to update skill")
	}

	// Cascade: Update ancestors for all descendants
	s.cascadeAncestorUpdate(ctx, id, newAncestors)

	return nil
}

// cascadeAncestorUpdate updates ancestors for all descendants of a given skill.
func (s *SkillService) cascadeAncestorUpdate(ctx context.Context, skillID primitive.ObjectID, newAncestors []primitive.ObjectID) {
	descendantFilter := bson.M{"ancestors": skillID}
	cursor, err := s.skillCollection.Find(ctx, descendantFilter)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	var descendants []models.Skill
	if err := cursor.All(ctx, &descendants); err != nil {
		return
	}

	for _, desc := range descendants {
		idx := -1
		for i, a := range desc.Ancestors {
			if a == skillID {
				idx = i
				break
			}
		}
		if idx != -1 {
			updatedAncestors := append(newAncestors, skillID)
			updatedAncestors = append(updatedAncestors, desc.Ancestors[idx+1:]...)
			s.skillCollection.UpdateOne(ctx, bson.M{"_id": desc.ID}, bson.M{"$set": bson.M{"ancestors": updatedAncestors}})
		}
	}
}

// Delete removes a skill from the database by its ObjectID.
func (s *SkillService) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.skillCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return apperror.Wrap(apperror.ErrInternal, "failed to delete skill")
	}
	return nil
}

// GetDashboard returns all root-level (master) skills with their direct sub-skills
// and per-sub-skill mastery percentages computed from progress items.
func (s *SkillService) GetDashboard(ctx context.Context) ([]models.DashboardSkill, error) {
	// 1. Find all root skills (no parent)
	rootFilter := bson.M{
		"$or": bson.A{
			bson.M{"parent_id": nil},
			bson.M{"parent_id": bson.M{"$exists": false}},
		},
	}

	cursor, err := s.skillCollection.Find(ctx, rootFilter)
	if err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to fetch root skills")
	}
	defer cursor.Close(ctx)

	var rootSkills []models.Skill
	if err := cursor.All(ctx, &rootSkills); err != nil {
		return nil, apperror.Wrap(apperror.ErrInternal, "failed to decode root skills")
	}

	var dashboard []models.DashboardSkill

	for _, root := range rootSkills {
		// 2. Find direct children of this root
		childCursor, err := s.skillCollection.Find(ctx, bson.M{"parent_id": root.ID})
		if err != nil {
			continue
		}

		var children []models.Skill
		if err := childCursor.All(ctx, &children); err != nil {
			childCursor.Close(ctx)
			continue
		}
		childCursor.Close(ctx)

		var subSkills []models.DashboardSubSkill
		var totalMastery float64

		for _, child := range children {
			// 3. Compute mastery for each sub-skill from progress items
			mastery := computeMastery(ctx, s.progressCollection, child.ID)

			subSkills = append(subSkills, models.DashboardSubSkill{
				ID:      child.ID,
				Name:    child.Name,
				Mastery: mastery,
			})
			totalMastery += mastery
		}

		if subSkills == nil {
			subSkills = []models.DashboardSubSkill{}
		}

		var overallMastery float64
		if len(subSkills) > 0 {
			overallMastery = totalMastery / float64(len(subSkills))
		}

		dashboard = append(dashboard, models.DashboardSkill{
			ID:             root.ID,
			Name:           root.Name,
			Category:       root.Category,
			Description:    root.Description,
			SubSkills:      subSkills,
			OverallMastery: math.Round(overallMastery*10) / 10,
		})
	}

	if dashboard == nil {
		dashboard = []models.DashboardSkill{}
	}
	return dashboard, nil
}

// computeMastery calculates the mastery % for a skill based on its progress items.
// Mastery = (sum of weightage for achieved items) / (sum of all weightage) * 100
func computeMastery(ctx context.Context, progressCol *mongo.Collection, skillID primitive.ObjectID) float64 {
	// Get all progress items for this skill
	cursor, err := progressCol.Find(ctx, bson.M{"parent_skill_id": skillID})
	if err != nil {
		return 0
	}
	defer cursor.Close(ctx)

	var items []models.ProgressItem
	if err := cursor.All(ctx, &items); err != nil {
		return 0
	}

	if len(items) == 0 {
		return 0
	}

	var totalWeight, achievedWeight int
	for _, item := range items {
		totalWeight += item.Weightage
		if item.Achieved {
			achievedWeight += item.Weightage
		}
	}

	if totalWeight == 0 {
		return 0
	}
	return math.Round((float64(achievedWeight)/float64(totalWeight)*100)*10) / 10
}
