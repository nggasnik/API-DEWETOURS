package handlers

import (
	"encoding/json"
	"erlangga/dto"
	"erlangga/models"
	"erlangga/repositories"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	// "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"

	// "github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"

	"gopkg.in/gomail.v2"
)

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) AddTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil data dari request form
	var request dto.TransactionRequest
	request.CounterQty, _ = strconv.Atoi(r.FormValue("counterQty"))
	request.Total, _ = strconv.Atoi(r.FormValue("total"))
	request.Status = r.FormValue("status")
	request.TripID, _ = strconv.Atoi(r.FormValue("trip_id"))
	// request.Attachment = r.Context().Value("attachment").(string)
	request.UserId, _ = strconv.Atoi(r.FormValue("userId"))

	fmt.Println(request.UserId)
	newTransaction := models.Transaction{
		CounterQty: request.CounterQty,
		Total:      request.Total,
		Status:     request.Status,
		TripID:     request.TripID,
		// Attachment: request.Attachment,
		UserID:  1,
	}

	newTransaction.BookingDate = time.Now()
	transactionAdded, err := h.TransactionRepository.CreateTransaction(newTransaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}


	transactionSaya, err := h.TransactionRepository.GetTransaction(transactionAdded.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var TransIdIsMatch = false
	var TransactionId int
	for !TransIdIsMatch {
		TransactionId = rand.Intn(10000) - rand.Intn(100)
		transactionData, _ := h.TransactionRepository.GetTransaction(TransactionId)
		if transactionData.ID == 0 {
			TransIdIsMatch = true
		}
	}
	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(("SB-Mid-server-nycCXqQZFycSu2OCeIUxGI8W"), midtrans.Sandbox)
	// Use to midtrans.Production if you want Production Environment (accept real transaction).

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transactionSaya.ID),
			GrossAmt: int64(transactionSaya.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transactionSaya.User.FullName,
			Email: transactionSaya.User.Email,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, _ := s.CreateTransaction(req)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: snapResp}
	json.NewEncoder(w).Encode(response)
}

// func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	id, _ := strconv.Atoi(mux.Vars(r)["id_transaction"])

// 	if err := r.ParseForm(); err != nil {
// 		panic(err.Error())
// 	}

// 	// mengambil data dari request form
// 	var request dto.TransactionRequest
// 	request.CounterQty, _ = strconv.Atoi(r.FormValue("counterQty"))
// 	request.Total, _ = strconv.Atoi(r.FormValue("total"))
// 	request.Status = r.FormValue("status")
// 	request.TripID, _ = strconv.Atoi(r.FormValue("tripId"))
// 	request.Attachment = r.Context().Value("attachment").(string)

// 	// mengambil data yang ingin diupdate berdasarkan id yang didapatkan dari url
// 	transactionWantToUpdate, _ := h.TransactionRepository.GetTransaction(id)

// 	// mengupdate data transaction dengan data dari request body
// 	transactionWantToUpdate.ID = id
// 	if r.FormValue("counterQty") != "" {
// 		transactionWantToUpdate.CounterQty = request.CounterQty
// 	}
// 	if r.FormValue("total") != "" {
// 		transactionWantToUpdate.Total = request.Total
// 	}
// 	if r.FormValue("status") != "" {
// 		if r.FormValue("status") != "new" && r.FormValue("status") != "pending" && r.FormValue("status") != "approve" && r.FormValue("status") != "reject" && r.FormValue("status") != "cancel" {
// 			w.WriteHeader(http.StatusBadRequest)
// 			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Transaction status does not match expected status"}
// 			json.NewEncoder(w).Encode(response)
// 			return
// 		} else if r.FormValue("status") == "pending" && request.Attachment == "" {
// 			w.WriteHeader(http.StatusBadRequest)
// 			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Can't find attachment when update status from 'new' to 'pending'"}
// 			json.NewEncoder(w).Encode(response)
// 			return
// 		}
// 		transactionWantToUpdate.Status = request.Status
// 	}
// 	if r.FormValue("tripId") != "" {
// 		transactionWantToUpdate.TripID = request.TripID
// 	}
// 	if request.Attachment != "" {
// 		transactionWantToUpdate.Attachment = request.Attachment
// 	}

