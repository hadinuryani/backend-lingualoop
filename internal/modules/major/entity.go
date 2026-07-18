package major

import "time"

// Major adalah representasi Entity dari tabel `majors` di database.
type Major struct {
	ID          string
	Code        string
	Name        string
	Description string
	LogoFileID  *string
	LogoPath    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
