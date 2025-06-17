package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/repositories"
	casesUseCases "relif/platform-bff/usecases/cases"

	"github.com/go-chi/chi/v5"
)

type CaseDocuments struct {
	docRepo                           *repositories.CaseDocumentRepository
	caseRepo                          repositories.CaseRepository
	userRepo                          repositories.Users
	createDocumentUseCase             casesUseCases.CreateDocumentUseCase
	updateDocumentUseCase             casesUseCases.UpdateDocumentUseCase
	generateDocumentUploadLinkUseCase casesUseCases.GenerateDocumentUploadLink
}

func NewCaseDocuments(
	docRepo *repositories.CaseDocumentRepository,
	caseRepo repositories.CaseRepository,
	userRepo repositories.Users,
	createDocumentUseCase casesUseCases.CreateDocumentUseCase,
	updateDocumentUseCase casesUseCases.UpdateDocumentUseCase,
	generateDocumentUploadLinkUseCase casesUseCases.GenerateDocumentUploadLink,
) *CaseDocuments {
	return &CaseDocuments{
		docRepo:                           docRepo,
		caseRepo:                          caseRepo,
		userRepo:                          userRepo,
		createDocumentUseCase:             createDocumentUseCase,
		updateDocumentUseCase:             updateDocumentUseCase,
		generateDocumentUploadLinkUseCase: generateDocumentUploadLinkUseCase,
	}
}

// GET /api/cases/{case_id}/documents
func (h *CaseDocuments) ListByCaseID(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	user := r.Context().Value("user").(entities.User)

	// Check authorization - get case and verify user has access to its organization
	caseEntity, err := h.caseRepo.GetByID(r.Context(), caseID)
	if err != nil {
		http.Error(w, "Case not found", http.StatusNotFound)
		return
	}

	if err := guards.IsOrganizationAdmin(user, caseEntity.Organization); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

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
		// Populate UploadedBy user information
		if docEntity.UploadedByID != "" {
			if user, err := h.userRepo.FindOneByID(docEntity.UploadedByID); err == nil {
				docEntity.UploadedBy = user
			}
		}
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

	docEntity := req.ToEntity() // Convert request to entity
	createdDoc, err := h.createDocumentUseCase.Execute(r.Context(), user, caseID, docEntity)
	if err != nil {
		// The use case now handles errors, we just map them
		// This part can be improved with custom error types like in the beneficiary handler
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.NewCaseDocumentResponse(createdDoc)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DELETE /api/cases/{case_id}/documents/{doc_id}
func (h *CaseDocuments) Delete(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	docID := chi.URLParam(r, "doc_id")
	user := r.Context().Value("user").(entities.User)

	// Check authorization - get case and verify user has access to its organization
	caseEntity, err := h.caseRepo.GetByID(r.Context(), caseID)
	if err != nil {
		http.Error(w, "Case not found", http.StatusNotFound)
		return
	}

	if err := guards.IsOrganizationAdmin(user, caseEntity.Organization); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete from database
	err = h.docRepo.Delete(r.Context(), docID)
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

	caseID := chi.URLParam(r, "case_id")
	docID := chi.URLParam(r, "doc_id")
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

	updateEntity := req.ToEntity()

	updatedDoc, err := h.updateDocumentUseCase.Execute(r.Context(), user, caseID, docID, updateEntity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.NewCaseDocumentResponse(updatedDoc)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