// 	// mengirim data transaction yang sudah diupdate ke database
// 	transactionUpdated, err := h.TransactionRepository.UpdateTransaction(transactionWantToUpdate)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	// mengambil detail transaction yang baru saja ditambahkan (perlu diambil ulang, karena hasil dari transactionAdded hanya ada country_id saja, tanpa ada detail country nya)
// 	getTransactionUpdated, _ := h.TransactionRepository.GetTransaction(transactionUpdated.ID)

// 	// menyiapkan response
// 	response := dto.SuccessResult{
// 		Code: http.StatusOK,
// 		Data: convertOneTransactionResponse(getTransactionUpdated, r),
// 	}

// 	// mengirim response
// 	json.NewEncoder(w).Encode(response)
// }

func (h *handlerTransaction) GetDetailTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil id dari url parameter
	id, _ := strconv.Atoi(mux.Vars(r)["id_transaction"])

	// mengambil satu data Transaction berdasarkan id
	transaction, err := h.TransactionRepository.GetTransaction(id)
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
		Data: convertOneTransactionResponse(transaction, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil seluruh data transaction
	transaction, err := h.TransactionRepository.FindTransactions()
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
		Data: convertMultipleTransactionResponse(transaction, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetAllTransactionByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value("userInfo").(jwt.MapClaims)
	id := int(claims["id"].(float64))

	// mengambil seluruh data transaction
	transaction, err := h.TransactionRepository.FindTransactionsByUser(id)
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
		Data: convertMultipleTransactionResponse(transaction, r),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

// membuat fungsi konversi data yang akan disajikan sebagai response sesuai requirement
func convertOneTransactionResponse(t models.Transaction, r *http.Request) dto.TransactionResponse {
	result := dto.TransactionResponse{
		ID:         t.ID,
		CounterQty: t.CounterQty,
		Total:      t.Total,
		Status:     t.Status,
		Attachment: fmt.Sprintf("http://%s/static/img/%s", r.Host, t.Attachment),
		User:       t.User,
		Trip: dto.TripResponse{
			ID:             t.Trip.ID,
			Title:          t.Trip.Title,
			Country:        t.Trip.Country,
			Accomodation:   t.Trip.Accomodation,
			Transportation: t.Trip.Transportation,
			Eat:            t.Trip.Eat,
			Day:            t.Trip.Day,
			Night:          t.Trip.Night,
			Price:          t.Trip.Price,
			Quota:          t.Trip.Quota,
			Description:    t.Trip.Description,
		},
	}
	if t.Attachment != "" {
		result.Attachment = fmt.Sprintf("http://%s/static/img/%s", r.Host, t.Attachment)
	} else {
		result.Attachment = ""
	}
	result.BookingDate = t.BookingDate.Format("Monday, 2 January 2006")
	result.Trip.DateTrip = t.Trip.DateTrip.Format("2 January 2006")
	for _, img := range t.Trip.Image {
		result.Trip.Images = append(result.Trip.Images, fmt.Sprintf("http://%s/static/img/%s", r.Host, img.FileName))
	}

	return result
}

// membuat fungsi konversi data yang akan disajikan sebagai response sesuai requirement
func convertMultipleTransactionResponse(t []models.Transaction, r *http.Request) []dto.TransactionResponse {
	var result []dto.TransactionResponse

	for _, trx := range t {
		transaction := dto.TransactionResponse{
			ID:         trx.ID,
			CounterQty: trx.CounterQty,
			Total:      trx.Total,
			Status:     trx.Status,
			User:       trx.User,
			Trip: dto.TripResponse{
				ID:             trx.Trip.ID,
				Title:          trx.Trip.Title,
				Country:        trx.Trip.Country,
				Accomodation:   trx.Trip.Accomodation,
				Transportation: trx.Trip.Transportation,
				Eat:            trx.Trip.Eat,
				Day:            trx.Trip.Day,
				Night:          trx.Trip.Night,
				Price:          trx.Trip.Price,
				Quota:          trx.Trip.Quota,
				Description:    trx.Trip.Description,
			},
		}
		if trx.Attachment != "" {
			transaction.Attachment = fmt.Sprintf("http://%s/uploads/img/%s", r.Host, trx.Attachment)
		} else {
			transaction.Attachment = ""
		}
		transaction.BookingDate = trx.BookingDate.Format("Monday, 2 January 2006")
		transaction.Trip.DateTrip = trx.Trip.DateTrip.Format("2 January 2006")
		for _, img := range trx.Trip.Image {
			transaction.Trip.Images = append(transaction.Trip.Images, fmt.Sprintf("http://%s/uploads/img/%s", r.Host, img.FileName))
		}
		result = append(result, transaction)
	}
	return result
}

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)
	id, _ := strconv.Atoi(orderId)

	transs, _ := h.TransactionRepository.GetTransaction(id)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			h.TransactionRepository.UpdateTransactionn("pending", transs.ID)
		} else if fraudStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			SendMail("success", transs)
			h.TransactionRepository.UpdateTransactionn("success", transs.ID)
		}
	} else if transactionStatus == "settlement" {
		// TODO set transaction status on your databaase to 'success'
		SendMail("success", transs)
		h.TransactionRepository.UpdateTransactionn("success", transs.ID)
	} else if transactionStatus == "deny" {
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become successSendMail("failed", transaction) 
		SendMail("failed", transs)
		h.TransactionRepository.UpdateTransactionn("failed", transs.ID)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		// TODO set transaction status on your databaase to 'failure'
		SendMail("failed", transs) 
		h.TransactionRepository.UpdateTransactionn("failed", transs.ID)
	} else if transactionStatus == "pending" {
		// TODO set transaction status on your databaase to 'pending' / waiting payment
		h.TransactionRepository.UpdateTransactionn("pending", transs.ID)
	}

	w.WriteHeader(http.StatusOK)
}

