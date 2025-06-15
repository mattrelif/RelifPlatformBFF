package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/repositories"
	"relif/platform-bff/usecases/cases"

	"github.com/go-chi/chi/v5"
)

type Cases struct {
	caseUC *cases.CaseUseCase
}

func NewCases(caseUC *cases.CaseUseCase) *Cases {
	return &Cases{
		caseUC: caseUC,
	}
}

// POST /api/cases
func (h *Cases) CreateCase(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateCase
	user := r.Context().Value("user").(entities.User)

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

	result, err := h.caseUC.CreateCase(r.Context(), req, user.ID, user.OrganizationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// GET /api/cases
func (h *Cases) FindManyByOrganization(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entities.User)

	// Helper function to convert string to *string if not empty
	stringPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	// Build filters from query parameters
	filters := repositories.CaseFilters{
		Status:       stringPtr(r.URL.Query().Get("status")),
		Priority:     stringPtr(r.URL.Query().Get("priority")),
		CaseType:     stringPtr(r.URL.Query().Get("case_type")),
		AssignedToID: stringPtr(r.URL.Query().Get("assigned_to_id")),
		Search:       stringPtr(r.URL.Query().Get("search")),
		SortBy:       r.URL.Query().Get("sort_by"),
		SortOrder:    r.URL.Query().Get("sort_order"),
	}

	// Handle organization_id from query params (for frontend compatibility)
	// But still use authenticated user's organization for security
	if orgIDParam := r.URL.Query().Get("organization_id"); orgIDParam != "" {
		if orgIDParam != user.OrganizationID {
			http.Error(w, "Access denied: organization mismatch", http.StatusForbidden)
			return
		}
	}

	if filters.SortBy == "" {
		filters.SortBy = "created_at"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "desc"
	}

	// Parse pagination (handle both page/limit and offset/limit patterns)
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filters.Page = page
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	// Handle offset-based pagination (frontend compatibility)
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && filters.Limit > 0 {
			filters.Page = (offset / filters.Limit) + 1
		}
	}

	result, err := h.caseUC.ListByOrganization(r.Context(), user.OrganizationID, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET /api/cases/:id
func (h *Cases) FindOne(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "id")

	result, err := h.caseUC.GetByID(r.Context(), caseID)
	if err != nil {
		http.Error(w, "Case not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// PUT /api/cases/:id
func (h *Cases) UpdateOne(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	var req requests.UpdateCase

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

	// ADD MISSING VALIDATION!
	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.caseUC.UpdateCase(r.Context(), caseID, req, user.OrganizationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DELETE /api/cases/:id
func (h *Cases) DeleteOne(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	err := h.caseUC.DeleteCase(r.Context(), caseID, user.OrganizationID)
	if err != nil {
		http.Error(w, "Case not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/cases/stats
func (h *Cases) GetStats(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entities.User)

	stats, err := h.caseUC.GetCaseStats(r.Context(), user.OrganizationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
