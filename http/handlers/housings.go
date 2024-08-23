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
	housingsUseCases "relif/platform-bff/usecases/housings"
	"relif/platform-bff/utils"
	"strconv"
)

type Housings struct {
	createUseCase                            housingsUseCases.Create
	findManyByOrganizationIDPaginatedUseCase housingsUseCases.FindManyByOrganizationIDPaginated
	findOneCompleteByIDUseCase               housingsUseCases.FindOneCompleteByID
	updateOneByIDUseCase                     housingsUseCases.UpdateOneByID
	deleteOneByIDUseCase                     housingsUseCases.DeleteOneByID
}

func NewHousings(
	createUseCase housingsUseCases.Create,
	findManyByOrganizationIDPaginatedUseCase housingsUseCases.FindManyByOrganizationIDPaginated,
	findOneCompleteByIDUseCase housingsUseCases.FindOneCompleteByID,
	updateOneByIDUseCase housingsUseCases.UpdateOneByID,
	deleteOneByIDUseCase housingsUseCases.DeleteOneByID,
) *Housings {
	return &Housings{
		createUseCase:                            createUseCase,
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		findOneCompleteByIDUseCase:               findOneCompleteByIDUseCase,
		updateOneByIDUseCase:                     updateOneByIDUseCase,
		deleteOneByIDUseCase:                     deleteOneByIDUseCase,
	}
}

func (handler *Housings) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateHousing

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
	housing, err := handler.createUseCase.Execute(user, data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewHousing(housing)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Housings) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
	organizationID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	search := r.URL.Query().Get("search")

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

	count, housings, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, search, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.Housings]{Data: responses.NewHousings(housings), Count: count}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Housings) FindOneByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	housing, err := handler.findOneCompleteByIDUseCase.Execute(user, id)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrHousingNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewHousing(housing)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Housings) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateHousing

	id := chi.URLParam(r, "id")
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

	if err = handler.updateOneByIDUseCase.Execute(user, id, data); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrHousingNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Housings) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.deleteOneByIDUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrHousingNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