func SendMail(status string, transaction models.Transaction) {

	if status != transaction.Status && (status == "success") {
	  var CONFIG_SMTP_HOST = "smtp.gmail.com"
	  var CONFIG_SMTP_PORT = 587
	  var CONFIG_SENDER_NAME = "DumbMerch <demo.dumbways@gmail.com>"
	  var CONFIG_AUTH_EMAIL = "steelangga@gmail.com"
	  var CONFIG_AUTH_PASSWORD = "thzmoylhtmeqetys"
  
	  var tripName = transaction.Trip.Title
	  var price = strconv.Itoa(transaction.Trip.Price)
  
	  mailer := gomail.NewMessage()
	  mailer.SetHeader("From", CONFIG_SENDER_NAME)
	  mailer.SetHeader("To", transaction.User.Email)
	  mailer.SetHeader("Subject", "Transaction Status")
	  mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	  <html lang="en">
		<head>
		<meta charset="UTF-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Document</title>
		<style>
		  h1 {
		  color: brown;
		  }
		</style>
		</head>
		<body>
		<h2>Product payment :</h2>
		<ul style="list-style-type:none;">
		  <li>Name : %s</li>
		  <li>Total payment: Rp.%s</li>
		  <li>Status : <b>%s</b></li>
		</ul>
		</body>
	  </html>`, tripName, price, status))
  
	  dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	  )
  
	  err := dialer.DialAndSend(mailer)
	  if err != nil {
		log.Fatal(err.Error())
	  }
  
	  log.Println("Mail sent! to " + transaction.User.Email)
	}
  }