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

type CaseNoteRepository struct {
	collection *mongo.Collection
}

func NewCaseNoteRepository(db *mongo.Database) *CaseNoteRepository {
	return &CaseNoteRepository{
		collection: db.Collection("case_notes"),
	}
}

type CaseNoteFilters struct {
	CaseID    string   `json:"case_id"`
	NoteType  *string  `json:"note_type,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Important *bool    `json:"is_important,omitempty"`
	Search    *string  `json:"search,omitempty"`
	CreatedBy *string  `json:"created_by,omitempty"`
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
}

func (r *CaseNoteRepository) Create(ctx context.Context, noteModel models.CaseNote) (string, error) {
	result, err := r.collection.InsertOne(ctx, noteModel)
	if err != nil {
		return "", fmt.Errorf("failed to create case note: %w", err)
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

func (r *CaseNoteRepository) GetByID(ctx context.Context, id string) (*models.CaseNote, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid note ID: %w", err)
	}

	var noteModel models.CaseNote
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&noteModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return &noteModel, nil
}

func (r *CaseNoteRepository) Update(ctx context.Context, id string, updates bson.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid note ID: %w", err)
	}

	updates["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

func (r *CaseNoteRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid note ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

func (r *CaseNoteRepository) ListByCaseID(ctx context.Context, filters CaseNoteFilters) ([]models.CaseNote, int, error) {
	// Build filter query
	query := bson.M{"case_id": filters.CaseID}

	if filters.NoteType != nil {
		query["note_type"] = *filters.NoteType
	}
	if filters.Important != nil {
		query["is_important"] = *filters.Important
	}
	if filters.CreatedBy != nil {
		query["created_by_id"] = *filters.CreatedBy
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
			{"title": searchRegex},
			{"content": searchRegex},
			{"tags": searchRegex},
		}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notes: %w", err)
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
		return nil, 0, fmt.Errorf("failed to list notes: %w", err)
	}
	defer cursor.Close(ctx)

	var notes []models.CaseNote
	if err = cursor.All(ctx, &notes); err != nil {
		return nil, 0, fmt.Errorf("failed to decode notes: %w", err)
	}

	return notes, int(total), nil
}
