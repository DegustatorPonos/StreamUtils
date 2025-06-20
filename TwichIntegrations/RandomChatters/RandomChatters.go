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
		CurrentCahtter: nil,
		Messages: make(chan string, _MESSAGES_BUF_LENGTH),
	}

	messagehandling.RegisterHandler(&messagehandling.Handler{
		Condition: HandlerCondition, 
		Action: HandlerAction,
	})
}

var IgnoredChatters []string = []string{"physickdev", "personthemanhumane"}

func GetRandomChatter() *twichcomm.ChannelInfo {
	var users, err = twichcomm.GetStreamViewers(ev.Enviroment.BroadcasterId, ev.Enviroment.UserId)
	if err != nil {
		return nil 
	}
	var possible = make([]string, 0)
	for _, u := range users.Data {
		if !(slices.Contains(IgnoredChatters, u.UserName) || slices.Contains(IgnoredChatters, u.UserLogin)) {
			possible = append(possible, u.UserLogin)
		}
	}
	if len(possible) == 0 {
		return &twichcomm.ChannelInfo{
			DisplayName: "Nobody",
		}
	}
	var userData, dataErr = twichcomm.GetChannelInfo(possible[rand.Intn(len(possible))])
	if dataErr != nil {
		fmt.Printf("An error occured during ercieving user data. Original error: %s\n", dataErr.Error())
		return &twichcomm.ChannelInfo{
			DisplayName: "Nobody",
		}
	}
	return &userData.Data[0]
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
	var toDelete = make([]int, 0)
	for i, conn := range _WSConnections {
		if conn == nil {
			fmt.Println("Deleted closed connection")
			toDelete = append(toDelete, i)
			continue
		}
		var _, err = conn.Write(*payload)
		if err != nil {
			fmt.Println("Deleted faulty connection")
			toDelete = append(toDelete, i)
		}
	}
	for i, v := range toDelete {
		_WSConnections = append(_WSConnections[:v-i], _WSConnections[v-i+1:]...)
	}
}
