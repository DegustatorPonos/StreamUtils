package messagehandling

import (
	models "StreamTTS/Models"
	"fmt"
)

const commandsLocation = "Commands.json"

var loadedCommands map[string]string = nil

type comandList struct {
	Commands map[string]string `json:"comands"`
}

func InitComands() error {
	return nil
}

func CreateComandsHandler() *Handler {
	return &Handler {
		Condition: commandCondition,
		Action: commandAction,
		Filtered: true,
	}
}

func commandAction(msg models.APIChatMessage) {
	// The loadedMessage is guaranteed to be loaded
	var message = msg.Payload.Event.Message.Text
	var commmand, _ = loadedCommands[message]
	fmt.Printf("Responce: %s\n", commmand)
}

func commandCondition(msg models.APIChatMessage) bool {
	if loadedCommands == nil {
		return false
	}
	var message = msg.Payload.Event.Message.Text
	var _, contains = loadedCommands[message]
	return contains
}
