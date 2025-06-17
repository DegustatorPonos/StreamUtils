package messagehandling

import (
	chatters "StreamTTS/Chatters"
	ev "StreamTTS/EnvVariables"
	"fmt"
	"os/exec"
)

type Handler struct {
	// Inputs are username and message. 
	// If it returns true antion will be invoked
	Condition func(string, string) bool
	// This function will be called when the condition check is passed
	Action func(string, string)
}

var registeredActions []Handler = []Handler{}

func HandleMessage(username, msg string) {
	var UserID = chatters.GetChatterID(username, ev.Enviroment.MainDB)
	if UserID < 0 {
		chatters.RegisterChatter(username, ev.Enviroment.MainDB)
		UserID = chatters.GetChatterID(username, ev.Enviroment.MainDB)
	}
	fmt.Printf("%d %v: %v\n", UserID, username, msg)
	if ev.Config.EnableTTS {
		go SayMsg(UserID, msg)
	}

	for _, handler := range registeredActions {
		if handler.Condition(username, msg) {
			handler.Action(username, msg)
		}
	}
}

func SayMsg(chatterID int, msg string) {
	var v = chatters.GetVoice(ev.Enviroment.MainDB, chatterID)
	var speedArg = fmt.Sprintf("%d", v.Speed)
	var pitchArg = fmt.Sprintf("%d", v.Pitch)
	var capArg = fmt.Sprintf("-k%d", v.Capital)
	// var cmd = exec.Command("espeak", fmt.Sprintf("\"%v\"", msg), "&")
	var cmd = exec.Command("espeak", fmt.Sprintf("\"%v\"", msg), "-s", speedArg, "-p", pitchArg, capArg, "&")
	cmd.Run()
}

func RegisterHandler(handler *Handler) {
	registeredActions = append(registeredActions, *handler)
}
