package commands

import (
	messagehandling "StreamTTS/MessageHandling"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const commandsLocation = "Commands.json"

var loadedCommands = make(map[string]string)
var helpCommand = ""

func readCommandsListFile(path string) error {
	var file, fopenerr = os.ReadFile(path)
	if fopenerr != nil {
		return fmt.Errorf("Failed to read commands list: %s", fopenerr.Error())
	}
	var jsonErr = json.Unmarshal(file, &loadedCommands)
	if jsonErr != nil {
		return fmt.Errorf("Failed to parse commands list: %s", jsonErr.Error())
	}
	helpCommand = formHelpCommand()
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

func formHelpCommand() string {
	var sb = strings.Builder{}
	sb.WriteString("Avaliable commands: ")
	for k := range loadedCommands {
		sb.WriteString(k)
		sb.WriteString(", ")
	}
	return sb.String()
}

func getComand(msg string) (string, bool) {
	switch msg {
	case "!help":
		return helpCommand, true
	}
	var found, exists = loadedCommands[msg]
	return found, exists
}
