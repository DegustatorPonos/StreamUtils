package main

import (
	"fmt"
	"net/http"

	chatters "StreamTTS/Chatters"
	ev "StreamTTS/EnvVariables"
	messagehandling "StreamTTS/MessageHandling"
	randomchatters "StreamTTS/RandomChatters"
	twichcomm "StreamTTS/TwichComm"
)

var TerminationChan = make(chan interface{}, 1)

func main() {
	go RunHTTPServer()

	var envErr = ev.ReadEnvVariables()
	if envErr != nil {
		return
	}

	// Setting up this app's API key
	if ev.Enviroment.AppAPIKey == "" {
		ev.RegenerateAPIKey()
	}

	var isUserTokenValid = twichcomm.AuthenticateApp() 
	// Waiting for auth 
	ev.Enviroment.UserToken = <- twichcomm.AuthKeyChan
	if !isUserTokenValid {
		ev.SetUserToken()
	}

	// Double-checking user token
	if !twichcomm.IsTokenValid() {
		panic("Token is somehow invalid")
	}

	// We are authenticated

	if ev.Config.EnableRandomChatter {
		randomchatters.Init()
	}
	if ev.Config.ActivityMetrics {
		messagehandling.Init()
	}

	ev.Enviroment.MainDB = chatters.EstablishDBConnection()

	// Setting up broadcaster ID
	ev.Enviroment.BroadcasterId = twichcomm.GetChannelId(ev.Enviroment.BroadcasterLogin)
	
	var SessionInfo, connectionErr = twichcomm.ConnectToWs(ev.Enviroment.TwichAPIKey)
	if connectionErr != nil {
		fmt.Println(connectionErr.Error())
		return
	}
	RegisterSubscriptions(SessionInfo)

	<- TerminationChan
}

func RunHTTPServer() {
	http.HandleFunc("/auth", twichcomm.AuthKeyHttpEndpoint) 
	if ev.Config.EnableRandomChatter {
		randomchatters.RegisterEndpoints()
	}
	if ev.Config.ActivityMetrics {
		messagehandling.RegisterEndpoints()
	}
	http.ListenAndServe(":3000", nil)
}

func RegisterSubscriptions(sessionInfo *twichcomm.ConnectionInfo) {
	twichcomm.ClearSubscriptions()
	twichcomm.SubscribeToChat(sessionInfo)
}
