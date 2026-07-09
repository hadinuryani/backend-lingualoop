package auth

type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	FullName     string
	Role         string
	AvatarURL    string
	IsActive     bool
}
