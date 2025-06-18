package randomchatters

import (
	ev "StreamTTS/EnvVariables"
	messagehandling "StreamTTS/MessageHandling"
	twichcomm "StreamTTS/TwichComm"
	"encoding/json"
	"fmt"
	"math/rand"
	"slices"
)

const _MESSAGES_BUF_LENGTH int = 16

type State struct {
	CurrentCahtter *twichcomm.ChannelInfo
	Messages chan string
}

type MessageEvent struct {
	Type string `json:"type"`
	Message string `json:"message"`
}

type DisconnectEvent struct {
	Type string `json:"type"`
}

type ConnectEvent struct {
	Type string `json:"type"`
	UserName string `json:"username"`
	UserPfp string `json:"userpfp"`
}

var CurrentState *State = nil

func Init() {
	CurrentState = &State{
		CurrentCahtter: &twichcomm.ChannelInfo{Login: "physickdev", DisplayName: "physickdev"},
		Messages: make(chan string, _MESSAGES_BUF_LENGTH),
	}

	messagehandling.RegisterHandler(&messagehandling.Handler{
		Condition: HandlerCondition, 
		Action: HandlerAction,
	})
}

var IgnoredChatters []string = []string{"physickdev", "PaketikCrew"}

func GetRandomChatterID() *twichcomm.UserInfo {
	var users, err = twichcomm.GetStreamViewers(ev.Enviroment.BroadcasterId, ev.Enviroment.UserId)
	if err != nil {
		return nil 
	}
	for i, u := range users.Data {
		if slices.Contains(IgnoredChatters, u.UserLogin) {
			users.Data = slices.Delete(users.Data, i, i+1)
		}
	}
	return &users.Data[rand.Intn(len(users.Data))]
}

// Is called when new chatter is selected through API
// New user info should be filled when invoking this function
func onConnect() {
	var event = ConnectEvent{
		Type: "connect",
		UserName: CurrentState.CurrentCahtter.DisplayName,
		UserPfp: CurrentState.CurrentCahtter.ProfileImageUrl,
	}
	var payload, marshalErr = json.Marshal(event)
	if marshalErr != nil {
		return
	}
	sendPayloadToWS(&payload)
}

// Is called when new chatter is disconnectd through API
func onDisconnect() {
	var event = DisconnectEvent{
		Type: "disconnect",
	}
	var payload, marshalErr = json.Marshal(event)
	if marshalErr != nil {
		return
	}
	sendPayloadToWS(&payload)
}

func HandlerCondition(username string, _ string) bool {
	if CurrentState.CurrentCahtter == nil {
		return false
	}
	// fmt.Printf("Comparing %v and %v - %v\n", username, CurrentState.CurrentCahtter.UserName, username == CurrentState.CurrentCahtter.UserName)
	return username == CurrentState.CurrentCahtter.DisplayName
}

func HandlerAction(_ string, message string) {
	var event = MessageEvent{Type: "message", Message: message}
	var payload, marshalErr = json.Marshal(event)
	if marshalErr != nil {
		return
	}
	sendPayloadToWS(&payload)
}

func sendPayloadToWS(payload *[]byte) {
	for i, conn := range _WSConnections {
		if conn == nil {
			fmt.Println("Deleted closed connection")
			slices.Delete(_WSConnections, i, i+1)
			continue
		}
		var _, err = conn.Write(*payload)
		if err != nil {
			fmt.Println("Deleted broken connection")
			slices.Delete(_WSConnections, i, i+1)
		}
	}
}
