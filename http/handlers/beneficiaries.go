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

	beneficiariesUseCases "relif/platform-bff/usecases/beneficiaries"
	"relif/platform-bff/utils"
	"strconv"
)

type Beneficiaries struct {
	createUseCase                            beneficiariesUseCases.Create
	findManyByOrganizationIDPaginatedUseCase beneficiariesUseCases.FindManyByOrganizationIDPaginated
	findManyByHousingIDPaginatedUseCase      beneficiariesUseCases.FindManyByHousingIDPaginated
	findManyByHousingRoomIDPaginatedUseCase  beneficiariesUseCases.FindManyByHousingRoomIDPaginated
	findOneCompleteByIDUseCase               beneficiariesUseCases.FindOneCompleteByID
	updateOneByIDUseCase                     beneficiariesUseCases.UpdateOneByID
	deleteOneByIDUseCase                     beneficiariesUseCases.DeleteOneByID
	generateProfileImageUploadLinkUseCase    beneficiariesUseCases.GenerateProfileImageUploadLink
}

func NewBeneficiaries(
	createUseCase beneficiariesUseCases.Create,
	findManyByOrganizationIDPaginatedUseCase beneficiariesUseCases.FindManyByOrganizationIDPaginated,
	findManyByHousingIDPaginatedUseCase beneficiariesUseCases.FindManyByHousingIDPaginated,
	findManyByHousingRoomIDPaginatedUseCase beneficiariesUseCases.FindManyByHousingRoomIDPaginated,
	findOneCompleteByIDUseCase beneficiariesUseCases.FindOneCompleteByID,
	updateOneByIDUseCase beneficiariesUseCases.UpdateOneByID,
	deleteOneByIDUseCase beneficiariesUseCases.DeleteOneByID,
	generateProfileImageUploadLinkUseCase beneficiariesUseCases.GenerateProfileImageUploadLink,
) *Beneficiaries {
	return &Beneficiaries{
		createUseCase:                            createUseCase,
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		findManyByHousingIDPaginatedUseCase:      findManyByHousingIDPaginatedUseCase,
		findManyByHousingRoomIDPaginatedUseCase:  findManyByHousingRoomIDPaginatedUseCase,
		findOneCompleteByIDUseCase:               findOneCompleteByIDUseCase,
		updateOneByIDUseCase:                     updateOneByIDUseCase,
		deleteOneByIDUseCase:                     deleteOneByIDUseCase,
		generateProfileImageUploadLinkUseCase:    generateProfileImageUploadLinkUseCase,
	}
}

func (handler *Beneficiaries) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateBeneficiary

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
	beneficiary, err := handler.createUseCase.Execute(user, organizationID, data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrBeneficiaryAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewBeneficiary(beneficiary)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) FindManyByHousingID(w http.ResponseWriter, r *http.Request) {
	housingID := chi.URLParam(r, "id")
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

	count, beneficiaries, err := handler.findManyByHousingIDPaginatedUseCase.Execute(user, housingID, search, int64(offset), int64(limit))

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

	res := responses.FindMany[responses.Beneficiaries]{Data: responses.NewBeneficiaries(beneficiaries), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) FindManyByRoomID(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
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

	count, beneficiaries, err := handler.findManyByHousingRoomIDPaginatedUseCase.Execute(user, roomID, search, int64(offset), int64(limit))

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

	res := responses.FindMany[responses.Beneficiaries]{Data: responses.NewBeneficiaries(beneficiaries), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, beneficiaries, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, search, int64(offset), int64(limit))

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

	res := responses.FindMany[responses.Beneficiaries]{Data: responses.NewBeneficiaries(beneficiaries), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) FindOneByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	beneficiary, err := handler.findOneCompleteByIDUseCase.Execute(user, id)

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

	res := responses.NewBeneficiary(beneficiary)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateBeneficiary

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
		case errors.Is(err, utils.ErrBeneficiaryNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Beneficiaries) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.deleteOneByIDUseCase.Execute(user, id); err != nil {
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

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Beneficiaries) GenerateProfileImageUploadLink(w http.ResponseWriter, r *http.Request) {
	var req requests.GenerateFileUploadLink

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
	link, err := handler.generateProfileImageUploadLinkUseCase.Execute(user, data.Type)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewGenerateFileUploadLink(link)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
