package routes

import (
	"erlangga/handlers"
	"erlangga/pkg/middleware"
	"erlangga/pkg/mysql"
	"erlangga/repositories"

	"github.com/gorilla/mux"
)

func User(r *mux.Router) {
	UserRepository := repositories.MakeRepository(mysql.DB)
	h := handlers.HandlerUser(UserRepository)

	// menghandle request dengan method GET pada endpoint /users
	r.HandleFunc("/users", middleware.UserAuth(h.GetAllUsers)).Methods("GET")

	// menghandle request dengan method DELETE pada endpoint /user
	r.HandleFunc("/user/{id_user}", middleware.AdminAuth(h.DeleteUser)).Methods("DELETE")

	// menghandle request dengan method GET pada endpoint /user <improve>
	r.HandleFunc("/user/{id}", middleware.UserAuth(h.GetDetailUser)).Methods("GET")

	// menghandle request dengan method PATCH pada endpoint /user <improve>
	r.HandleFunc("/user", middleware.UserAuth(middleware.UpdateUserImage(h.UpdateImageUser))).Methods("PATCH")
}
