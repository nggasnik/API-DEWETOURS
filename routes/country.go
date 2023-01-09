package routes

import (
	"erlangga/handlers"
	"erlangga/pkg/middleware"
	"erlangga/pkg/mysql"
	"erlangga/repositories"

	"github.com/gorilla/mux"
)

func Country(r *mux.Router) {
	// membuat object dari struct repository yang berisikan koneksi database dan beberapa method untuk komunikasi dengan databaase
	countryRepository := repositories.MakeRepository(mysql.DB)

	// membuat object dari struct handlerCountry yang berisikan method-method milik struct handleCountry, dan interface CountryRepository yang didalamnya terdapat beberapa method milik struct repository.
	h := handlers.HandlerCountry(countryRepository)

	// menghandle request dengan method GET pada endpoint /country
	r.HandleFunc("/country", h.GetAllCountry).Methods("GET")

	// menghandle request dengan method GET pada endpoint /country/{id_country}
	r.HandleFunc("/country/{id_country}", h.GetDetailCountry).Methods("GET")

	// menghandle request dengan method POST pada endpoint /country
	r.HandleFunc("/country", middleware.UserAuth(h.AddCountry)).Methods("POST")

	// menghandle request dengan method PATCH pada endpoint /country/{id_country}
	r.HandleFunc("/country/{id_country}", middleware.UserAuth(h.UpdateCountry)).Methods("PATCH")

	// menghandle request dengan method DELETE pada endpoint /country
	r.HandleFunc("/country/{id_country}", middleware.UserAuth(h.DeleteCountry)).Methods("DELETE")

}
