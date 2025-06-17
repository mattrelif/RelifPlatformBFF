package models

import (
	"relif/platform-bff/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BeneficiaryStatusAudit struct {
	ID             string    `bson:"_id,omitempty"`
	BeneficiaryID  string    `bson:"beneficiary_id,omitempty"`
	PreviousStatus string    `bson:"previous_status,omitempty"`
	NewStatus      string    `bson:"new_status,omitempty"`
	ChangedBy      string    `bson:"changed_by,omitempty"`
	ChangedAt      time.Time `bson:"changed_at,omitempty"`
	Reason         string    `bson:"reason,omitempty"`
	OrganizationID string    `bson:"organization_id,omitempty"`
}

func (audit *BeneficiaryStatusAudit) ToEntity() entities.BeneficiaryStatusAudit {
	return entities.BeneficiaryStatusAudit{
		ID:             audit.ID,
		BeneficiaryID:  audit.BeneficiaryID,
		PreviousStatus: audit.PreviousStatus,
		NewStatus:      audit.NewStatus,
		ChangedBy:      audit.ChangedBy,
		ChangedAt:      audit.ChangedAt,
		Reason:         audit.Reason,
		OrganizationID: audit.OrganizationID,
	}
}

func NewBeneficiaryStatusAudit(entity entities.BeneficiaryStatusAudit) BeneficiaryStatusAudit {
	return BeneficiaryStatusAudit{
		ID:             primitive.NewObjectID().Hex(),
		BeneficiaryID:  entity.BeneficiaryID,
		PreviousStatus: entity.PreviousStatus,
		NewStatus:      entity.NewStatus,
		ChangedBy:      entity.ChangedBy,
		ChangedAt:      time.Now(),
		Reason:         entity.Reason,
		OrganizationID: entity.OrganizationID,
	}
}
