package routes

import (
	"github.com/gorilla/mux"
)

// membuat middleware Router untuk menampung/mengumpulkan semua router
func RouterInit(r *mux.Router) {

	Auth(r)
	User(r)
	Country(r)
	Trip(r)
	Transaction(r)
}
