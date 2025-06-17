package randomchatters

import (
	ev "StreamTTS/EnvVariables"
	messagehandling "StreamTTS/MessageHandling"
	twichcomm "StreamTTS/TwichComm"
	"math/rand"
	"slices"
)

const _MESSAGES_BUF_LENGTH int = 16

type State struct {
	CurrentCahtter *twichcomm.UserInfo
	Messages chan string
}

var CurrentState *State = nil

func Init() {
	CurrentState = &State{
		CurrentCahtter: &twichcomm.UserInfo{UserLogin: "physickdev", UserName: "physickdev"},
		Messages: make(chan string, _MESSAGES_BUF_LENGTH),
	}

	messagehandling.RegisterHandler(&messagehandling.Handler{
		Condition: ChechForEvent, 
		Action: Action,
	})
}

var ChattersIgnore []string = []string{"physickdev"}

func GetRandomChatterID() *twichcomm.UserInfo {
	var users, err = twichcomm.GetStreamViewers(ev.Enviroment.BroadcasterId, ev.Enviroment.UserId)
	if err != nil {
		return nil 
	}
	for i, u := range users.Data {
		if slices.Contains(ChattersIgnore, u.UserLogin) {
			users.Data = slices.Delete(users.Data, i, i+1)
		}
	}
	return &users.Data[rand.Intn(len(users.Data))]
}

func ChechForEvent(username string, _ string) bool {
	// fmt.Printf("Comparing %v and %v - %v\n", username, CurrentState.CurrentCahtter.UserName, username == CurrentState.CurrentCahtter.UserName)
	return username == CurrentState.CurrentCahtter.UserName
}

func Action(_ string, message string) {
	CurrentState.Messages <- message
}
