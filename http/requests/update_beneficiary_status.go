package requests

import (
	"relif/platform-bff/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateBeneficiaryStatus struct {
	Status string `json:"status"`
}

func (req *UpdateBeneficiaryStatus) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Status, validation.Required, validation.In(
			utils.ActiveStatus,
			utils.InactiveStatus,
			utils.PendingStatus,
			utils.ArchivedStatus,
		)),
	)
}

func (req *UpdateBeneficiaryStatus) IsValidStatusTransition(currentStatus string) bool {
	// ARCHIVED status cannot be changed to any other status
	if currentStatus == utils.ArchivedStatus {
		return false
	}

	// All other transitions are allowed based on business rules
	return true
}
