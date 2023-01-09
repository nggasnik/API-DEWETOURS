package repositories

import (
	"erlangga/models"
)

type CountryRepository interface {
	FindCountry() ([]models.Country, error)
	GetCountry(ID int) (models.Country, error)
	CreateCountry(newCountry models.Country) (models.Country, error)
	DeleteCountry(country models.Country) (models.Country, error)
	UpdateCountry(country models.Country) (models.Country, error)
}

func (r *repository) FindCountry() ([]models.Country, error) {
	var country []models.Country
	err := r.db.Find(&country).Error
	return country, err
}

func (r *repository) GetCountry(ID int) (models.Country, error) {
	var country models.Country

	err := r.db.First(&country, ID).Error
	return country, err
}

func (r *repository) CreateCountry(newCountry models.Country) (models.Country, error) {
	err := r.db.Create(&newCountry).Error
	return newCountry, err
}

func (r *repository) DeleteCountry(country models.Country) (models.Country, error) {
	err := r.db.Delete(&country).Error
	return country, err
}

func (r *repository) UpdateCountry(country models.Country) (models.Country, error) {
	err := r.db.Save(&country).Error
	return country, err
}
