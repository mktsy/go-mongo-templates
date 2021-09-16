package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func readEnv(val string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return strings.ReplaceAll(os.Getenv(val), "\n", "")
}

var (
	MongoDBProtocol = readEnv("MONGO_DB_PROTOCOL")
	MongoDBHost     = readEnv("MONGO_DB_HOST")
	MongoDBUsername = readEnv("MONGO_DB_USERNAME")
	MongoDBPassword = readEnv("MONGO_DB_PASSWORD")
	MongoDBName     = readEnv("MONGO_DB_NAME")
	MongoDBOptions  = readEnv("MONGO_DB_OPTIONS")
	Port            = readEnv("PORT")
)
