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
	beneficiaryAllocationsUseCases "relif/platform-bff/usecases/beneficiary_allocations"
	"relif/platform-bff/utils"
	"strconv"
)

type BeneficiaryAllocations struct {
	createEntranceUseCase                   beneficiaryAllocationsUseCases.CreateEntrance
	createReallocationUseCase               beneficiaryAllocationsUseCases.CreateReallocation
	findManyByBeneficiaryIDPaginatedUseCase beneficiaryAllocationsUseCases.FindManyByBeneficiaryIDPaginated
	findManyByHousingIDPaginatedUseCase     beneficiaryAllocationsUseCases.FindManyByHousingIDPaginated
	findManyByHousingRoomIDPaginatedUseCase beneficiaryAllocationsUseCases.FindManyByHousingRoomIDPaginated
}

func NewBeneficiaryAllocations(
	createEntranceUseCase beneficiaryAllocationsUseCases.CreateEntrance,
	createReallocationUseCase beneficiaryAllocationsUseCases.CreateReallocation,
	findManyByBeneficiaryIDPaginatedUseCase beneficiaryAllocationsUseCases.FindManyByBeneficiaryIDPaginated,
	findManyByHousingIDPaginatedUseCase beneficiaryAllocationsUseCases.FindManyByHousingIDPaginated,
	findManyByHousingRoomIDPaginatedUseCase beneficiaryAllocationsUseCases.FindManyByHousingRoomIDPaginated,
) *BeneficiaryAllocations {
	return &BeneficiaryAllocations{
		createEntranceUseCase:                   createEntranceUseCase,
		createReallocationUseCase:               createReallocationUseCase,
		findManyByBeneficiaryIDPaginatedUseCase: findManyByBeneficiaryIDPaginatedUseCase,
		findManyByHousingIDPaginatedUseCase:     findManyByHousingIDPaginatedUseCase,
		findManyByHousingRoomIDPaginatedUseCase: findManyByHousingRoomIDPaginatedUseCase,
	}
}

func (handler *BeneficiaryAllocations) Allocate(w http.ResponseWriter, r *http.Request) {
	var req requests.AllocateBeneficiary

	user := r.Context().Value("user").(entities.User)
	beneficiaryID := chi.URLParam(r, "id")

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
	allocation, err := handler.createEntranceUseCase.Execute(user, beneficiaryID, data)

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

	res := responses.NewBeneficiaryAllocation(allocation)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *BeneficiaryAllocations) Reallocate(w http.ResponseWriter, r *http.Request) {
	var req requests.ReallocateBeneficiary

	user := r.Context().Value("user").(entities.User)
	beneficiaryID := chi.URLParam(r, "id")

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
	allocation, err := handler.createReallocationUseCase.Execute(user, beneficiaryID, data)

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

	res := responses.NewBeneficiaryAllocation(allocation)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *BeneficiaryAllocations) FindManyByBeneficiaryID(w http.ResponseWriter, r *http.Request) {
	beneficiaryID := chi.URLParam(r, "id")
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

	count, allocations, err := handler.findManyByBeneficiaryIDPaginatedUseCase.Execute(user, beneficiaryID, int64(offset), int64(limit))

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

	res := responses.FindMany[responses.BeneficiaryAllocations]{Data: responses.NewBeneficiaryAllocations(allocations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *BeneficiaryAllocations) FindManyByHousingID(w http.ResponseWriter, r *http.Request) {
	housingID := chi.URLParam(r, "id")
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

	count, allocations, err := handler.findManyByHousingIDPaginatedUseCase.Execute(user, housingID, int64(offset), int64(limit))

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

	res := responses.FindMany[responses.BeneficiaryAllocations]{Data: responses.NewBeneficiaryAllocations(allocations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *BeneficiaryAllocations) FindManyByRoomID(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
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

	count, allocations, err := handler.findManyByHousingRoomIDPaginatedUseCase.Execute(user, roomID, int64(offset), int64(limit))

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

	res := responses.FindMany[responses.BeneficiaryAllocations]{Data: responses.NewBeneficiaryAllocations(allocations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
