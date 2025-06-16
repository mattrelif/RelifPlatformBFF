package repositories

import (
	"context"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CaseRepository interface {
	Create(ctx context.Context, case_ entities.Case) (entities.Case, error)
	GetByID(ctx context.Context, id string) (entities.Case, error)
	GetByOrganization(ctx context.Context, organizationID string, filters CaseFilters) ([]entities.Case, int64, error)
	Update(ctx context.Context, id string, case_ entities.Case) (entities.Case, error)
	Delete(ctx context.Context, id string) error
	GetStats(ctx context.Context, organizationID string) (CaseStats, error)
	UpdateNotesCount(ctx context.Context, caseID string, increment int) error
	UpdateDocumentsCount(ctx context.Context, caseID string, increment int) error
}

type CaseFilters struct {
	OrganizationID string     `json:"organization_id"`
	Status         *string    `json:"status,omitempty"`
	Priority       *string    `json:"priority,omitempty"`
	CaseType       *string    `json:"case_type,omitempty"`
	UrgencyLevel   *string    `json:"urgency_level,omitempty"`
	AssignedToID   *string    `json:"assigned_to_id,omitempty"`
	BeneficiaryID  *string    `json:"beneficiary_id,omitempty"`
	FromDate       *time.Time `json:"from_date,omitempty"`
	ToDate         *time.Time `json:"to_date,omitempty"`
	Search         *string    `json:"search,omitempty"`
	Tags           []string   `json:"tags,omitempty"`
	SortBy         string     `json:"sort_by"`
	SortOrder      string     `json:"sort_order"`
	Page           int        `json:"page"`
	Limit          int        `json:"limit"`
	Offset         int        `json:"offset"`
}

type CaseStats struct {
	TotalCases        int `json:"total_cases"`
	OpenCases         int `json:"open_cases"`
	InProgressCases   int `json:"in_progress_cases"`
	OverdueCases      int `json:"overdue_cases"`
	ClosedThisMonth   int `json:"closed_this_month"`
	AvgResolutionDays int `json:"avg_resolution_days"`
}

type caseRepository struct {
	collection *mongo.Collection
}

func NewCaseRepository(database *mongo.Database) CaseRepository {
	return &caseRepository{
		collection: database.Collection("cases"),
	}
}

func (r *caseRepository) Create(ctx context.Context, case_ entities.Case) (entities.Case, error) {
	model := models.NewCase(case_)

	_, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return entities.Case{}, err
	}

	return model.ToEntity(), nil
}

func (r *caseRepository) GetByID(ctx context.Context, id string) (entities.Case, error) {
	var model models.Case

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model)
	if err != nil {
		return entities.Case{}, err
	}

	return model.ToEntity(), nil
}

func (r *caseRepository) GetByOrganization(ctx context.Context, organizationID string, filters CaseFilters) ([]entities.Case, int64, error) {
	// Build filter query
	filter := bson.M{"organization_id": organizationID}

	if filters.Status != nil {
		filter["status"] = *filters.Status
	}
	if filters.Priority != nil {
		filter["priority"] = *filters.Priority
	}
	if filters.CaseType != nil {
		filter["case_type"] = *filters.CaseType
	}
	if filters.AssignedToID != nil {
		filter["assigned_to_id"] = *filters.AssignedToID
	}
	if filters.BeneficiaryID != nil {
		filter["beneficiary_id"] = *filters.BeneficiaryID
	}
	if filters.Search != nil && *filters.Search != "" {
		searchRegex := primitive.Regex{
			Pattern: *filters.Search,
			Options: "i", // Case insensitive
		}
		filter["$or"] = []bson.M{
			{"title": searchRegex},
			{"description": searchRegex},
			{"case_number": searchRegex},
		}
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort
	sortField := "created_at"
	sortOrder := -1
	if filters.SortBy != "" {
		sortField = filters.SortBy
	}
	if filters.SortOrder == "asc" {
		sortOrder = 1
	}

	// Build pagination
	skip := 0
	limit := 20
	if filters.Page > 0 {
		skip = (filters.Page - 1) * filters.Limit
	}
	if filters.Limit > 0 {
		limit = filters.Limit
	}

	opts := options.Find().
		SetSort(bson.D{{sortField, sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var models []models.Case
	if err = cursor.All(ctx, &models); err != nil {
		return nil, 0, err
	}

	cases := make([]entities.Case, len(models))
	for i, model := range models {
		cases[i] = model.ToEntity()
	}

	return cases, total, nil
}

func (r *caseRepository) Update(ctx context.Context, id string, case_ entities.Case) (entities.Case, error) {
	model := models.NewUpdatedCase(case_)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": model}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return entities.Case{}, err
	}

	return r.GetByID(ctx, id)
}

func (r *caseRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *caseRepository) GetStats(ctx context.Context, organizationID string) (CaseStats, error) {
	filter := bson.M{"organization_id": organizationID}

	// Total cases
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return CaseStats{}, err
	}

	// Pending cases (replacing "open")
	pendingFilter := bson.M{"organization_id": organizationID, "status": "PENDING"}
	pendingCount, err := r.collection.CountDocuments(ctx, pendingFilter)
	if err != nil {
		return CaseStats{}, err
	}

	// In progress cases
	progressFilter := bson.M{"organization_id": organizationID, "status": "IN_PROGRESS"}
	progressCount, err := r.collection.CountDocuments(ctx, progressFilter)
	if err != nil {
		return CaseStats{}, err
	}

	// Overdue cases
	now := time.Now()
	overdueFilter := bson.M{
		"organization_id": organizationID,
		"status":          bson.M{"$nin": []string{"CLOSED"}},
		"due_date":        bson.M{"$lt": now},
	}
	overdueCount, err := r.collection.CountDocuments(ctx, overdueFilter)
	if err != nil {
		return CaseStats{}, err
	}

	// Closed this month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	closedFilter := bson.M{
		"organization_id": organizationID,
		"status":          "CLOSED",
		"updated_at":      bson.M{"$gte": startOfMonth},
	}
	closedCount, err := r.collection.CountDocuments(ctx, closedFilter)
	if err != nil {
		return CaseStats{}, err
	}

	// Average resolution time (simplified calculation)
	avgResolutionDays := 0.0

	return CaseStats{
		TotalCases:        int(total),
		OpenCases:         int(pendingCount),
		InProgressCases:   int(progressCount),
		OverdueCases:      int(overdueCount),
		ClosedThisMonth:   int(closedCount),
		AvgResolutionDays: int(avgResolutionDays),
	}, nil
}

func (r *caseRepository) UpdateNotesCount(ctx context.Context, caseID string, increment int) error {
	filter := bson.M{"_id": caseID}
	update := bson.M{
		"$inc": bson.M{"notes_count": increment},
		"$set": bson.M{"last_activity": time.Now(), "updated_at": time.Now()},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *caseRepository) UpdateDocumentsCount(ctx context.Context, caseID string, increment int) error {
	filter := bson.M{"_id": caseID}
	update := bson.M{
		"$inc": bson.M{"documents_count": increment},
		"$set": bson.M{"last_activity": time.Now(), "updated_at": time.Now()},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
