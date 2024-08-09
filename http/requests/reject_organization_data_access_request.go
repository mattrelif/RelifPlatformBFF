package requests

import (
	"relif/platform-bff/entities"
)

type RejectOrganizationDataAccessRequest struct {
	RejectReason string `json:"reject_reason"`
}

func (req *RejectOrganizationDataAccessRequest) ToEntity() entities.OrganizationDataAccessRequest {
	return entities.OrganizationDataAccessRequest{
		RejectReason: req.RejectReason,
	}
}
