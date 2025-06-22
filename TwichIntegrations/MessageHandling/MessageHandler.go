package messagehandling

import (
	chatters "StreamTTS/Chatters"
	ev "StreamTTS/EnvVariables"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const FilterPath string = "../Filters.txt"

type Filter struct {
	BannedWords []string `json:"bannedwords"`
}

func (base *Filter) ChechString(inp string) bool {
	for _, v := range base.BannedWords {
		if v == "" {
			continue
		}
		var nomalized = strings.ToLower(inp)
		if strings.Contains(nomalized, v) {
			fmt.Printf("The message '%v' contaings a banned word '%v'\n", inp, v)
			return false
		}
	}
	return true
}

type Handler struct {
	// Inputs are username and message. 
	// If it returns true antion will be invoked
	Condition func(string, string) bool
	// This function will be called when the condition check is passed
	Action func(string, string)
	// Set to true if the word should go through 
	// additional filter before the invocation
	Filtered bool
}

var registeredActions []Handler = []Handler{}

var GlobalFilter Filter

func LoadFilter() {
	var contents, err = os.ReadFile(FilterPath)
	if err != nil {
		panic("Could not load filter file")

	}
	GlobalFilter = Filter{
		BannedWords: strings.Split(string(contents), "\n"),
	}
}

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
			if handler.Filtered && !GlobalFilter.ChechString(msg) { 
				continue
			}
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

func ReloadBannedWordList(w http.ResponseWriter, r *http.Request) {
	LoadFilter()
	var payload, err = json.Marshal(GlobalFilter)
	if err != nil {
		return
	}
	w.Write(payload)
}
