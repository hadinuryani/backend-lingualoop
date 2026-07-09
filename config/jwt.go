package config

import "strconv"

type JWT struct {
	Secret          string
	ExpirationHours int
}

func LoadJWTConfig() JWT {
	jwtExpiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	return JWT{
		Secret:          getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
		ExpirationHours: jwtExpiration,
	}
}
