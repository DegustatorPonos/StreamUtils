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
	go SayMsg(UserID, msg)
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
