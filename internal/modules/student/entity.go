package student

import "time"

const (
	StudentActive    = "ACTIVE"
	StudentGraduated = "GRADUATED"
	StudentTransfer  = "TRANSFER"
	StudentInactive  = "INACTIVE"
	RoleStudent      = "student"
)

type Student struct {
	ID            string
	UserID        string
	Username      string // Join dari users
	Email         string // Join dari users
	NIS           string
	FullName      string
	Gender        string
	BirthPlace    *string
	BirthDate     *time.Time
	Phone         *string
	AddressRegion *string
	AddressDetail *string
	Photo         *string
	MajorID        *string // Dari fk majors
	ClassLevel     *string // Dari fk levels
	CurrentClassID *string // Ditambahkan dari JOIN student_classes
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// User struct hanya merepresentasikan model data pada users table saat Create.
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

type StudentClass struct {
	ID             string
	StudentID      string
	ClassID        string
	AcademicYearID string
	IsActive       bool
	CreatedAt      time.Time
}
