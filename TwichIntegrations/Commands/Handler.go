package commands

import (
	messagehandling "StreamTTS/MessageHandling"
	models "StreamTTS/Models"
	twichcomm "StreamTTS/TwichComm"
	"fmt"
	"strings"
)

func CreateComandsHandler() *messagehandling.Handler {
	return &messagehandling.Handler {
		Condition: commandCondition,
		Action: commandAction,
		Filtered: true,
	}
}

func separateCommandFromMessage(msg models.APIChatMessage) string {
	var messageParts = strings.Split(msg.Payload.Event.Message.Text, " ")
	if (len(messageParts) == 0) {
		return ""
	}
	return messageParts[0]
}

func commandAction(msg models.APIChatMessage) {
	// The loadedMessage is guaranteed to be loaded
	var msgRaw = separateCommandFromMessage(msg)
	var commandResp, exists = getComand(msgRaw)
	if (!exists) {
		commandResp = fmt.Sprintf("%s is not a thing around here", 
			strings.TrimLeft(msgRaw, "!"))
	}
	var err = twichcomm.SendMessage(commandResp, msg.Payload.Event.MessageID)
	if err != nil {
		fmt.Printf("Command execution error: %s", err.Error())
	}
}

func commandCondition(msg models.APIChatMessage) bool {
	var message = msg.Payload.Event.Message.Text
	return len(message) != 0 && message[0] == '!'
	// fmt.Printf("Checking message %s\n", message)
	// var _, contains = getComand(separateCommandFromMessage(msg))
	// return contains
}
