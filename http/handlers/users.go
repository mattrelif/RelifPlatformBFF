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
	usersUseCases "relif/platform-bff/usecases/users"
	"relif/platform-bff/utils"
	"strconv"
)

type Users struct {
	findOneCompleteByIDUseCase               usersUseCases.FindOneCompleteByID
	findManyByOrganizationIDPaginatedUseCase usersUseCases.FindManyByOrganizationIDPaginated
	findManyRelifMembersPaginatedUseCase     usersUseCases.FindManyRelifMembersPaginated
	updateOneByIDUseCase                     usersUseCases.UpdateOneByID
	inactivateOneByIDUseCase                 usersUseCases.InactivateOneByID
	reactivateOneByIDUseCase                 usersUseCases.ReactivateOneByID
}

func NewUsers(
	findOneCompleteByIDUseCase usersUseCases.FindOneCompleteByID,
	findManyByOrganizationIDPaginatedUseCase usersUseCases.FindManyByOrganizationIDPaginated,
	findManyRelifMembersPaginatedUseCase usersUseCases.FindManyRelifMembersPaginated,
	updateOneByIDUseCase usersUseCases.UpdateOneByID,
	inactivateOneByIDUseCase usersUseCases.InactivateOneByID,
	reactivateOneByIDUseCase usersUseCases.ReactivateOneByID,
) *Users {
	return &Users{
		findOneCompleteByIDUseCase:               findOneCompleteByIDUseCase,
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		findManyRelifMembersPaginatedUseCase:     findManyRelifMembersPaginatedUseCase,
		updateOneByIDUseCase:                     updateOneByIDUseCase,
		inactivateOneByIDUseCase:                 inactivateOneByIDUseCase,
		reactivateOneByIDUseCase:                 reactivateOneByIDUseCase,
	}
}

func (handler *Users) FindOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := handler.findOneCompleteByIDUseCase.Execute(id)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewUser(user)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Users) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
	organizationID := chi.URLParam(r, "id")
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

	count, users, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.Users]{Data: responses.NewUsers(users), Count: count}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Users) FindManyRelifMembers(w http.ResponseWriter, r *http.Request) {
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

	count, users, err := handler.findManyRelifMembersPaginatedUseCase.Execute(user, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.Users]{Data: responses.NewUsers(users), Count: count}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Users) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateUser

	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := req.ToEntity()
	if err = handler.updateOneByIDUseCase.Execute(user, id, data); err != nil {
		switch {
		case errors.Is(err, utils.ErrUserNotFound):
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

func (handler *Users) InactivateOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.inactivateOneByIDUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrUserNotFound):
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

func (handler *Users) ReactivateOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.reactivateOneByIDUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrUserNotFound):
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
