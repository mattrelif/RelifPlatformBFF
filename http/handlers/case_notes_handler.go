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

type CaseNotes struct {
	noteRepo *repositories.CaseNoteRepository
	caseRepo repositories.CaseRepository
	userRepo repositories.Users

	createNoteUseCase casesUseCases.CreateNoteUseCase
	updateNoteUseCase casesUseCases.UpdateNoteUseCase
}

func NewCaseNotes(
	noteRepo *repositories.CaseNoteRepository,
	caseRepo repositories.CaseRepository,
	userRepo repositories.Users,
	createNoteUseCase casesUseCases.CreateNoteUseCase,
	updateNoteUseCase casesUseCases.UpdateNoteUseCase,
) *CaseNotes {
	return &CaseNotes{
		noteRepo:          noteRepo,
		caseRepo:          caseRepo,
		userRepo:          userRepo,
		createNoteUseCase: createNoteUseCase,
		updateNoteUseCase: updateNoteUseCase,
	}
}

// GET /api/cases/{case_id}/notes
func (h *CaseNotes) ListByCaseID(w http.ResponseWriter, r *http.Request) {
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
	filters := repositories.CaseNoteFilters{
		CaseID: caseID,
	}

	// Parse optional filters
	if noteType := r.URL.Query().Get("note_type"); noteType != "" {
		filters.NoteType = &noteType
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filters.Search = &search
	}
	if important := r.URL.Query().Get("is_important"); important != "" {
		if important == "true" {
			isImportant := true
			filters.Important = &isImportant
		} else if important == "false" {
			isImportant := false
			filters.Important = &isImportant
		}
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

	notes, total, err := h.noteRepo.ListByCaseID(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	noteResponses := make([]responses.CaseNoteResponse, len(notes))
	for i, note := range notes {
		noteEntity := note.ToEntity()

		// Populate CreatedBy user information
		if noteEntity.CreatedByID != "" {
			createdBy, err := h.userRepo.FindOneByID(noteEntity.CreatedByID)
			if err == nil {
				noteEntity.CreatedBy = createdBy
			}
		}

		noteResponses[i] = responses.NewCaseNoteResponse(noteEntity)
	}

	response := responses.CaseNoteListResponse{
		Count: total,
		Data:  noteResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /api/cases/{case_id}/notes
func (h *CaseNotes) Create(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	user := r.Context().Value("user").(entities.User)

	var req requests.CreateCaseNoteRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Create note entity
	noteEntity := entities.CaseNote{
		CaseID:      caseID,
		Title:       req.Title,
		Content:     req.Content,
		Tags:        req.Tags,
		NoteType:    req.NoteType,
		IsImportant: req.IsImportant,
	}

	createdNote, err := h.createNoteUseCase.Execute(r.Context(), user, caseID, noteEntity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.NewCaseNoteResponse(createdNote)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// PUT /api/cases/{case_id}/notes/{note_id}
func (h *CaseNotes) Update(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	noteID := chi.URLParam(r, "note_id")
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

	var req requests.UpdateCaseNoteRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build updates entity based on optional fields
	updatesEntity := entities.CaseNote{}
	if req.Title != nil {
		updatesEntity.Title = *req.Title
	}
	if req.Content != nil {
		updatesEntity.Content = *req.Content
	}
	if req.NoteType != nil {
		updatesEntity.NoteType = *req.NoteType
	}
	if req.IsImportant != nil {
		updatesEntity.IsImportant = *req.IsImportant
	}
	if req.Tags != nil {
		updatesEntity.Tags = req.Tags
	}

	updatedNote, err := h.updateNoteUseCase.Execute(r.Context(), user, caseID, noteID, updatesEntity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.NewCaseNoteResponse(updatedNote)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DELETE /api/cases/{case_id}/notes/{note_id}
func (h *CaseNotes) Delete(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	noteID := chi.URLParam(r, "note_id")
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

	err = h.noteRepo.Delete(r.Context(), noteID)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Decrement case notes count
	h.caseRepo.UpdateNotesCount(r.Context(), caseID, -1)

	w.WriteHeader(http.StatusNoContent)
}
