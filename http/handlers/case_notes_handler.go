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

	"github.com/go-chi/chi/v5"
)

type CaseNotes struct {
	noteRepo *repositories.CaseNoteRepository
	caseRepo repositories.CaseRepository
}

func NewCaseNotes(noteRepo *repositories.CaseNoteRepository, caseRepo repositories.CaseRepository) *CaseNotes {
	return &CaseNotes{
		noteRepo: noteRepo,
		caseRepo: caseRepo,
	}
}

// GET /api/cases/{case_id}/notes
func (h *CaseNotes) ListByCaseID(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	_ = r.Context().Value("user").(entities.User) // TODO: Use for authorization

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
		// Convert model to entity (you'd need this method)
		noteEntity := note.ToEntity() // This method needs to be implemented
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

	// Create note entity
	noteEntity := entities.CaseNote{
		CaseID:      caseID,
		Title:       req.Title,
		Content:     req.Content,
		Tags:        req.Tags,
		NoteType:    req.NoteType,
		IsImportant: req.IsImportant,
		CreatedByID: user.ID,
		CreatedBy:   user,
	}

	// Convert to model and create
	noteModel := models.NewCaseNoteFromEntity(noteEntity) // This needs to be implemented
	noteID, err := h.noteRepo.Create(r.Context(), *noteModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Increment case notes count
	h.caseRepo.UpdateNotesCount(r.Context(), caseID, 1)

	// Get created note
	createdNote, err := h.noteRepo.GetByID(r.Context(), noteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	noteEntity = createdNote.ToEntity()
	response := responses.NewCaseNoteResponse(noteEntity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// PUT /api/cases/{case_id}/notes/{note_id}
func (h *CaseNotes) Update(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "note_id")
	_ = r.Context().Value("user").(entities.User) // TODO: Use for authorization

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

	// Build updates map
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.NoteType != nil {
		updates["note_type"] = *req.NoteType
	}
	if req.IsImportant != nil {
		updates["is_important"] = *req.IsImportant
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	err = h.noteRepo.Update(r.Context(), noteID, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get updated note
	updatedNote, err := h.noteRepo.GetByID(r.Context(), noteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	noteEntity := updatedNote.ToEntity()
	response := responses.NewCaseNoteResponse(noteEntity)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DELETE /api/cases/{case_id}/notes/{note_id}
func (h *CaseNotes) Delete(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	noteID := chi.URLParam(r, "note_id")
	_ = r.Context().Value("user").(entities.User) // TODO: Use for authorization

	err := h.noteRepo.Delete(r.Context(), noteID)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Decrement case notes count
	h.caseRepo.UpdateNotesCount(r.Context(), caseID, -1)

	w.WriteHeader(http.StatusNoContent)
}
