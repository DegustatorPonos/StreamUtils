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
	ev.Enviroment.MainDB = chatters.EstablishDBConnection()
	messagehandling.LoadFilter()
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

	// Setting up broadcaster ID
	ev.Enviroment.BroadcasterId = twichcomm.GetChannelId(ev.Enviroment.BroadcasterLogin)
	
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
	http.HandleFunc("/api/filters/reload", messagehandling.ReloadBannedWordList) 
	if ev.Config.EnableRandomChatter {
		randomchatters.RegisterEndpoints()
	}
	http.ListenAndServe(":3000", nil)
}

func RegisterSubscriptions(sessionInfo *twichcomm.ConnectionInfo) {
	twichcomm.ClearSubscriptions()
	twichcomm.SubscribeToChat(sessionInfo)
}
