package messagehandling

import (
	chatters "StreamTTS/Chatters"
	ev "StreamTTS/EnvVariables"
	"fmt"
	"os/exec"
)

func HandleMessage(username, msg string) {
	var UserID = chatters.GetChatterID(username, ev.Enviroment.MainDB)
	if UserID < 0 {
		chatters.RegisterChatter(username, ev.Enviroment.MainDB)
		UserID = chatters.GetChatterID(username, ev.Enviroment.MainDB)
	}
	fmt.Printf("%d %v: %v\n", UserID, username, msg)
	go SayMsg(msg)
}

func SayMsg(msg string) {
	var cmd = exec.Command("espeak", fmt.Sprintf("\"%v\"", msg), "&")
	cmd.Run()
}
