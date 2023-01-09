package handlers

import (
	"encoding/json"
	"erlangga/dto"
	"erlangga/models"
	"erlangga/pkg/bcrypt"
	jwtToken "erlangga/pkg/jwt"
	"erlangga/repositories"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

func (h *handlerAuth) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request dto.RegisterRequest
	json.NewDecoder(r.Body).Decode(&request)
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	hashedPassword, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	newUser := models.User{
		FullName: request.FullName,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: hashedPassword,
		Address:  request.Address,
		Role:     "user",
	}

	userAdded, err := h.AuthRepository.Register(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertRegisterResponse(userAdded),
	}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerAuth) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request dto.LoginRequest
	json.NewDecoder(r.Body).Decode(&request)

	userLogin, err := h.AuthRepository.Login(request.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	isPasswordValid := bcrypt.CheckPassword(request.Password, userLogin.Password)
	if !isPasswordValid {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "wrong email or password"}
		json.NewEncoder(w).Encode(response)
		return
	}

	myClaims := jwt.MapClaims{}
	myClaims["id"] = userLogin.ID
	myClaims["role"] = userLogin.Role
	myClaims["exp"] = time.Now().Add(time.Hour * 8).Unix() // 8 jam expired

	token, errGenerateToken := jwtToken.GenerateToken(&myClaims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		fmt.Println("Unauthorize")
		return
	}

	var loginResponse dto.LoginResponse
	loginResponse.Email = userLogin.Email
	loginResponse.Token = token
	loginResponse.Role = userLogin.Role

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: loginResponse,
	}
	json.NewEncoder(w).Encode(response)
}

func convertRegisterResponse(user models.User) dto.RegisterResponse {
	var result dto.RegisterResponse

	result.FullName = user.FullName
	result.Email = user.Email
	result.Phone = user.Phone
	result.Address = user.Address
	result.Role = user.Role


	return result
}

func (h *handlerAuth) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request dto.RegisterRequest
	json.NewDecoder(r.Body).Decode(&request)

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	hashedPassword, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	newUser := models.User{
		FullName: request.FullName,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: hashedPassword,
		Address:  request.Address,
		Role:     "admin",
	}

	userAdded, err := h.AuthRepository.Register(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertRegisterResponse(userAdded),
	}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerAuth) CheckAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("test")
	claims := r.Context().Value("userInfo").(jwt.MapClaims)
	fmt.Println(claims)
	id := int(claims["id"].(float64))
	
	User, err := h.AuthRepository.GetUserLogin(id)
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
		Data: User,
	}
	json.NewEncoder(w).Encode(response)
}
