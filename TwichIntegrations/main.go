package main

import (
	"fmt"
	"net/http"

	chatters "StreamTTS/Chatters"
	ev "StreamTTS/EnvVariables"
	randomchatters "StreamTTS/RandomChatters"
	twichcomm "StreamTTS/TwichComm"
)

var TerminationChan = make(chan interface{}, 1)

func main() {
	go RunHTTPServer()
	ev.Enviroment.MainDB = chatters.EstablishDBConnection()
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
	
	var SessionInfo, connectionErr = twichcomm.ConnectToWs(ev.Enviroment.TwichAPIKey)
	if connectionErr != nil {
		fmt.Println(connectionErr.Error())
		return
	}

	if ev.Config.EnableRandomChatter {
		randomchatters.Init()
	}

	RegisterSubscriptions(SessionInfo)

	<- TerminationChan
}

func RunHTTPServer() {
	http.HandleFunc("/auth", twichcomm.AuthKeyHttpEndpoint) 
	if ev.Config.EnableRandomChatter {
		randomchatters.RegisterEndpoints()
	}
	http.ListenAndServe(":3000", nil)
}

func RegisterSubscriptions(sessionInfo *twichcomm.ConnectionInfo) {
	twichcomm.ClearSubscriptions()
	twichcomm.SubscribeToChat(sessionInfo)
}
