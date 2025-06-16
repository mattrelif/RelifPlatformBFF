package cases

import (
	"relif/platform-bff/entities"
	filesUseCases "relif/platform-bff/usecases/files"
)

type GenerateDocumentUploadLink interface {
	Execute(user entities.User, caseID, fileType string) (string, error)
}

type generateDocumentUploadLinkImpl struct {
	generateUploadLinkUseCase filesUseCases.GenerateUploadLink
}

func NewGenerateDocumentUploadLink(generateUploadLinkUseCase filesUseCases.GenerateUploadLink) GenerateDocumentUploadLink {
	return &generateDocumentUploadLinkImpl{
		generateUploadLinkUseCase: generateUploadLinkUseCase,
	}
}

func (uc *generateDocumentUploadLinkImpl) Execute(user entities.User, caseID, fileType string) (string, error) {
	// Use "case-documents" as the domain with case ID for better organization
	domain := "case-documents/" + caseID
	return uc.generateUploadLinkUseCase.Execute(domain, fileType)
}
