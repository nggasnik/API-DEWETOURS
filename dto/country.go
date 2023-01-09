package dto

// membuat struct untuk digunakan sebagai tipe data response country
type CountryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// membuat struct untuk digunakan sebagai tipe data request country
type CountryRequest struct {
	Name string `json:"name" validate:"required"`
}
