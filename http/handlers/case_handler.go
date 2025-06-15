package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/repositories"
	"relif/platform-bff/usecases/cases"

	"github.com/go-chi/chi/v5"
)

type Cases struct {
	caseUC cases.CaseUseCase
}

func NewCases(caseUC cases.CaseUseCase) *Cases {
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

	// Build filters from query parameters
	filters := repositories.CaseFilters{
		Status:       r.URL.Query().Get("status"),
		Priority:     r.URL.Query().Get("priority"),
		CaseType:     r.URL.Query().Get("case_type"),
		AssignedToID: r.URL.Query().Get("assigned_to_id"),
		Search:       r.URL.Query().Get("search"),
		SortBy:       r.URL.Query().Get("sort_by"),
		SortOrder:    r.URL.Query().Get("sort_order"),
	}

	if filters.SortBy == "" {
		filters.SortBy = "created_at"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "desc"
	}

	// Parse pagination
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

	cases, total, err := h.caseUC.ListCases(r.Context(), user.OrganizationID, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[[]responses.CaseResponse]{
		Data:  cases,
		Count: total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GET /api/cases/:id
func (h *Cases) FindOne(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	result, err := h.caseUC.GetCase(r.Context(), caseID, user.OrganizationID)
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
