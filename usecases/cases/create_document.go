package cases

import (
	"context"
	"fmt"

	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"relif/platform-bff/repositories"
)

type CreateDocumentUseCase interface {
	Execute(ctx context.Context, user entities.User, caseID string, document entities.CaseDocument) (entities.CaseDocument, error)
}

type createDocumentUseCase struct {
	caseRepo repositories.CaseRepository
	docRepo  *repositories.CaseDocumentRepository
	userRepo repositories.Users
}

func NewCreateDocumentUseCase(
	caseRepo repositories.CaseRepository,
	docRepo *repositories.CaseDocumentRepository,
	userRepo repositories.Users,
) CreateDocumentUseCase {
	return &createDocumentUseCase{
		caseRepo: caseRepo,
		docRepo:  docRepo,
		userRepo: userRepo,
	}
}

func (uc *createDocumentUseCase) Execute(ctx context.Context, user entities.User, caseID string, document entities.CaseDocument) (entities.CaseDocument, error) {
	// 1. Authorize
	caseEntity, err := uc.caseRepo.GetByID(ctx, caseID)
	if err != nil {
		return entities.CaseDocument{}, fmt.Errorf("case not found: %w", err)
	}

	// Example guard, you can make this more specific
	if caseEntity.OrganizationID != user.OrganizationID {
		return entities.CaseDocument{}, fmt.Errorf("user not authorized for this case")
	}

	// 2. Prepare entity
	document.CaseID = caseID
	document.UploadedByID = user.ID

	// 3. Create in database
	docModel := models.NewCaseDocumentFromEntity(document)
	docID, err := uc.docRepo.Create(ctx, *docModel)
	if err != nil {
		return entities.CaseDocument{}, fmt.Errorf("failed to create document in database: %w", err)
	}

	// 4. Increment documents count
	if err := uc.caseRepo.UpdateDocumentsCount(ctx, caseID, 1); err != nil {
		// Log this error, but don't fail the entire operation since the doc was created
		fmt.Printf("Warning: failed to update documents count for case %s: %v\n", caseID, err)
	}

	// 5. Return the full entity
	createdDocModel, err := uc.docRepo.GetByID(ctx, docID)
	if err != nil {
		return entities.CaseDocument{}, fmt.Errorf("failed to retrieve created document: %w", err)
	}

	createdDocEntity := createdDocModel.ToEntity()
	createdDocEntity.UploadedBy = user // Ensure the full user object is attached

	return createdDocEntity, nil
}

type UpdateDocumentUseCase interface {
	Execute(ctx context.Context, user entities.User, caseID, documentID string, document entities.CaseDocument) (entities.CaseDocument, error)
}

type updateDocumentUseCase struct {
	caseRepo repositories.CaseRepository
	docRepo  *repositories.CaseDocumentRepository
	userRepo repositories.Users
}

func NewUpdateDocumentUseCase(
	caseRepo repositories.CaseRepository,
	docRepo *repositories.CaseDocumentRepository,
	userRepo repositories.Users,
) UpdateDocumentUseCase {
	return &updateDocumentUseCase{
		caseRepo: caseRepo,
		docRepo:  docRepo,
		userRepo: userRepo,
	}
}

func (uc *updateDocumentUseCase) Execute(ctx context.Context, user entities.User, caseID, documentID string, document entities.CaseDocument) (entities.CaseDocument, error) {
	// 1. Authorize
	caseEntity, err := uc.caseRepo.GetByID(ctx, caseID)
	if err != nil {
		return entities.CaseDocument{}, fmt.Errorf("case not found: %w", err)
	}
	if caseEntity.OrganizationID != user.OrganizationID {
		return entities.CaseDocument{}, fmt.Errorf("user not authorized for this case")
	}

	// 2. Update in database
	updateModel := models.NewCaseDocumentFromEntity(document)
	if err := uc.docRepo.Update(ctx, documentID, *updateModel); err != nil {
		return entities.CaseDocument{}, fmt.Errorf("failed to update document: %w", err)
	}

	// 3. Return the full entity
	updatedDocModel, err := uc.docRepo.GetByID(ctx, documentID)
	if err != nil {
		return entities.CaseDocument{}, fmt.Errorf("failed to retrieve updated document: %w", err)
	}

	updatedDocEntity := updatedDocModel.ToEntity()
	// Populate uploaded by user since it's not in the model
	if updatedDocEntity.UploadedByID != "" {
		if uploader, err := uc.userRepo.FindOneByID(updatedDocEntity.UploadedByID); err == nil {
			updatedDocEntity.UploadedBy = uploader
		}
	}

	return updatedDocEntity, nil
}
