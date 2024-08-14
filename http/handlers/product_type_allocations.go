package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
)

type ProductTypeAllocations struct {
	service              services.ProductTypeAllocations
	authorizationService services.Authorization
}

func NewProductTypeAllocations(service services.ProductTypeAllocations, authorizationService services.Authorization) *ProductTypeAllocations {
	return &ProductTypeAllocations{
		service:              service,
		authorizationService: authorizationService,
	}
}

func (handler *ProductTypeAllocations) Allocate(w http.ResponseWriter, r *http.Request) {
	var req requests.AllocateProductType

	productTypeID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeCreateProductTypeResource(productTypeID, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrProductTypeNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

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

	allocation, err := handler.service.CreateEntrance(productTypeID, req.ToEntity())

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrProductTypeNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewProductTypeAllocation(allocation)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *ProductTypeAllocations) Reallocate(w http.ResponseWriter, r *http.Request) {
	var req requests.ReallocateProductType

	productTypeID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeCreateProductTypeResource(productTypeID, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrProductTypeNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

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

	allocation, err := handler.service.CreateReallocation(productTypeID, req.ToEntity())

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrProductTypeNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewProductTypeAllocation(allocation)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
