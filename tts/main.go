package main

import (
	"fmt"
	"net/http"

	ev "StreamTTS/EnvVariables"
	twichcomm "StreamTTS/TwichComm"
)

var TerminationChan = make(chan interface{}, 1)

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

	// Setting up user ID
	if twichcomm.IsTokenValid() {
		fmt.Printf("User ID: '%v'\n", ev.Enviroment.UserId)
	} else {
		panic("Token is somehow invalid")
	}

	// Setting up broadcaster ID
	fmt.Printf("Broadcaster: %v\n", ev.Enviroment.BroadcasterLogin)
	var broadcasterID = twichcomm.GetChannelId(ev.Enviroment.BroadcasterLogin)
	fmt.Printf("Broadcaster ID: '%v'\n", broadcasterID)
	ev.Enviroment.BroadcasterId = broadcasterID
	
	// fmt.Printf("User token: %v\n", ev.Enviroment.UserToken)
	var SessionInfo, connectionErr = twichcomm.ConnectToWs(ev.Enviroment.TwichAPIKey)
	if connectionErr != nil {
		fmt.Println(connectionErr.Error())
		return
	}
	// fmt.Printf("Session ID: %v\n", SessionInfo.SessionId)
	RegisterSubscriptions(SessionInfo)

	<- TerminationChan
}

func RunHTTPServer() {
	http.HandleFunc("/auth", twichcomm.AuthKeyHttpEndpoint) 
	http.ListenAndServe(":3000", nil)
}

func RegisterSubscriptions(sessionInfo *twichcomm.ConnectionInfo) {
	twichcomm.ClearSubscriptions()
	twichcomm.SubscribeToChat(sessionInfo)
}
