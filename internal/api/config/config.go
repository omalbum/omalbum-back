package config

import (
	"os"
)

const (
	AuthKey string = "lionelmessi"

	Port string = ":8080"

	DevelopmentDatabaseHost     string = "localhost"
	DevelopmentDatabaseDatabase string = "teleoma"
	DevelopmentDatabaseUser     string = "teleoma"
	DevelopmentDatabasePassword string = "teleoma"
)

func getEnv() string {
	return os.Getenv("ENVIRONMENT")
}

func IsProd() bool {
	return getEnv() == "PROD"
}

func IsStage() bool {
	return getEnv() == "STAGE"
}

func GetSendGridApiKey() string {
	return os.Getenv("SENDGRID_API_KEY")
}

type Database struct {
	Host     string
	Database string
	User     string
	Password string
}

func GetDatabase() Database {
	if IsProd() || IsStage() {
		return Database{
			Host:     os.Getenv("DATABASE_HOST"),
			Database: os.Getenv("DATABASE_DATABASE"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
		}
	}

	return Database{
		Host:     DevelopmentDatabaseHost,
		Database: DevelopmentDatabaseDatabase,
		User:     DevelopmentDatabaseUser,
		Password: DevelopmentDatabasePassword,
	}
}

func IsCorsEnabled() bool {
	cors := os.Getenv("CORS")
	return cors == "ENABLED"
}

func ShouldSendMails() bool {
	return IsProd()
}

func GetMailsPath() string {
	if IsProd() || IsStage() {
		return os.Getenv("MAILS_PATH")
	}

	return "static/mails/"
}
