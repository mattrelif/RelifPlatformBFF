package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/models"
	"relif/platform-bff/repositories"
	casesUseCases "relif/platform-bff/usecases/cases"

	"github.com/go-chi/chi/v5"
)

type CaseDocuments struct {
	docRepo                           *repositories.CaseDocumentRepository
	caseRepo                          repositories.CaseRepository
	generateDocumentUploadLinkUseCase casesUseCases.GenerateDocumentUploadLink
}

func NewCaseDocuments(
	docRepo *repositories.CaseDocumentRepository,
	caseRepo repositories.CaseRepository,
	generateDocumentUploadLinkUseCase casesUseCases.GenerateDocumentUploadLink,
) *CaseDocuments {
	return &CaseDocuments{
		docRepo:                           docRepo,
		caseRepo:                          caseRepo,
		generateDocumentUploadLinkUseCase: generateDocumentUploadLinkUseCase,
	}
}

// GET /api/cases/{case_id}/documents
func (h *CaseDocuments) ListByCaseID(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	_ = r.Context().Value("user").(entities.User) // TODO: Use for authorization

	// Build filters from query parameters
	filters := repositories.CaseDocumentFilters{
		CaseID: caseID,
	}

	// Parse optional filters
	if docType := r.URL.Query().Get("document_type"); docType != "" {
		filters.DocumentType = &docType
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filters.Search = &search
	}
	if uploadedBy := r.URL.Query().Get("uploaded_by"); uploadedBy != "" {
		filters.UploadedBy = &uploadedBy
	}

	// Parse pagination
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	// Default limit
	if filters.Limit == 0 {
		filters.Limit = 20
	}

	documents, total, err := h.docRepo.ListByCaseID(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	docResponses := make([]responses.CaseDocumentResponse, len(documents))
	for i, doc := range documents {
		docEntity := doc.ToEntity()
		docResponses[i] = responses.NewCaseDocumentResponse(docEntity)
	}

	response := responses.CaseDocumentListResponse{
		Count: total,
		Data:  docResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /api/cases/{case_id}/documents/generate-upload-link
func (h *CaseDocuments) GenerateUploadLink(w http.ResponseWriter, r *http.Request) {
	var req requests.GenerateFileUploadLink

	caseID := chi.URLParam(r, "case_id")
	user := r.Context().Value("user").(entities.User)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := req.ToEntity()
	link, err := h.generateDocumentUploadLinkUseCase.Execute(user, caseID, data.Type)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewGenerateFileUploadLink(link)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /api/cases/{case_id}/documents
func (h *CaseDocuments) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateCaseDocumentRequest

	caseID := chi.URLParam(r, "case_id")
	user := r.Context().Value("user").(entities.User)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create document entity - S3 URL will be constructed from file key
	docEntity := entities.CaseDocument{
		CaseID:       caseID,
		DocumentName: req.DocumentName,
		FileName:     req.FileName,
		DocumentType: req.DocumentType,
		FileSize:     req.FileSize,
		MimeType:     req.MimeType,
		Description:  req.Description,
		Tags:         req.Tags,
		UploadedByID: user.ID,
		UploadedBy:   user,
		FilePath:     req.FileKey, // S3 object key
	}

	// Convert to model and create
	docModel := models.NewCaseDocumentFromEntity(docEntity)
	docID, err := h.docRepo.Create(r.Context(), *docModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Increment case documents count
	h.caseRepo.UpdateDocumentsCount(r.Context(), caseID, 1)

	// Get created document
	createdDoc, err := h.docRepo.GetByID(r.Context(), docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	docEntity = createdDoc.ToEntity()
	response := responses.NewCaseDocumentResponse(docEntity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DELETE /api/cases/{case_id}/documents/{doc_id}
func (h *CaseDocuments) Delete(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	docID := chi.URLParam(r, "doc_id")
	_ = r.Context().Value("user").(entities.User) // TODO: Use for authorization

	// Delete from database
	err := h.docRepo.Delete(r.Context(), docID)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Decrement case documents count
	h.caseRepo.UpdateDocumentsCount(r.Context(), caseID, -1)

	w.WriteHeader(http.StatusNoContent)
}

// PUT /api/cases/{case_id}/documents/{doc_id}
func (h *CaseDocuments) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateCaseDocumentRequest

	docID := chi.URLParam(r, "doc_id")
	_ = r.Context().Value("user").(entities.User) // TODO: Use for authorization

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to entity and then model for update
	updateEntity := entities.CaseDocument{
		DocumentName: req.DocumentName,
		DocumentType: req.DocumentType,
		Description:  req.Description,
		Tags:         req.Tags,
	}

	updateModel := models.NewCaseDocumentFromEntity(updateEntity)
	err = h.docRepo.Update(r.Context(), docID, *updateModel)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
