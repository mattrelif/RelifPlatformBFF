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
	productTypeAllocationsUseCases "relif/platform-bff/usecases/product_type_alloctions"
	"relif/platform-bff/utils"
	"strconv"
)

type ProductTypeAllocations struct {
	createEntranceUseCase                   productTypeAllocationsUseCases.CreateEntrance
	createReallocationUseCase               productTypeAllocationsUseCases.CreateReallocation
	findManyByProductTypeIDPaginatedUseCase productTypeAllocationsUseCases.FindManyByProductTypeIDPaginated
}

func NewProductTypeAllocations(
	createEntranceUseCase productTypeAllocationsUseCases.CreateEntrance,
	createReallocationUseCase productTypeAllocationsUseCases.CreateReallocation,
	findManyByProductTypeIDPaginatedUseCase productTypeAllocationsUseCases.FindManyByProductTypeIDPaginated,
) *ProductTypeAllocations {
	return &ProductTypeAllocations{
		createEntranceUseCase:                   createEntranceUseCase,
		createReallocationUseCase:               createReallocationUseCase,
		findManyByProductTypeIDPaginatedUseCase: findManyByProductTypeIDPaginatedUseCase,
	}
}

func (handler *ProductTypeAllocations) Allocate(w http.ResponseWriter, r *http.Request) {
	var req requests.AllocateProductType

	productTypeID := chi.URLParam(r, "id")
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
	allocation, err := handler.createEntranceUseCase.Execute(user, productTypeID, data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrProductTypeNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewProductTypeAllocation(allocation)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *ProductTypeAllocations) Reallocate(w http.ResponseWriter, r *http.Request) {
	var req requests.ReallocateProductType

	productTypeID := chi.URLParam(r, "id")
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
	allocation, err := handler.createReallocationUseCase.Execute(user, productTypeID, data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrProductTypeNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewProductTypeAllocation(allocation)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *ProductTypeAllocations) FindManyByProductTypeID(w http.ResponseWriter, r *http.Request) {
	productTypeID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limitParam := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	count, allocations, err := handler.findManyByProductTypeIDPaginatedUseCase.Execute(user, productTypeID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrBeneficiaryNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.ProductTypeAllocations]{Data: responses.NewProductTypeAllocations(allocations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
