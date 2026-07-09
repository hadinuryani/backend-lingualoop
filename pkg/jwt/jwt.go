package jwt

import (
	"errors"
	"fmt"
	"time"

	"backend-lingualoop/config"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"` // "admin", "teacher", "student"
	gojwt.RegisteredClaims
}

func GenerateToken(userID, email, username, role string) (string, error) {
	cfg := config.GetConfig()

	expirationHours := cfg.JWT.ExpirationHours
	if expirationHours <= 0 {
		expirationHours = 24
	}

	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: gojwt.RegisteredClaims{
			Issuer:    "lingualoop",
			Subject:   userID,
			IssuedAt:  gojwt.NewNumericDate(now),
			ExpiresAt: gojwt.NewNumericDate(now.Add(time.Duration(expirationHours) * time.Hour)),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("gagal membuat token: %w", err)
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token tidak boleh kosong")
	}

	cfg := config.GetConfig()

	token, err := gojwt.ParseWithClaims(tokenString, &Claims{}, func(token *gojwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak valid: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token tidak valid: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid atau sudah expired")
	}

	return claims, nil
}