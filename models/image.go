package models

type Image struct {
	ID       int    `json:"id"`
	FileName string `json:"file_name" gorm:"type:varchar(255)"`
	TripID   int
	Trip     TripResponse
}

type ImageResponse struct {
	ID       int    `json:"-"`
	FileName string `json:"file_name"`
	TripID   int    `json:"-"`
}

func (ImageResponse) TableName() string {
	return "images"
}
