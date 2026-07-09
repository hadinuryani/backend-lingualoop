package config

type Gmail struct {
	Email       string
	AppPassword string
}

func LoadGmailConfig() Gmail {
	return Gmail{
		Email:       getEnv("GMAIL_EMAIL", ""),
		AppPassword: getEnv("GMAIL_APP_PASSWORD", ""),
	}
}
