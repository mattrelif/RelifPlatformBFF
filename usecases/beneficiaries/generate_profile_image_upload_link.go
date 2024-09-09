package beneficiaries

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	filesUseCases "relif/platform-bff/usecases/files"
)

type GenerateProfileImageUploadLink interface {
	Execute(actor entities.User, fileType string) (string, error)
}

type generateProfileImageUploadLinkImpl struct {
	generateUploadLinkUseCase filesUseCases.GenerateUploadLink
}

func NewGenerateProfileImageUploadLink(generateUploadLinkUseCase filesUseCases.GenerateUploadLink) GenerateProfileImageUploadLink {
	return &generateProfileImageUploadLinkImpl{
		generateUploadLinkUseCase: generateUploadLinkUseCase,
	}
}

func (uc *generateProfileImageUploadLinkImpl) Execute(actor entities.User, fileType string) (string, error) {
	if err := guards.IsOrganizationAdmin(actor, actor.Organization); err != nil {
		return "", err
	}

	return uc.generateUploadLinkUseCase.Execute("beneficiaries", fileType)
}
