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
	donationsUseCases "relif/platform-bff/usecases/donations"
	"relif/platform-bff/utils"
	"strconv"
)

type Donations struct {
	createUseCase                           donationsUseCases.Create
	findManyByBeneficiaryIDPaginatedUseCase donationsUseCases.FindManyByBeneficiaryIDPaginated
}

func NewDonations(
	createUseCase donationsUseCases.Create,
	findManyByBeneficiaryIDPaginatedUseCase donationsUseCases.FindManyByBeneficiaryIDPaginated,
) *Donations {
	return &Donations{
		createUseCase:                           createUseCase,
		findManyByBeneficiaryIDPaginatedUseCase: findManyByBeneficiaryIDPaginatedUseCase,
	}
}

func (handler *Donations) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateDonation

	beneficiaryID := chi.URLParam(r, "id")
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
	donation, err := handler.createUseCase.Execute(user, beneficiaryID, data)

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

	res := responses.NewDonation(donation)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Donations) FindManyByBeneficiaryID(w http.ResponseWriter, r *http.Request) {
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

	count, donations, err := handler.findManyByBeneficiaryIDPaginatedUseCase.Execute(user, beneficiaryID, int64(offset), int64(limit))

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

	res := responses.FindMany[responses.Donations]{Data: responses.NewDonations(donations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
