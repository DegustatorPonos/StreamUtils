package main

import (
	twichcomm "StreamTTS/TwichComm"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	TwichAPIKey string
}

func main() {
	var EnvVars, envErr = ReadEnvVariables()
	if envErr != nil {
		return
	}
	fmt.Printf("API key: %v\n", EnvVars.TwichAPIKey)

	var err = twichcomm.ConnectToWs(EnvVars.TwichAPIKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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
