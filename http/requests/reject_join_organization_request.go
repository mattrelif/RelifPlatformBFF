package requests

import (
	"relif/platform-bff/entities"
)

type RejectJoinOrganizationRequest struct {
	RejectReason string `json:"reject_reason"`
}

func (req *RejectJoinOrganizationRequest) ToEntity() entities.JoinOrganizationRequest {
	return entities.JoinOrganizationRequest{
		RejectReason: req.RejectReason,
	}
}
