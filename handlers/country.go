package handlers

import (
	"erlangga/dto"
	"erlangga/models"
	"erlangga/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type handlerCountry struct {
	CountryRepository repositories.CountryRepository
}

func HandlerCountry(CountryRepository repositories.CountryRepository) *handlerCountry {
	return &handlerCountry{CountryRepository}
}

// membuat fungsi yang digunakan untuk menampilkan list country
func (h *handlerCountry) GetAllCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil seluruh data country
	country, err := h.CountryRepository.FindCountry()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menyiapkan response
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertMultipleCountryResponse(country),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

// membuat fungsi yang digunakan untuk menampilkan detail country
func (h *handlerCountry) GetDetailCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil id dari url parameter
	id, _ := strconv.Atoi(mux.Vars(r)["id_country"])

	// mengambil satu data country berdasarkan id
	country, err := h.CountryRepository.GetCountry(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menyiapkan response
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOneCountryResponse(country),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

// membuat fungsi yang digunakan untuk menambahkan data country baru
func (h *handlerCountry) AddCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil data dari request body
	var request dto.CountryRequest
	json.NewDecoder(r.Body).Decode(&request)

	// memvalidasi inputan dari request body berdasarkan struct dto.CountryRequest
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// membuat object country baru dengan cetakan models.Country
	newCountry := models.Country{
		Name: request.Name,
	}

	// mengirim data country baru ke database
	countryAdded, err := h.CountryRepository.CreateCountry(newCountry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menyiapkan response
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOneCountryResponse(countryAdded),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

// membuat fungsi yang digunakan untuk mengupdate/mengedit data country
func (h *handlerCountry) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil id dari url parameter
	id, _ := strconv.Atoi(mux.Vars(r)["id_country"])

	// mengambil data dari request body
	var request dto.CountryRequest
	json.NewDecoder(r.Body).Decode(&request)

	// mengambil satu data country berdasarkan id
	country, err := h.CountryRepository.GetCountry(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// mengubah property Name pada object country dengan nilai yang didapatkan dari request
	if request.Name != "" {
		country.Name = request.Name
	}

	// mengirim data country yang sudah diupdate ke database
	countryUpdated, err := h.CountryRepository.UpdateCountry(country)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menyiapkan response
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOneCountryResponse(countryUpdated),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

// membuat fungsi yang digunakan untuk menghapus sebuah data country
func (h *handlerCountry) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id_country"])

	country, err := h.CountryRepository.GetCountry(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	countryDeleted, err := h.CountryRepository.DeleteCountry(country)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOneCountryResponse(countryDeleted),
	}

	json.NewEncoder(w).Encode(response)
}

// membuat fungsi konversi data yang akan disajikan sebagai response sesuai requirement
func convertOneCountryResponse(c models.Country) dto.CountryResponse {
	return dto.CountryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}

// membuat fungsi konversi data yang akan disajikan sebagai response sesuai requirement
func convertMultipleCountryResponse(c []models.Country) []dto.CountryResponse {
	var result []dto.CountryResponse

	for _, country := range c {
		result = append(result, dto.CountryResponse{
			ID:   country.ID,
			Name: country.Name,
		})
	}

	return result
}
