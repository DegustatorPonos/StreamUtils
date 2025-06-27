package messagehandling

import (
	"fmt"
	"os"
	"strings"

	ev "StreamTTS/EnvVariables"
)


type Filter struct {
	BannedWords []string `json:"bannedwords"`
}

func LoadFilter() {
	var contents, err = os.ReadFile(ev.Config.FilterSettings.FilterFileLocation)
	if err != nil {
		panic("Could not load filter file")
	}
	GlobalFilter = Filter{
		BannedWords: strings.Split(string(contents), "\n"),
	}
}

func (base *Filter) CheckString(inp string) bool {
	if ev.Config.FilterSettings.BannedWords {
		// TODO: REwrite this shit
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
	}
	return true
}
