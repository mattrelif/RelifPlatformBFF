package beneficiaries

import (
	"errors"
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type UpdateStatus interface {
	Execute(actor entities.User, beneficiaryID, newStatus string) (entities.Beneficiary, error)
}

type updateStatusImpl struct {
	beneficiariesRepository          repositories.Beneficiaries
	organizationsRepository          repositories.Organizations
	beneficiaryStatusAuditRepository repositories.BeneficiaryStatusAudit
}

func NewUpdateStatus(
	beneficiariesRepository repositories.Beneficiaries,
	organizationsRepository repositories.Organizations,
	beneficiaryStatusAuditRepository repositories.BeneficiaryStatusAudit,
) UpdateStatus {
	return &updateStatusImpl{
		beneficiariesRepository:          beneficiariesRepository,
		organizationsRepository:          organizationsRepository,
		beneficiaryStatusAuditRepository: beneficiaryStatusAuditRepository,
	}
}

func (uc *updateStatusImpl) Execute(actor entities.User, beneficiaryID, newStatus string) (entities.Beneficiary, error) {
	// Find the beneficiary
	beneficiary, err := uc.beneficiariesRepository.FindOneByID(beneficiaryID)
	if err != nil {
		return entities.Beneficiary{}, utils.ErrBeneficiaryNotFound
	}

	// Check permissions
	organization, err := uc.organizationsRepository.FindOneByID(beneficiary.CurrentOrganizationID)
	if err != nil {
		return entities.Beneficiary{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.Beneficiary{}, utils.ErrForbiddenAction
	}

	// Validate status values
	if !isValidStatus(newStatus) {
		return entities.Beneficiary{}, errors.New("INVALID_STATUS")
	}

	// Check status transition rules
	if !isValidStatusTransition(beneficiary.Status, newStatus) {
		return entities.Beneficiary{}, errors.New("INVALID_STATUS_TRANSITION")
	}

	// Store previous status for audit
	previousStatus := beneficiary.Status

	// Update beneficiary status
	beneficiary.Status = newStatus
	beneficiary.UpdatedAt = time.Now()

	if err = uc.beneficiariesRepository.UpdateOneByID(beneficiary.ID, beneficiary); err != nil {
		return entities.Beneficiary{}, err
	}

	// Create audit record
	auditRecord := entities.BeneficiaryStatusAudit{
		BeneficiaryID:  beneficiary.ID,
		PreviousStatus: previousStatus,
		NewStatus:      newStatus,
		ChangedBy:      actor.ID,
		ChangedAt:      time.Now(),
		OrganizationID: beneficiary.CurrentOrganizationID,
	}

	if err = uc.beneficiaryStatusAuditRepository.Create(auditRecord); err != nil {
		// Log error but don't fail the operation
		// In production, you might want to implement proper logging
	}

	return beneficiary, nil
}

func isValidStatus(status string) bool {
	validStatuses := []string{
		utils.ActiveStatus,
		utils.InactiveStatus,
		utils.PendingStatus,
		utils.ArchivedStatus,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}

	return false
}

func isValidStatusTransition(currentStatus, newStatus string) bool {
	// ARCHIVED status cannot be changed to any other status
	if currentStatus == utils.ArchivedStatus {
		return false
	}

	// All other transitions are allowed based on business rules
	return true
}
