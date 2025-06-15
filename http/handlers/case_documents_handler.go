package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/models"
	"relif/platform-bff/repositories"

	"github.com/go-chi/chi/v5"
)

type CaseDocuments struct {
	docRepo  *repositories.CaseDocumentRepository
	caseRepo repositories.CaseRepository
}

func NewCaseDocuments(docRepo *repositories.CaseDocumentRepository, caseRepo repositories.CaseRepository) *CaseDocuments {
	return &CaseDocuments{
		docRepo:  docRepo,
		caseRepo: caseRepo,
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

// POST /api/cases/{case_id}/documents
func (h *CaseDocuments) Upload(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	user := r.Context().Value("user").(entities.User)

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Create request from form data
	req := requests.CreateCaseDocumentRequest{
		DocumentName: r.FormValue("document_name"),
		DocumentType: r.FormValue("document_type"),
		Description:  r.FormValue("description"),
		FileName:     fileHeader.Filename,
		FileSize:     fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		FileData:     fileData,
	}

	// Parse tags if provided
	if tagsStr := r.FormValue("tags"); tagsStr != "" {
		var tags []string
		if err := json.Unmarshal([]byte(tagsStr), &tags); err == nil {
			req.Tags = tags
		}
	}

	// Validate request
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Upload file to storage service (S3, etc.) and get download URL
	downloadURL := fmt.Sprintf("/api/files/%s", fileHeader.Filename) // Placeholder

	// Create document entity
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
		DownloadURL:  downloadURL,
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

	err := h.docRepo.Delete(r.Context(), docID)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Decrement case documents count
	h.caseRepo.UpdateDocumentsCount(r.Context(), caseID, -1)

	w.WriteHeader(http.StatusNoContent)
}
