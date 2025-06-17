package envvariables

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	// TTS will say all the messages
	EnableTTS bool
	// Enables random chatter calls functionality
	EnableRandomChatter bool
}

type EnvVariables struct {
	// Stored in .env file
	TwichAPIKey string
	TwichAPISecret string
	BroadcasterLogin string
	UserCode string
	// Stored in temp
	UserToken string
	// Got every time
	WsSessionID string
	BroadcasterId string
	UserId string
	MainDB *sql.DB
}

// Functionality
var Config AppConfig = AppConfig{
	EnableTTS: false,
	EnableRandomChatter: true,
}

var Enviroment EnvVariables = EnvVariables{}

// Reads variables from .env file
func ReadEnvVariables() (error) {
	var loadErr = godotenv.Load()
	if loadErr != nil {
		return loadErr
	}
	Enviroment.TwichAPIKey = os.Getenv("TWICH_API_KEY")
	Enviroment.UserCode = os.Getenv("USER_CODE")
	Enviroment.TwichAPISecret = os.Getenv("TWICH_API_SECRET")
	Enviroment.BroadcasterLogin = os.Getenv("BROADCASTER_LOGIN")
	return nil
}

func SetUserToken() {
	var newLine = fmt.Sprintf("USER_CODE=%v", Enviroment.UserCode)
	var content, ferr = os.ReadFile(".env")
	if ferr != nil {
		return
	}
	var vars = strings.Split(string(content), "\n")
	for i, v := range vars {
		if strings.HasPrefix(v, "USER_CODE=") {
			vars[i] = newLine
			writeToEnv(vars)
			return
		}
	}
	vars = append(vars, newLine)
	writeToEnv(vars)
	// var _ = os.Setenv("USER_TOKEN", Enviroment.UserToken)
}

func writeToEnv(val []string) {
	var contents = strings.Join(val, "\n")
	var err = os.WriteFile(".env", []byte(contents), 0666)
	if err != nil {
		fmt.Println("Error while writing user code to .env file")
	}
}
