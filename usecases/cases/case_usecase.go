package cases

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/models"
	"relif/platform-bff/repositories"
	"relif/platform-bff/usecases/beneficiaries"
	"relif/platform-bff/usecases/users"
)

// Simple interfaces to avoid circular dependencies
type BeneficiaryService interface {
	GetByID(ctx context.Context, id string) (entities.Beneficiary, error)
}

type UserService interface {
	GetByID(ctx context.Context, id string) (entities.User, error)
}

type CaseUseCase struct {
	caseRepo      *repositories.CaseRepository
	noteRepo      *repositories.CaseNoteRepository
	beneficiaryUC *beneficiaries.BeneficiaryUseCase
	userUC        *users.UserUseCase
}

func NewCaseUseCase(
	caseRepo *repositories.CaseRepository,
	noteRepo *repositories.CaseNoteRepository,
	beneficiaryUC *beneficiaries.BeneficiaryUseCase,
	userUC *users.UserUseCase,
) *CaseUseCase {
	return &CaseUseCase{
		caseRepo:      caseRepo,
		noteRepo:      noteRepo,
		beneficiaryUC: beneficiaryUC,
		userUC:        userUC,
	}
}

func (uc *CaseUseCase) CreateCase(ctx context.Context, req requests.CreateCaseRequest) (*responses.CaseResponse, error) {
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

	// Create case entity
	caseEntity := entities.Case{
		Title:          req.Title,
		Description:    req.Description,
		Status:         "OPEN", // Always start as OPEN
		Priority:       req.Priority,
		UrgencyLevel:   req.UrgencyLevel,
		CaseType:       req.CaseType,
		BeneficiaryID:  req.BeneficiaryID,
		Beneficiary:    *beneficiary,
		AssignedToID:   req.AssignedToID,
		AssignedTo:     *assignedUser,
		Tags:           req.Tags,
		OrganizationID: beneficiary.CurrentOrganizationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastActivity:   time.Now(),
	}

	// Handle optional fields
	if req.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due date format: %w", err)
		}
		caseEntity.DueDate = &dueDate
	}
	if req.EstimatedDuration != nil {
		caseEntity.EstimatedDuration = *req.EstimatedDuration
	}
	if req.BudgetAllocated != nil {
		caseEntity.BudgetAllocated = *req.BudgetAllocated
	}

	// Convert to model and create
	caseModel := models.NewCaseFromEntity(caseEntity)
	caseID, err := uc.caseRepo.Create(ctx, *caseModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create case: %w", err)
	}

	// If initial note provided, create it
	if req.InitialNote != nil {
		note := entities.CaseNote{
			CaseID:      caseID,
			Title:       req.InitialNote.Title,
			Content:     req.InitialNote.Content,
			Tags:        req.InitialNote.Tags,
			NoteType:    req.InitialNote.NoteType,
			IsImportant: req.InitialNote.IsImportant,
			CreatedByID: assignedUser.ID,
			CreatedBy:   *assignedUser,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		noteModel := models.NewCaseNoteFromEntity(note)
		_, err = uc.noteRepo.Create(ctx, *noteModel)
		if err != nil {
			// Log error but don't fail case creation
			fmt.Printf("Warning: failed to create initial note: %v\n", err)
		} else {
			// Increment notes count
			uc.caseRepo.IncrementNotesCount(ctx, caseID)
		}
	}

	// Get the created case with populated relations
	createdCase, err := uc.GetByID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created case: %w", err)
	}

	return createdCase, nil
}

func (uc *CaseUseCase) GetByID(ctx context.Context, id string) (*responses.CaseResponse, error) {
	// Get case from repository
	caseModel, err := uc.caseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get beneficiary details
	beneficiary, err := uc.beneficiaryUC.GetByID(ctx, caseModel.BeneficiaryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get beneficiary: %w", err)
	}

	// Get assigned user details
	assignedUser, err := uc.userUC.GetByID(ctx, caseModel.AssignedToID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned user: %w", err)
	}

	// Convert to entity with populated relations
	caseEntity := caseModel.ToEntity()
	caseEntity.Beneficiary = *beneficiary
	caseEntity.AssignedTo = *assignedUser

	// Convert to response
	response := responses.NewCaseResponse(caseEntity)
	return &response, nil
}

func (uc *CaseUseCase) Update(ctx context.Context, id string, req requests.UpdateCaseRequest) (*responses.CaseResponse, error) {
	// Build updates map
	updates := bson.M{}

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.UrgencyLevel != nil {
		updates["urgency_level"] = *req.UrgencyLevel
	}
	if req.CaseType != nil {
		updates["case_type"] = *req.CaseType
	}
	if req.AssignedToID != nil {
		// Validate user exists
		_, err := uc.userUC.GetByID(ctx, *req.AssignedToID)
		if err != nil {
			return nil, fmt.Errorf("assigned user not found: %w", err)
		}
		updates["assigned_to_id"] = *req.AssignedToID
	}
	if req.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due date format: %w", err)
		}
		updates["due_date"] = dueDate
	}
	if req.EstimatedDuration != nil {
		updates["estimated_duration"] = *req.EstimatedDuration
	}
	if req.BudgetAllocated != nil {
		updates["budget_allocated"] = *req.BudgetAllocated
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	// Update last activity
	updates["last_activity"] = time.Now()

	// Perform update
	err := uc.caseRepo.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	// Return updated case
	return uc.GetByID(ctx, id)
}

func (uc *CaseUseCase) Delete(ctx context.Context, id string) error {
	return uc.caseRepo.Delete(ctx, id)
}

func (uc *CaseUseCase) ListByOrganization(ctx context.Context, organizationID string, filters repositories.CaseFilters) (*responses.CaseListResponse, error) {
	// Set organization ID in filters
	filters.OrganizationID = organizationID

	// Get cases from repository
	caseModels, total, err := uc.caseRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Convert to responses with populated relations
	caseResponses := make([]responses.CaseResponse, len(caseModels))
	for i, caseModel := range caseModels {
		// Get beneficiary details
		beneficiary, err := uc.beneficiaryUC.GetByID(ctx, caseModel.BeneficiaryID)
		if err != nil {
			// Handle missing beneficiary gracefully
			beneficiary = &entities.Beneficiary{
				ID:       caseModel.BeneficiaryID,
				FullName: "Unknown Beneficiary",
			}
		}

		// Get assigned user details
		assignedUser, err := uc.userUC.GetByID(ctx, caseModel.AssignedToID)
		if err != nil {
			// Handle missing user gracefully
			assignedUser = &entities.User{
				ID:        caseModel.AssignedToID,
				FirstName: "Unknown",
				LastName:  "User",
				FullName:  "Unknown User",
				Email:     "unknown@example.com",
			}
		}

		// Convert to entity with populated relations
		caseEntity := caseModel.ToEntity()
		caseEntity.Beneficiary = *beneficiary
		caseEntity.AssignedTo = *assignedUser

		// Convert to response
		caseResponses[i] = responses.NewCaseResponse(caseEntity)
	}

	return &responses.CaseListResponse{
		Count: total,
		Data:  caseResponses,
	}, nil
}

func (uc *CaseUseCase) GetStats(ctx context.Context, organizationID string) (*responses.CaseStatsResponse, error) {
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
