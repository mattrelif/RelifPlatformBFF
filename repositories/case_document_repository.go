package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"relif/platform-bff/models"
)

type CaseDocumentRepository struct {
	collection *mongo.Collection
}

func NewCaseDocumentRepository(db *mongo.Database) *CaseDocumentRepository {
	return &CaseDocumentRepository{
		collection: db.Collection("case_documents"),
	}
}

type CaseDocumentFilters struct {
	CaseID       string   `json:"case_id"`
	DocumentType *string  `json:"document_type,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Search       *string  `json:"search,omitempty"`
	UploadedBy   *string  `json:"uploaded_by,omitempty"`
	Limit        int      `json:"limit"`
	Offset       int      `json:"offset"`
}

func (r *CaseDocumentRepository) Create(ctx context.Context, docModel models.CaseDocument) (string, error) {
	result, err := r.collection.InsertOne(ctx, docModel)
	if err != nil {
		return "", fmt.Errorf("failed to create case document: %w", err)
	}

	switch id := result.InsertedID.(type) {
	case primitive.ObjectID:
		return id.Hex(), nil
	case string:
		return id, nil
	default:
		return "", fmt.Errorf("unexpected insertedID type %T", id)
	}
}

func (r *CaseDocumentRepository) GetByID(ctx context.Context, id string) (*models.CaseDocument, error) {
	// Accept both ObjectID and plain string IDs
	var filter bson.M
	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": objID}
	} else {
		filter = bson.M{"_id": id}
	}

	var docModel models.CaseDocument
	err := r.collection.FindOne(ctx, filter).Decode(&docModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	return &docModel, nil
}

func (r *CaseDocumentRepository) Update(ctx context.Context, id string, docModel models.CaseDocument) error {
	// Accept both ObjectID and plain string IDs
	var filter bson.M
	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": objID}
	} else {
		filter = bson.M{"_id": id}
	}

	// Create update document excluding ID and timestamps
	update := bson.M{
		"$set": bson.M{
			"document_name": docModel.DocumentName,
			"document_type": docModel.DocumentType,
			"description":   docModel.Description,
			"tags":          docModel.Tags,
			"updated_at":    time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("document not found")
	}

	return nil
}

func (r *CaseDocumentRepository) Delete(ctx context.Context, id string) error {
	// Accept both ObjectID and plain string IDs
	var filter bson.M
	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": objID}
	} else {
		filter = bson.M{"_id": id}
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("document not found")
	}

	return nil
}

func (r *CaseDocumentRepository) ListByCaseID(ctx context.Context, filters CaseDocumentFilters) ([]models.CaseDocument, int, error) {
	// Build filter query
	query := bson.M{"case_id": filters.CaseID}

	if filters.DocumentType != nil {
		query["document_type"] = *filters.DocumentType
	}
	if filters.UploadedBy != nil {
		query["uploaded_by_id"] = *filters.UploadedBy
	}

	// Tags filter
	if len(filters.Tags) > 0 {
		query["tags"] = bson.M{"$in": filters.Tags}
	}

	// Search filter
	if filters.Search != nil && *filters.Search != "" {
		searchRegex := primitive.Regex{
			Pattern: *filters.Search,
			Options: "i", // Case insensitive
		}
		query["$or"] = []bson.M{
			{"document_name": searchRegex},
			{"file_name": searchRegex},
			{"description": searchRegex},
			{"tags": searchRegex},
		}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// Build options
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by newest first
	opts.SetSkip(int64(filters.Offset))
	if filters.Limit > 0 {
		opts.SetLimit(int64(filters.Limit))
	}

	// Execute query
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list documents: %w", err)
	}
	defer cursor.Close(ctx)

	var documents []models.CaseDocument
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, 0, fmt.Errorf("failed to decode documents: %w", err)
	}

	return documents, int(total), nil
}
