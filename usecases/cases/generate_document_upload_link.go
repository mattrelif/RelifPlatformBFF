package cases

import (
	"context"
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	filesUseCases "relif/platform-bff/usecases/files"
)

type GenerateDocumentUploadLink interface {
	Execute(user entities.User, caseID, fileType string) (string, error)
}

type generateDocumentUploadLinkImpl struct {
	generateUploadLinkUseCase filesUseCases.GenerateUploadLink
	casesRepository           repositories.CaseRepository
}

func NewGenerateDocumentUploadLink(generateUploadLinkUseCase filesUseCases.GenerateUploadLink, casesRepository repositories.CaseRepository) GenerateDocumentUploadLink {
	return &generateDocumentUploadLinkImpl{
		generateUploadLinkUseCase: generateUploadLinkUseCase,
		casesRepository:           casesRepository,
	}
}

func (uc *generateDocumentUploadLinkImpl) Execute(user entities.User, caseID, fileType string) (string, error) {
	// First, get the case to check authorization
	caseEntity, err := uc.casesRepository.GetByID(context.Background(), caseID)
	if err != nil {
		return "", err
	}

	// Check if user has access to this case's organization
	if err := guards.IsOrganizationAdmin(user, caseEntity.Organization); err != nil {
		return "", err
	}

	// Use "case-documents" as the domain with case ID for better organization
	domain := "case-documents/" + caseID
	return uc.generateUploadLinkUseCase.Execute(domain, fileType)
}
