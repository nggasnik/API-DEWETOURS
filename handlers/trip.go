package handlers

import (
	"context"
	"encoding/json"
	"erlangga/dto"
	"erlangga/models"
	"erlangga/repositories"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// membuat struct handler sebagai tipe data untuk menampung interface TripRepository dari package repositories
type handlerTrip struct {
	TripRepository repositories.TripRepository
}

// membuat function yang mengembalikan object berbentuk struct handleTrip, fungsi ini membutuhkan interface TripRepository sebagai parameter, karena method-method dalam interface tersebut dibutuhkan di dalam method struct handleTrip
func HandlerTrip(TripRepository repositories.TripRepository) *handlerTrip {
	return &handlerTrip{TripRepository}
}

func (h *handlerTrip) AddTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Tes")

	if err := r.ParseMultipartForm(1024); err != nil {
		panic(err.Error())
	}

	// mengambil data dari request body
	var requestData dto.AddTripRequest
	requestData.Title = r.FormValue("title")
	requestData.CountryID, _ = strconv.Atoi(r.FormValue("id_country"))
	requestData.Accomodation = r.FormValue("accomodation")
	requestData.Transportation = r.FormValue("transportation")
	requestData.Eat = r.FormValue("eat")
	requestData.Day, _ = strconv.Atoi(r.FormValue("day"))
	requestData.Night, _ = strconv.Atoi(r.FormValue("night"))
	requestData.DateTrip = r.FormValue("dateTrip")
	requestData.Price, _ = strconv.Atoi(r.FormValue("price"))
	requestData.Quota, _ = strconv.Atoi(r.FormValue("quota"))
	requestData.Description = r.FormValue("description")

	fmt.Println(r.FormValue("title"))

	fmt.Println(requestData)

	// mengambil array image dari context yang dikirim oleh middleware, lalu append ke request Image
	dataContex := r.Context().Value("arrImages").([]string)

	arrImages := append(requestData.Images, dataContex...)

	for _, img := range arrImages {
		// create empty context
		var ctx = context.Background()

		// setup cloudinary credentials
		var CLOUD_NAME = os.Getenv("dn6qjddsi")
		var API_KEY = os.Getenv("886395169118187")
		var API_SECRET = os.Getenv("8HT8ttZ1pieoceVQgFIyfM9Wxe4")

		// create new instance of cloudinary object using cloudinary credentials
		cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

		// Upload file to Cloudinary
		resp, err := cld.Upload.Upload(ctx, img, uploader.UploadParams{Folder: "nggasnik"})
		if err != nil {
			fmt.Println(err.Error())
		}
		requestData.Images = append(requestData.Images, resp.SecureURL)
	}

	// memvalidasi inputan dari request body berdasarkan struct dto.TripRequest
	validation := validator.New()
	err := validation.Struct(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// membuat data trip baru yang berisikan data dari request body
	newTrip := models.Trip{
		Title:          requestData.Title,
		CountryID:      1,
		Accomodation:   requestData.Accomodation,
		Transportation: requestData.Transportation,
		Eat:            requestData.Eat,
		Day:            requestData.Day,
		Night:          requestData.Night,
		Price:          requestData.Price,
		Quota:          requestData.Quota,
		Description:    requestData.Description,
	}
	// mengkonversi data request body berbentuk string sebelum dimasukkan ke dateTrip
	newTrip.DateTrip, err = time.Parse("2 January 2006", requestData.DateTrip)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println(newTrip)

	// mengisikan array image file ke array image milik object newTrip, image yang ditambahkan disini nantinya akan otomatis ditambahkan pula ke tabel Image yang berelasi (sesuai dengan association method yang ada di doc gorm)
	for _, img := range requestData.Images {
		imgData := models.ImageResponse{
			FileName: img,
		}
		newTrip.Image = append(newTrip.Image, imgData)
	}

	// mengirim data trip baru ke database
	tripAdded, err := h.TripRepository.CreateTrip(newTrip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// mengambil detail trip yang baru saja ditambahkan (perlu diambil ulang, karena hasil dari tripAdded hanya ada country_id saja, tanpa ada detail country nya)
	getTripAdded, _ := h.TripRepository.GetTrip(tripAdded.ID)

	// menyiapkan response
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOneTripResponse(getTripAdded, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrip) GetAllTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Tes find trip")

	// mengambil seluruh data trip
	trip, err := h.TripRepository.FindTrips()
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
		Data: convertMultipleTripResponse(trip, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrip) GetDetailTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id_trip"])

	// mengambil seluruh data trip
	trip, err := h.TripRepository.GetTrip(id)
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
		Data: convertOneTripResponse(trip, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrip) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id_trip"])

	if err := r.ParseForm(); err != nil {
		panic(err.Error())
	}

	// mengambil data dari request body
	var requestData dto.UpdateTripRequest
	requestData.Title = r.FormValue("title")
	requestData.CountryID, _ = strconv.Atoi(r.FormValue("id_country"))
	requestData.Accomodation = r.FormValue("accomodation")
	requestData.Transportation = r.FormValue("transportation")
	requestData.Eat = r.FormValue("eat")
	requestData.Day, _ = strconv.Atoi(r.FormValue("day"))
	requestData.Night, _ = strconv.Atoi(r.FormValue("night"))
	requestData.DateTrip = r.FormValue("dateTrip")
	requestData.Price, _ = strconv.Atoi(r.FormValue("price"))
	requestData.Quota, _ = strconv.Atoi(r.FormValue("quota"))
	requestData.Description = r.FormValue("description")

	// mengambil array image dari context yang dikirim oleh middleware, lalu append ke request Image
	dataContex := r.Context().Value("arrImages").([]string)
	requestData.Images = append(requestData.Images, dataContex...)

	// mengambil data yang ingin diupdate berdasarkan id yang didapatkan dari url
	tripWantToUpdate, err := h.TripRepository.GetTrip(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusNotFound, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// mengupdate data trip dengan data dari request body
	tripWantToUpdate.ID = id
	if requestData.Title != "" {
		tripWantToUpdate.Title = requestData.Title
	}
	if r.FormValue("id_country") != "" {
		tripWantToUpdate.CountryID = requestData.CountryID
	}
	if r.FormValue("accomodation") != "" {
		tripWantToUpdate.Accomodation = requestData.Accomodation
	}
	if r.FormValue("transportation") != "" {
		tripWantToUpdate.Transportation = requestData.Transportation
	}
	if r.FormValue("eat") != "" {
		tripWantToUpdate.Eat = requestData.Eat
	}
	if r.FormValue("day") != "" {
		tripWantToUpdate.Day = requestData.Day
	}
	if r.FormValue("night") != "" {
		tripWantToUpdate.Night = requestData.Night
	}
	if r.FormValue("dateTrip") != "" {
		// mengkonversi data request body berbentuk string sebelum dimasukkan ke dateTrip
		tripWantToUpdate.DateTrip, err = time.Parse("2 January 2006", requestData.DateTrip)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	if r.FormValue("price") != "" {
		tripWantToUpdate.Price = requestData.Price
	}
	if r.FormValue("quota") != "" {
		tripWantToUpdate.Quota = requestData.Quota
	}
	if r.FormValue("description") != "" {
		tripWantToUpdate.Description = requestData.Description
	}
	// mereplace gambar jika ada gambar yang diuplad
	arrReqImagesLength := len(requestData.Images)
	if arrReqImagesLength > 0 {
		// mengambil panjang array image dari data yang akan diupdate
		arrPrevImagesLength := len(tripWantToUpdate.Image)
		// fmt.Println(arrPrevImagesLength)

		// mengupdate array image file ke array image milik object tripWantToUpdate
		for i, img := range requestData.Images {
			// mengganti gambar sebelumnya (sesuai jumlah gambar)
			if i < arrPrevImagesLength {
				tripWantToUpdate.Image[i].FileName = img
			} else { // jika gambar yang diupload lebih banyak dari gambar sebelumnya, maka tambahkan gambar baru
				imgData := models.ImageResponse{
					FileName: img,
				}
				tripWantToUpdate.Image = append(tripWantToUpdate.Image, imgData)
			}
		}

		// menghapus gambar lama jika gambar yang baru yang direquest lebih sedikit
		for i := range tripWantToUpdate.Image {
			if i >= arrReqImagesLength {
				tripWantToUpdate.Image[i].FileName = "unused"
			}
		}
	}

	// mengirim data trip yang sudah diupdate ke database
	tripUpdated, err := h.TripRepository.UpdateTrip(tripWantToUpdate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// mengambil detail trip yang baru saja ditambahkan (perlu diambil ulang, karena hasil dari tripAdded hanya ada country_id saja, tanpa ada detail country nya)
	getTripUpdated, _ := h.TripRepository.GetTrip(tripUpdated.ID)

	// menyiapkan response
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOneTripResponse(getTripUpdated, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrip) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil id dari url
	id, _ := strconv.Atoi(mux.Vars(r)["id_trip"])

	// mencari data trip berdasarkan id
	trip, err := h.TripRepository.GetTrip(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menghapus data yang dimaksud
	tripDeleted, err := h.TripRepository.DeleteTrip(trip)
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
		Data: tripDeleted.ID,
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func convertMultipleTripResponse(t []models.Trip, r *http.Request) []dto.TripResponse {
	var results []dto.TripResponse

	for _, trip := range t {
		var result dto.TripResponse

		result.ID = trip.ID
		result.Title = trip.Title
		result.Country = trip.Country
		result.Accomodation = trip.Accomodation
		result.Transportation = trip.Transportation
		result.Eat = trip.Eat
		result.Day = trip.Day
		result.Night = trip.Night
		result.DateTrip = trip.DateTrip.Format("02 January 2006")
		result.Price = trip.Price
		result.Quota = trip.Quota
		result.Description = trip.Description

		for _, img := range trip.Image {
			imgUrl := fmt.Sprintf("http://%s/uploads/img/%s", r.Host, img.FileName)
			result.Images = append(result.Images, imgUrl)
		}

		results = append(results, result)
	}

	return results
}

func convertOneTripResponse(t models.Trip, r *http.Request) dto.TripResponse {
	var result dto.TripResponse
	result.ID = t.ID
	result.Title = t.Title
	result.Country = t.Country
	result.Accomodation = t.Accomodation
	result.Transportation = t.Transportation
	result.Eat = t.Eat
	result.Day = t.Day
	result.Night = t.Night
	result.DateTrip = t.DateTrip.Format("02 January 2006")
	result.Price = t.Price
	result.Quota = t.Quota
	result.Description = t.Description

	for _, img := range t.Image {
		imgUrl := fmt.Sprintf("http://%s/uploads/img/%s", r.Host, img.FileName)
		result.Images = append(result.Images, imgUrl)
	}

	// fmt.Println(result.Images)
	return result
}
