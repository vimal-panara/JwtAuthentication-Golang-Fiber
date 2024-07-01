package initializer

import (
	"JwtAuthentication/database"
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvFile() {
	fmt.Println("loading .env file")
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func ConnectToDatabase() {
	fmt.Println("connecting to database")
	if err := database.ConnectToDatabase(); err != nil {
		panic(err)
	}
}
