package messagehandling

import (
	chatters "StreamTTS/Chatters"
	envvariables "StreamTTS/EnvVariables"
	models "StreamTTS/Models"
	"fmt"
)

func CreateTTSHandler() *Handler {
	return &Handler {
		Condition: ttsCondition,
		Action: ttsAction,
		Filtered: true,
	}
}

func ttsCondition(_ models.APIChatMessage) bool {
	return envvariables.Config.EnableTTS
}

func ttsAction(msg models.APIChatMessage) {
	var name = msg.Payload.Event.Broadcaster_User_Name
	var message = msg.Payload.Event.Message.Text

	var UserID = chatters.GetChatterID(name, envvariables.Enviroment.MainDB)
	if UserID < 0 {
		chatters.RegisterChatter(name, envvariables.Enviroment.MainDB)
		UserID = chatters.GetChatterID(name, envvariables.Enviroment.MainDB)
	}
	fmt.Printf("%d %v: %v\n", UserID, name, message)
	go SayMsg(UserID, message)
}
