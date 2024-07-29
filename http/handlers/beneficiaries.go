package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"relif/bff/http/requests"
	"relif/bff/http/responses"
	"relif/bff/services"
)

type Beneficiaries struct {
	service services.Beneficiaries
}

func NewBeneficiaries(service services.Beneficiaries) *Beneficiaries {
	return &Beneficiaries{
		service: service,
	}
}

func (handler *Beneficiaries) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateBeneficiary

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

	beneficiary, err := handler.service.Create(req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewBeneficiary(beneficiary)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) FindManyByHousingId(w http.ResponseWriter, r *http.Request) {
	housingId := chi.URLParam(r, "id")

	beneficiaries, err := handler.service.FindManyByHousingId(housingId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewBeneficiaries(beneficiaries)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) FindManyByRoomId(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "id")

	beneficiaries, err := handler.service.FindManyByRoomId(roomId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewBeneficiaries(beneficiaries)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) FindOneById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	beneficiary, err := handler.service.FindOneById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	updated, err := handler.service.FindOneAndUpdateById(id, req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewBeneficiary(updated)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Beneficiaries) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := handler.service.DeleteOneById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
