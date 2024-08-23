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
	housingRoomsUseCases "relif/platform-bff/usecases/housing_rooms"
	"relif/platform-bff/utils"
	"strconv"
)

type HousingRooms struct {
	createManyUseCase                   housingRoomsUseCases.CreateMany
	findOneCompleteByIDUseCase          housingRoomsUseCases.FindOneCompleteByID
	findManyByHousingIDPaginatedUseCase housingRoomsUseCases.FindManyByHousingIDPaginated
	updateOneByIDUseCase                housingRoomsUseCases.UpdateOneByID
	deleteOneByIDUseCase                housingRoomsUseCases.DeleteOneByID
}

func NewHousingRooms(
	createManyUseCase housingRoomsUseCases.CreateMany,
	findOneCompleteByIDUseCase housingRoomsUseCases.FindOneCompleteByID,
	findManyByHousingIDPaginatedUseCase housingRoomsUseCases.FindManyByHousingIDPaginated,
	updateOneByIDUseCase housingRoomsUseCases.UpdateOneByID,
	deleteOneByIDUseCase housingRoomsUseCases.DeleteOneByID,
) *HousingRooms {
	return &HousingRooms{
		createManyUseCase:                   createManyUseCase,
		findOneCompleteByIDUseCase:          findOneCompleteByIDUseCase,
		findManyByHousingIDPaginatedUseCase: findManyByHousingIDPaginatedUseCase,
		updateOneByIDUseCase:                updateOneByIDUseCase,
		deleteOneByIDUseCase:                deleteOneByIDUseCase,
	}
}

func (handler *HousingRooms) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateManyHousingRooms

	housingID := chi.URLParam(r, "id")
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
	count, rooms, err := handler.createManyUseCase.Execute(user, housingID, data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, utils.ErrOrganizationNotFound), errors.Is(err, utils.ErrHousingNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.HousingRooms]{Data: responses.NewHousingRooms(rooms), Count: count}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *HousingRooms) FindOneByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	room, err := handler.findOneCompleteByIDUseCase.Execute(user, id)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, utils.ErrHousingRoomNotFound), errors.Is(err, utils.ErrHousingNotFound), errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewHousingRoom(room)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *HousingRooms) FindManyByHousingID(w http.ResponseWriter, r *http.Request) {
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

	count, rooms, err := handler.findManyByHousingIDPaginatedUseCase.Execute(user, housingID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, utils.ErrHousingNotFound), errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.HousingRooms]{Data: responses.NewHousingRooms(rooms), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *HousingRooms) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateHousingRoom

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
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, utils.ErrHousingRoomNotFound), errors.Is(err, utils.ErrHousingNotFound), errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *HousingRooms) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.deleteOneByIDUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, utils.ErrHousingRoomNotFound), errors.Is(err, utils.ErrHousingNotFound), errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
