package envvariables

import (
	"crypto"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	// Stored in .env file
	TwichAPIKey string
	TwichAPISecret string
	UserCode string
	AppAPIKey string
	// Stored in temp
	UserToken string
	// Got every time
	WsSessionID string
	BroadcasterId string
	UserId string
	MainDB *sql.DB
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
	// Enviroment.BroadcasterLogin = os.Getenv("BROADCASTER_LOGIN")
	Enviroment.AppAPIKey = os.Getenv("APP_TOKEN")
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

// Regenerates token used to aceess this app's API
func RegenerateAPIKey() {
	var newToken = crypto.SHA256.New()
	var newLine = fmt.Sprintf("APP_TOKEN=%v", base64.RawURLEncoding.EncodeToString(newToken.Sum(nil)))
	var content, ferr = os.ReadFile(".env")
	if ferr != nil {
		return
	}
	var vars = strings.Split(string(content), "\n")
	for i, v := range vars {
		if strings.HasPrefix(v, "APP_TOKEN=") {
			vars[i] = newLine
			writeToEnv(vars)
			return
		}
	}
	vars = append(vars, newLine)
	writeToEnv(vars)
}

func writeToEnv(val []string) {
	var contents = strings.Join(val, "\n")
	var err = os.WriteFile(".env", []byte(contents), 0666)
	if err != nil {
		fmt.Println("Error while writing user code to .env file")
	}
}
