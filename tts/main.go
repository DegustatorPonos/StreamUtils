package main

import (
	twichcomm "StreamTTS/TwichComm"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	TwichAPIKey string
	UserToken string
}

func main() {
	go RunHTTPServer()
	var EnvVars, envErr = ReadEnvVariables()
	if envErr != nil {
		return
	}

	twichcomm.PrintAuthRequest(EnvVars.TwichAPIKey);
	// Waiting for auth 
	EnvVars.UserToken = <- twichcomm.AuthKeyChan
	fmt.Printf("User token: %v\n", EnvVars.UserToken)
	var SessionInfo, connectionErr = twichcomm.ConnectToWs(EnvVars.TwichAPIKey)
	if connectionErr != nil {
		fmt.Println(connectionErr.Error())
		return
	}
	fmt.Printf("Session ID: %v\n", SessionInfo.SessionId)

	for {
	}
}

// Reads variables from .env file
func ReadEnvVariables() (*EnvVariables, error) {
	var loadErr = godotenv.Load()
	if loadErr != nil {
		return nil, loadErr
	}
	var outp = EnvVariables{}
	outp.TwichAPIKey = os.Getenv("TWICH_API_KEY")
	return &outp, nil
}

func RunHTTPServer() {
	http.HandleFunc("/auth", twichcomm.AuthKeyHttpEndpoint) 
	http.ListenAndServe(":3000", nil)
}
