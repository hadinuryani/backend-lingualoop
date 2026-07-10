package teacher

import "time"

const (
	TeacherActive   = "ACTIVE"
	TeacherInactive = "INACTIVE"
	RoleTeacher     = "teacher"
)

type Teacher struct {
	ID            string
	UserID        string
	Username      string // Di-join dari tabel users
	Email         string // Di-join dari tabel users
	NIP           string
	FullName      string
	Gender        string
	BirthPlace    *string
	BirthDate     *time.Time
	Phone         *string
	AddressRegion *string
	AddressDetail *string
	Photo         *string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// User struct hanya untuk merepresentasikan tabel users saat membuat guru
type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	FullName     string
	Role         string
	AvatarURL    *string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
