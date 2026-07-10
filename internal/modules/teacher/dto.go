package teacher

import "time"

type TeacherRequest struct {
	NIP           string `json:"nip" binding:"required"`
	FullName      string `json:"full_name" binding:"required"`
	Gender        string `json:"gender" binding:"required,oneof=L P"`
	BirthPlace    string `json:"birth_place"`
	BirthDate     string `json:"birth_date"` // YYYY-MM-DD
	Phone         string `json:"phone"`
	AddressRegion string `json:"address_region"`
	AddressDetail string `json:"address_detail"`
}

type TeacherResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	NIP           string    `json:"nip"`
	FullName      string    `json:"full_name"`
	Gender        string    `json:"gender"`
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
