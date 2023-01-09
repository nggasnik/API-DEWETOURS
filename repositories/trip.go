package repositories

import (
	"erlangga/models"

	"gorm.io/gorm"
)

type TripRepository interface {
	FindTrips() ([]models.Trip, error)
	GetTrip(ID int) (models.Trip, error)
	CreateTrip(newTrip models.Trip) (models.Trip, error)
	UpdateTrip(trip models.Trip) (models.Trip, error)
	DeleteTrip(trip models.Trip) (models.Trip, error)
}

func (r *repository) FindTrips() ([]models.Trip, error) {
	var trips []models.Trip
	err := r.db.Preload("Country").Preload("Image").Find(&trips).Error
	return trips, err
}

func (r *repository) GetTrip(ID int) (models.Trip, error) {
	var trip models.Trip
	err := r.db.Preload("Country").Preload("Image").First(&trip, ID).Error
	return trip, err
}

func (r *repository) CreateTrip(newTrip models.Trip) (models.Trip, error) {
	err := r.db.Create(&newTrip).Preload("Country").Error
	return newTrip, err
}

func (r *repository) UpdateTrip(trip models.Trip) (models.Trip, error) {
	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&trip).Updates(trip).Error
	r.db.Exec("DELETE from images where file_name = ?", "unused")
	return trip, err
}

func (r *repository) DeleteTrip(trip models.Trip) (models.Trip, error) {
	err := r.db.Select("Image").Delete(&trip).Error
	return trip, err
}
