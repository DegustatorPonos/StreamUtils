package envvariables

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	// Stored in .env file
	TwichAPIKey string
	TwichAPISecret string
	UserToken string
	// Stored in temp
	UserCode string
}

var Enviroment EnvVariables = EnvVariables{}

// Reads variables from .env file
func ReadEnvVariables() (error) {
	var loadErr = godotenv.Load()
	if loadErr != nil {
		return loadErr
	}
	Enviroment.TwichAPIKey = os.Getenv("TWICH_API_KEY")
	Enviroment.UserToken = os.Getenv("USER_TOKEN")
	Enviroment.TwichAPISecret = os.Getenv("TWICH_API_SECRET")
	return nil
}

func SetUserToken() {
	var _ = os.Setenv("USER_TOKEN", Enviroment.UserToken)
}

