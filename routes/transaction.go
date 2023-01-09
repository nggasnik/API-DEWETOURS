package routes

import (
	"erlangga/handlers"
	"erlangga/pkg/middleware"
	"erlangga/pkg/mysql"
	"erlangga/repositories"

	"github.com/gorilla/mux"
)

func Transaction(r *mux.Router) {
	transactionRepository := repositories.MakeRepository(mysql.DB)
	h := handlers.HandlerTransaction(transactionRepository)


	r.HandleFunc("/transaction", h.AddTransaction).Methods("POST")

	// menghandle request dengan method PATCH pada endpoint /transaction/{id_transaction}
	// r.HandleFunc("/transaction/{id_transaction}", middleware.UserAuth(middleware.UploadTransactionImage(h.UpdateTransaction))).Methods("PATCH")

	// menghandle request dengan method GET pada endpoint /transaction/{id_transaction}
	r.HandleFunc("/transaction/{id_transaction}", middleware.UserAuth(h.GetDetailTransaction)).Methods("GET")

	// menghandle request dengan method GET pada endpoint /orders
	r.HandleFunc("/orders", middleware.AdminAuth(h.GetAllTransaction)).Methods("GET")

	// menghandle request dengan method GET pada endpoint /orders
	r.HandleFunc("/transaction", middleware.UserAuth(h.GetAllTransactionByUser)).Methods("GET")


	r.HandleFunc("/notification", h.Notification).Methods("POST")

	

	// menghandle request endpoint /transaction dengan method DELETE (digunakan untuk menghapus transaction dengan id tertentu)
	// r.HandleFunc("/transaction/{id_transaction}", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode("Ini endpoint /trip/{id_trip} dengan method delete")
	// }).Methods("DELETE")
}
