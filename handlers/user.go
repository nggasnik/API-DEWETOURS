package handlers

import (
	"erlangga/dto"
	"erlangga/models"
	"erlangga/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// membuat struct handler sebagai tipe data untuk menampung interface UserRepository dari package repositories
type handlerUser struct {
	UserRepository repositories.UserRepository
}

// membuat function yang mengembalikan object berbentuk struct handleUser, fungsi ini membutuhkan interface UserRepository sebagai parameter, karena method-method dalam interface tersebut dibutuhkan di dalam method struct handleUser
func HandlerUser(UserRepository repositories.UserRepository) *handlerUser {
	return &handlerUser{UserRepository}
}

func (h *handlerUser) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil seluruh data user
	user, err := h.UserRepository.FindUsers()
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
		Data: convertMultipleUserResponse(user, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id_user"])

	user, err := h.UserRepository.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	userDeleted, err := h.UserRepository.DeleteUser(user)
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
		Data: convertDeleteUserResponse(userDeleted),
	}

	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) GetDetailUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// id, _ := strconv.Atoi(mux.Vars(r)["id_user"])

	claims := r.Context().Value("userInfo").(jwt.MapClaims)
	id := int(claims["id"].(float64))

	user, err := h.UserRepository.GetUser(id)
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
		Data: convertOneUserResponse(user, r),
	}

	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) UpdateImageUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil id dari url parameter
	// id, _ := strconv.Atoi(mux.Vars(r)["id_user"])

	claims := r.Context().Value("userInfo").(jwt.MapClaims)
	id := int(claims["id"].(float64))

	// mengambil satu data User berdasarkan id
	User, err := h.UserRepository.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// mengambil data dari context yang dikirim oleh middleware
	dataContex := r.Context().Value("imageUploaded")

	// mengubah property Name pada object User dengan nilai yang didapatkan dari request
	if dataContex != "" {
		User.Image = dataContex.(string)
	}

	// mengirim data User yang sudah diupdate ke database
	UserUpdated, err := h.UserRepository.UpdateUser(User)
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
		Data: convertOneUserResponse(UserUpdated, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

// membuat fungsi konversi data yang akan disajikan sebagai response sesuai requirement
func convertMultipleUserResponse(u []models.User, r *http.Request) []dto.UserResponse {
	var result []dto.UserResponse

	for _, user := range u {
		if user.Image != "" {
			result = append(result, dto.UserResponse{
				ID:       user.ID,
				FullName: user.FullName,
				Email:    user.Email,
				Phone:    user.Phone,
				Address:  user.Address,
				Role:     user.Role,
				Image:    fmt.Sprintf("http://%s/uploads/img/%s", r.Host, user.Image),
			})
		} else {
			result = append(result, dto.UserResponse{
				ID:       user.ID,
				FullName: user.FullName,
				Email:    user.Email,
				Phone:    user.Phone,
				Address:  user.Address,
				Role:     user.Role,
				Image:    user.Image,
			})
		}
	}

	return result
}

func convertDeleteUserResponse(u models.User) dto.UserDeleteResponse {
	return dto.UserDeleteResponse{
		ID: u.ID,
	}
}

func convertOneUserResponse(u models.User, r *http.Request) dto.UserResponse {
	if u.Image != "" {
		return dto.UserResponse{
			ID:       u.ID,
			FullName: u.FullName,
			Email:    u.Email,
			Phone:    u.Phone,
			Address:  u.Address,
			Role:     u.Role,
			Image:    fmt.Sprintf("http://%s/uploads/img/%s", r.Host, u.Image),
		}
	} else {
		return dto.UserResponse{
			ID:       u.ID,
			FullName: u.FullName,
			Email:    u.Email,
			Phone:    u.Phone,
			Address:  u.Address,
			Role:     u.Role,
			Image:    u.Image,
		}
	}
}
