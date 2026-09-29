package commands

import (
	messagehandling "StreamTTS/MessageHandling"
	models "StreamTTS/Models"
	twichcomm "StreamTTS/TwichComm"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const commandsLocation = "Commands.json"

var loadedCommands = make(map[string]string)

func readCommandsListFile(path string) error {
	var file, fopenerr = os.ReadFile(path)
	if fopenerr != nil {
		return fmt.Errorf("Failed to read commands list: %s", fopenerr.Error())
	}
	var jsonErr = json.Unmarshal(file, &loadedCommands)
	if jsonErr != nil {
		return fmt.Errorf("Failed to parse commands list: %s", jsonErr.Error())
	}
	return nil
}

func InitComands() error {
	var err = readCommandsListFile(commandsLocation)
	if err != nil {
		return err
	}
	messagehandling.RegisterHandler(CreateComandsHandler())
	return nil
}

func CreateComandsHandler() *messagehandling.Handler {
	return &messagehandling.Handler {
		Condition: commandCondition,
		Action: commandAction,
		Filtered: true,
	}
}

func commandAction(msg models.APIChatMessage) {
	// The loadedMessage is guaranteed to be loaded
	var messageParts = strings.Split(msg.Payload.Event.Message.Text, " ")
	if (len(messageParts) == 0) {
		return
	}
	var commandResp, _ = loadedCommands[messageParts[0]]
	var err = twichcomm.SendMessage(commandResp, msg.Payload.Event.MessageID)
	if err != nil {
		fmt.Printf("Command execution error: %s", err.Error())
	}
}

func commandCondition(msg models.APIChatMessage) bool {
	var message = msg.Payload.Event.Message.Text
	fmt.Printf("Checking message %s\n", message)
	var _, contains = loadedCommands[message]
	return contains
}
