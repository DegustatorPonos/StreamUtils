package randomchatters

import (
	ev "StreamTTS/EnvVariables"
	twichcomm "StreamTTS/TwichComm"
	"math/rand"
	"slices"
)

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
