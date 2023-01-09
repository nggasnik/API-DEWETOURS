package repositories

import "erlangga/models"

type TransactionRepository interface {
	CreateTransaction(newTransaction models.Transaction) (models.Transaction, error)
	FindTransactions() ([]models.Transaction, error)
	GetTransaction(ID int) (models.Transaction, error)
	// UpdateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransactionn(status string, ID int) (models.Transaction, error)
	FindTransactionsByUser(UserId int) ([]models.Transaction, error)
	AddTransaction(newTransaction models.Transaction) (models.Transaction, error)
	GetOneTransaction(ID string) (models.Transaction, error)
}

func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("Trip.Country").Preload("Trip.Image").Preload("Trip").Preload("User").Find(&transaction).Error

	return transaction, err
}

func (r *repository) FindTransactionsByUser(UserId int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("Trip.Country").Preload("Trip.Image").Preload("Trip").Preload("User").Where("user_id = ?", UserId).Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetTransaction(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Trip.Country").Preload("Trip.Image").Preload("Trip").Preload("User").First(&transaction, ID).Error

	return transaction, err
}

func (r *repository) CreateTransaction(newTransaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&newTransaction).Error

	return newTransaction, err
}

// func (r *repository) UpdateTransaction(transaction models.Transaction) (models.Transaction, error) {
// 	err := r.db.Model(&transaction).Updates(transaction).Error

// 	return transaction, err
// }

func (r *repository) DeleteTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Delete(&transaction).Error

	return transaction, err
}

func (r *repository) AddTransaction(newTransaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&newTransaction).Error
	return newTransaction, err
}

func (r *repository) UpdateTransactionn(status string, ID int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Product").First(&transaction, ID)

	// If is different & Status is "success" decrement product quantity
	if status != transaction.Status && status == "success" {
		var trip models.Trip
		r.db.First(&trip, transaction.Trip.ID)
		trip.Quota = trip.Quota - 1
		r.db.Save(&trip)
	}

	transaction.Status = status

	err := r.db.Save(&transaction).Error

	return transaction, err
}

func (r *repository) GetOneTransaction(ID string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Trip").Preload("User").First(&transaction, "id = ?", ID).Error

	return transaction, err
}
