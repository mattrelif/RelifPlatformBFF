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
	voluntaryPeopleUseCases "relif/platform-bff/usecases/voluntary_people"
	"relif/platform-bff/utils"
	"strconv"
)

type VoluntaryPeople struct {
	createUseCase                            voluntaryPeopleUseCases.Create
	findManyByOrganizationIDPaginatedUseCase voluntaryPeopleUseCases.FindManyByOrganizationIDPaginated
	findOneByIDUseCase                       voluntaryPeopleUseCases.FindOneByID
	updateOneByIDUseCase                     voluntaryPeopleUseCases.UpdateOneByID
	deleteOneByIDUseCase                     voluntaryPeopleUseCases.DeleteOneByID
}

func NewVoluntaryPeople(
	createUseCase voluntaryPeopleUseCases.Create,
	findManyByOrganizationIDPaginatedUseCase voluntaryPeopleUseCases.FindManyByOrganizationIDPaginated,
	findOneByIDUseCase voluntaryPeopleUseCases.FindOneByID,
	updateOneByIDUseCase voluntaryPeopleUseCases.UpdateOneByID,
	deleteOneByIDUseCase voluntaryPeopleUseCases.DeleteOneByID,
) *VoluntaryPeople {
	return &VoluntaryPeople{
		createUseCase:                            createUseCase,
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		findOneByIDUseCase:                       findOneByIDUseCase,
		updateOneByIDUseCase:                     updateOneByIDUseCase,
		deleteOneByIDUseCase:                     deleteOneByIDUseCase,
	}
}

func (handler *VoluntaryPeople) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateVoluntaryPerson

	organizationID := chi.URLParam(r, "id")
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
	voluntary, err := handler.createUseCase.Execute(user, organizationID, data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrVoluntaryPersonAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewVoluntaryPerson(voluntary)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *VoluntaryPeople) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, voluntaries, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, search, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.VoluntaryPeople]{Data: responses.NewVoluntaryPeople(voluntaries), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *VoluntaryPeople) FindOneByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	voluntary, err := handler.findOneByIDUseCase.Execute(user, id)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrVoluntaryPersonNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	res := responses.NewVoluntaryPerson(voluntary)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *VoluntaryPeople) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateVoluntaryPerson

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
		case errors.Is(err, utils.ErrVoluntaryPersonNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *VoluntaryPeople) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.deleteOneByIDUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrVoluntaryPersonNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
