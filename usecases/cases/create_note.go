package cases

import (
	"context"
	"fmt"

	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"relif/platform-bff/repositories"
)

// CreateNoteUseCase defines behaviour for creating a note that belongs to a case.
type CreateNoteUseCase interface {
	Execute(ctx context.Context, user entities.User, caseID string, note entities.CaseNote) (entities.CaseNote, error)
}

type createNoteUseCase struct {
	caseRepo repositories.CaseRepository
	noteRepo *repositories.CaseNoteRepository
	userRepo repositories.Users
}

// NewCreateNoteUseCase wires a new implementation.
func NewCreateNoteUseCase(
	caseRepo repositories.CaseRepository,
	noteRepo *repositories.CaseNoteRepository,
	userRepo repositories.Users,
) CreateNoteUseCase {
	return &createNoteUseCase{
		caseRepo: caseRepo,
		noteRepo: noteRepo,
		userRepo: userRepo,
	}
}

func (uc *createNoteUseCase) Execute(ctx context.Context, user entities.User, caseID string, note entities.CaseNote) (entities.CaseNote, error) {
	// 1. Authorize – ensure user belongs to same organisation as the case.
	caseEntity, err := uc.caseRepo.GetByID(ctx, caseID)
	if err != nil {
		return entities.CaseNote{}, fmt.Errorf("case not found: %w", err)
	}
	if caseEntity.OrganizationID != user.OrganizationID {
		return entities.CaseNote{}, fmt.Errorf("user not authorized for this case")
	}

	// 2. Prepare entity.
	note.CaseID = caseID
	note.CreatedByID = user.ID

	// 3. Persist.
	noteModel := models.NewCaseNoteFromEntity(note)
	noteID, err := uc.noteRepo.Create(ctx, *noteModel)
	if err != nil {
		return entities.CaseNote{}, fmt.Errorf("failed to create note in database: %w", err)
	}

	// 4. Increment notes count (best-effort).
	if err := uc.caseRepo.UpdateNotesCount(ctx, caseID, 1); err != nil {
		fmt.Printf("Warning: failed to update notes count for case %s: %v\n", caseID, err)
	}

	// 5. Retrieve full entity.
	createdNoteModel, err := uc.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return entities.CaseNote{}, fmt.Errorf("failed to retrieve created note: %w", err)
	}
	createdNote := createdNoteModel.ToEntity()
	createdNote.CreatedBy = user

	return createdNote, nil
}

// UpdateNoteUseCase defines behaviour for updating a case note.
type UpdateNoteUseCase interface {
	Execute(ctx context.Context, user entities.User, caseID, noteID string, updates entities.CaseNote) (entities.CaseNote, error)
}

type updateNoteUseCase struct {
	caseRepo repositories.CaseRepository
	noteRepo *repositories.CaseNoteRepository
	userRepo repositories.Users
}

// NewUpdateNoteUseCase builds a new updater implementation.
func NewUpdateNoteUseCase(
	caseRepo repositories.CaseRepository,
	noteRepo *repositories.CaseNoteRepository,
	userRepo repositories.Users,
) UpdateNoteUseCase {
	return &updateNoteUseCase{
		caseRepo: caseRepo,
		noteRepo: noteRepo,
		userRepo: userRepo,
	}
}

func (uc *updateNoteUseCase) Execute(ctx context.Context, user entities.User, caseID, noteID string, updates entities.CaseNote) (entities.CaseNote, error) {
	// 1. Authorize.
	caseEntity, err := uc.caseRepo.GetByID(ctx, caseID)
	if err != nil {
		return entities.CaseNote{}, fmt.Errorf("case not found: %w", err)
	}
	if caseEntity.OrganizationID != user.OrganizationID {
		return entities.CaseNote{}, fmt.Errorf("user not authorized for this case")
	}

	// 2. Build update map – only include provided fields (non-zero).
	set := make(map[string]interface{})
	if updates.Title != "" {
		set["title"] = updates.Title
	}
	if updates.Content != "" {
		set["content"] = updates.Content
	}
	if len(updates.Tags) > 0 {
		set["tags"] = updates.Tags
	}
	if updates.NoteType != "" {
		set["note_type"] = updates.NoteType
	}
	// Bool has default false; we distinguish by using a pointer in request => here expecting caller to express update via value. We'll include IsImportant as-is.
	set["is_important"] = updates.IsImportant

	if len(set) == 0 {
		// Nothing to update, just return current note.
		current, err := uc.noteRepo.GetByID(ctx, noteID)
		if err != nil {
			return entities.CaseNote{}, fmt.Errorf("note not found: %w", err)
		}
		noteEntity := current.ToEntity()
		if noteEntity.CreatedByID != "" {
			if creator, err := uc.userRepo.FindOneByID(noteEntity.CreatedByID); err == nil {
				noteEntity.CreatedBy = creator
			}
		}
		return noteEntity, nil
	}

	// 3. Persist.
	if err := uc.noteRepo.Update(ctx, noteID, set); err != nil {
		return entities.CaseNote{}, err
	}

	// 4. Retrieve full entity.
	updatedModel, err := uc.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return entities.CaseNote{}, fmt.Errorf("failed to retrieve updated note: %w", err)
	}

	updatedEntity := updatedModel.ToEntity()
	if updatedEntity.CreatedByID != "" {
		if creator, err := uc.userRepo.FindOneByID(updatedEntity.CreatedByID); err == nil {
			updatedEntity.CreatedBy = creator
		}
	}

	return updatedEntity, nil
}
