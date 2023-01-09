package routes

import (
	"erlangga/handlers"
	"erlangga/pkg/middleware"
	"erlangga/pkg/mysql"
	"erlangga/repositories"

	"github.com/gorilla/mux"
)

func Auth(r *mux.Router) {
	authRepository := repositories.MakeRepository(mysql.DB)
	h := handlers.HandlerAuth(authRepository)

	// menghandle request dengan method POST pada endpoint /register
	r.HandleFunc("/register", h.Register).Methods("POST")

	// menghandle request dengan method POST pada endpoint /register
	r.HandleFunc("/register-admin", h.RegisterAdmin).Methods("POST")

	// menghandle request dengan method POST pada endpoint /login
	r.HandleFunc("/login", h.Login).Methods("POST")

	r.HandleFunc("/check-auth",middleware.UserAuth(h.CheckAuth)).Methods("GET")
}
