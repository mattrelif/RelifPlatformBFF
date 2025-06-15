package cases

import (
	"context"
	"fmt"
	"time"

	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/models"
	"relif/platform-bff/repositories"
)

// Simple interfaces to avoid circular dependencies
type BeneficiaryService interface {
	GetByID(ctx context.Context, id string) (entities.Beneficiary, error)
}

type UserService interface {
	GetByID(ctx context.Context, id string) (entities.User, error)
}

type CaseUseCase struct {
	caseRepo      repositories.CaseRepository
	noteRepo      *repositories.CaseNoteRepository
	beneficiaryUC BeneficiaryService
	userUC        UserService
}

func NewCaseUseCase(
	caseRepo repositories.CaseRepository,
	noteRepo *repositories.CaseNoteRepository,
	beneficiaryUC BeneficiaryService,
	userUC UserService,
) *CaseUseCase {
	return &CaseUseCase{
		caseRepo:      caseRepo,
		noteRepo:      noteRepo,
		beneficiaryUC: beneficiaryUC,
		userUC:        userUC,
	}
}

func (uc *CaseUseCase) CreateCase(ctx context.Context, req requests.CreateCase, userID, organizationID string) (*responses.CaseResponse, error) {
	// Validate beneficiary exists
	beneficiary, err := uc.beneficiaryUC.GetByID(ctx, req.BeneficiaryID)
	if err != nil {
		return nil, fmt.Errorf("beneficiary not found: %w", err)
	}

	// Validate assigned user exists
	assignedUser, err := uc.userUC.GetByID(ctx, req.AssignedToID)
	if err != nil {
		return nil, fmt.Errorf("assigned user not found: %w", err)
	}

	// Validate organization boundaries
	if err := req.ValidateOrganizationBoundaries(organizationID, beneficiary, assignedUser); err != nil {
		return nil, err
	}

	// Create case entity using the request's ToEntity method
	caseEntity := req.ToEntity(organizationID)
	caseEntity.Status = "OPEN" // Always start as OPEN
	caseEntity.CreatedAt = time.Now()
	caseEntity.UpdatedAt = time.Now()
	caseEntity.LastActivity = time.Now()

	// Create the case
	createdCase, err := uc.caseRepo.Create(ctx, caseEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to create case: %w", err)
	}

	// If initial note provided, create it
	if req.InitialNote != nil {
		noteEntity := req.ToInitialNoteEntity(createdCase.ID, userID)
		if noteEntity != nil {
			noteEntity.CreatedAt = time.Now()
			noteEntity.UpdatedAt = time.Now()

			// Convert entity to model for repository
			noteModel := models.NewCaseNoteFromEntity(*noteEntity)
			_, err = uc.noteRepo.Create(ctx, *noteModel)
			if err != nil {
				// Log error but don't fail case creation
				fmt.Printf("Warning: failed to create initial note: %v\n", err)
			} else {
				// Increment notes count
				uc.caseRepo.UpdateNotesCount(ctx, createdCase.ID, 1)
			}
		}
	}

	// Get the created case with populated relations
	return uc.GetByID(ctx, createdCase.ID)
}

func (uc *CaseUseCase) GetByID(ctx context.Context, id string) (*responses.CaseResponse, error) {
	// Get case from repository
	caseEntity, err := uc.caseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get beneficiary details
	beneficiary, err := uc.beneficiaryUC.GetByID(ctx, caseEntity.BeneficiaryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get beneficiary: %w", err)
	}

	// Get assigned user details
	assignedUser, err := uc.userUC.GetByID(ctx, caseEntity.AssignedToID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned user: %w", err)
	}

	// Populate relations in the entity
	caseEntity.Beneficiary = beneficiary
	caseEntity.AssignedTo = assignedUser

	// Convert to response
	response := responses.NewCaseResponse(caseEntity)
	return &response, nil
}

func (uc *CaseUseCase) UpdateCase(ctx context.Context, id string, req requests.UpdateCase, organizationID string) (*responses.CaseResponse, error) {
	// Get existing case
	existingCase, err := uc.caseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate organization boundaries for reassignment
	if req.AssignedToID != nil {
		assignedUser, err := uc.userUC.GetByID(ctx, *req.AssignedToID)
		if err != nil {
			return nil, fmt.Errorf("assigned user not found: %w", err)
		}
		if err := req.ValidateOrganizationBoundaries(organizationID, &assignedUser); err != nil {
			return nil, err
		}
	}

	// Convert request to entity with updates
	updateEntity := req.ToEntity()
	updateEntity.ID = existingCase.ID
	updateEntity.UpdatedAt = time.Now()
	updateEntity.LastActivity = time.Now()

	// Update the case
	updatedCase, err := uc.caseRepo.Update(ctx, id, updateEntity)
	if err != nil {
		return nil, err
	}

	// Return updated case with populated relations
	return uc.GetByID(ctx, updatedCase.ID)
}

func (uc *CaseUseCase) DeleteCase(ctx context.Context, id, organizationID string) error {
	// Get case to verify organization
	caseEntity, err := uc.caseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if caseEntity.OrganizationID != organizationID {
		return fmt.Errorf("case not found in organization")
	}

	return uc.caseRepo.Delete(ctx, id)
}

func (uc *CaseUseCase) ListByOrganization(ctx context.Context, organizationID string, filters repositories.CaseFilters) (*responses.CaseListResponse, error) {
	// Ensure organization ID is set in filters
	filters.OrganizationID = organizationID

	// Get cases from repository
	cases, total, err := uc.caseRepo.GetByOrganization(ctx, organizationID, filters)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	caseResponses := make([]responses.CaseResponse, len(cases))
	for i, caseEntity := range cases {
		// Get beneficiary and assigned user details for each case
		if caseEntity.BeneficiaryID != "" {
			beneficiary, err := uc.beneficiaryUC.GetByID(ctx, caseEntity.BeneficiaryID)
			if err == nil {
				caseEntity.Beneficiary = beneficiary
			}
		}
		if caseEntity.AssignedToID != "" {
			assignedUser, err := uc.userUC.GetByID(ctx, caseEntity.AssignedToID)
			if err == nil {
				caseEntity.AssignedTo = assignedUser
			}
		}

		caseResponses[i] = responses.NewCaseResponse(caseEntity)
	}

	return &responses.CaseListResponse{
		Count: int(total),
		Data:  caseResponses,
	}, nil
}

func (uc *CaseUseCase) GetCaseStats(ctx context.Context, organizationID string) (*responses.CaseStatsResponse, error) {
	stats, err := uc.caseRepo.GetStats(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	return &responses.CaseStatsResponse{
		TotalCases:        stats.TotalCases,
		OpenCases:         stats.OpenCases,
		InProgressCases:   stats.InProgressCases,
		OverdueCases:      stats.OverdueCases,
		ClosedThisMonth:   stats.ClosedThisMonth,
		AvgResolutionDays: stats.AvgResolutionDays,
	}, nil
}
