package keys

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetEnv(str string) string {
	envPath, err := filepath.Abs("internals/server/keys/.env")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	switch str {
	case "key":
		return os.Getenv("key")
	case "login":
		return os.Getenv("ADMIN_USERNAME")
	case "password":
		return os.Getenv("ADMIN_PASSWORD")
	default:
		return ""
	}
}
