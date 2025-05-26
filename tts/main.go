package main

import (
	twichcomm "StreamTTS/TwichComm"
	ev "StreamTTS/EnvVariables"
	"fmt"
	"net/http"
)
func main() {
	go RunHTTPServer()
	var envErr = ev.ReadEnvVariables()
	if envErr != nil {
		return
	}

	var isUserTokenValid = twichcomm.AuthenticateApp() 
	// Waiting for auth 
	ev.Enviroment.UserToken = <- twichcomm.AuthKeyChan
	if !isUserTokenValid {
		ev.SetUserToken()
	}

	// fmt.Printf("User token: %v\n", ev.Enviroment.UserToken)
	var SessionInfo, connectionErr = twichcomm.ConnectToWs(ev.Enviroment.TwichAPIKey)
	if connectionErr != nil {
		fmt.Println(connectionErr.Error())
		return
	}
	fmt.Printf("Session ID: %v\n", SessionInfo.SessionId)
	RegisterSubscriptions(SessionInfo)

	for {
	}
}

func RunHTTPServer() {
	http.HandleFunc("/auth", twichcomm.AuthKeyHttpEndpoint) 
	http.ListenAndServe(":3000", nil)
}

func RegisterSubscriptions(sessionInfo *twichcomm.ConnectionInfo) {
	twichcomm.SubscribeToChat(sessionInfo)
}
