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
	organizationsUseCases "relif/platform-bff/usecases/organizations"
	"relif/platform-bff/utils"
	"strconv"
)

type Organizations struct {
	createUseCase            organizationsUseCases.Create
	findManyUseCase          organizationsUseCases.FindManyPaginated
	findOneByIDUseCase       organizationsUseCases.FindOneByID
	updateOneByIDUseCase     organizationsUseCases.UpdateOneByID
	inactivateOneByIDUseCase organizationsUseCases.InactivateOneByID
	reactivateOneByIDUseCase organizationsUseCases.ReactivateOneByID
}

func NewOrganizations(
	createUseCase organizationsUseCases.Create,
	findManyUseCase organizationsUseCases.FindManyPaginated,
	findOneByIDUseCase organizationsUseCases.FindOneByID,
	updateOneByIDUseCase organizationsUseCases.UpdateOneByID,
	inactivateOneByIDUseCase organizationsUseCases.InactivateOneByID,
	reactivateOneByIDUseCase organizationsUseCases.ReactivateOneByID,
) *Organizations {
	return &Organizations{
		createUseCase:            createUseCase,
		findManyUseCase:          findManyUseCase,
		findOneByIDUseCase:       findOneByIDUseCase,
		updateOneByIDUseCase:     updateOneByIDUseCase,
		inactivateOneByIDUseCase: inactivateOneByIDUseCase,
		reactivateOneByIDUseCase: reactivateOneByIDUseCase,
	}
}

func (handler *Organizations) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateOrganization

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
	organization, err := handler.createUseCase.Execute(user, data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewOrganization(organization)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Organizations) FindMany(w http.ResponseWriter, r *http.Request) {
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

	count, organizations, err := handler.findManyUseCase.Execute(int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.Organizations]{Data: responses.NewOrganizations(organizations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Organizations) FindOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	organization, err := handler.findOneByIDUseCase.Execute(id)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewOrganization(organization)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Organizations) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateOrganization

	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Organizations) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.inactivateOneByIDUseCase.Execute(user, id); err != nil {
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

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Organizations) ReactivateOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.reactivateOneByIDUseCase.Execute(user, id); err != nil {
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

	w.WriteHeader(http.StatusNoContent)
}
