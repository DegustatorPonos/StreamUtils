package messagehandling

import (
	chatters "StreamTTS/Chatters"
	envvariables "StreamTTS/EnvVariables"
	"fmt"
)

func CreateTTSHandler() *Handler {
	return &Handler {
		Condition: ttsCondition,
		Action: ttsAction,
		Filtered: true,
	}
}

func ttsCondition(_ string, _ string) bool {
	return true
}

func ttsAction(name string, message string) {
	var UserID = chatters.GetChatterID(name, envvariables.Enviroment.MainDB)
	if UserID < 0 {
		chatters.RegisterChatter(name, envvariables.Enviroment.MainDB)
		UserID = chatters.GetChatterID(name, envvariables.Enviroment.MainDB)
	}
	fmt.Printf("%d %v: %v\n", UserID, name, message)
	SayMsg(UserID, message)
}
