package student

import "time"

type StudentRequest struct {
	NIS           string `json:"nis" binding:"required"`
	FullName      string `json:"full_name" binding:"required"`
	Gender        string `json:"gender" binding:"required,oneof=L P"`
	MajorID       string `json:"major_id" binding:"required"`
	ClassLevel    string `json:"class_level" binding:"required"`
	BirthPlace    string `json:"birth_place"`
	BirthDate     string `json:"birth_date"` // YYYY-MM-DD
	Phone         string `json:"phone"`
	AddressRegion string `json:"address_region"`
	AddressDetail string `json:"address_detail"`
	Photo         string `json:"photo"`
}

type StudentResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	NIS           string    `json:"nis"`
	FullName      string    `json:"full_name"`
	Gender        string    `json:"gender"`
	MajorID       string    `json:"major_id,omitempty"`
	ClassLevel    string    `json:"class_level,omitempty"`
	BirthPlace    string    `json:"birth_place,omitempty"`
	BirthDate     string    `json:"birth_date,omitempty"` // YYYY-MM-DD
	Phone         string    `json:"phone,omitempty"`
	AddressRegion string    `json:"address_region,omitempty"`
	AddressDetail string    `json:"address_detail,omitempty"`
	Photo         string    `json:"photo,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
