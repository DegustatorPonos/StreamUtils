package envvariables

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const ConfigLocation string = "Config.json"

type AppConfig struct {
	// TTS will say all the messages
	EnableTTS bool `json:"enabletts"`
	// Enables random chatter calls functionality
	EnableRandomChatter bool `json:"enablerandomchatter"`
	// A calculation of chatter activity
	ActivityMetrics bool `json:"activitymetrics"`
	// The target channel
	BroadcasterLogin string `json:"broadcasterlogin"`
}

var Config AppConfig

var defaultConfig AppConfig = AppConfig{
	EnableTTS: true,
	EnableRandomChatter: true,
	ActivityMetrics: true,
	BroadcasterLogin: "physickdev",
}

func LoadConfig() error {
	var config, err = readConfig()
	if err != nil {
		return err
	}
	Config = config
	return nil
}

func readConfig() (AppConfig, error) {
	var file, fopenerr = os.ReadFile(ConfigLocation)
	if fopenerr != nil {
		var createerr, inintBody = createConfigFile()
		if createerr != nil {
			file = inintBody
		} else {
			return AppConfig{}, fmt.Errorf("Unable to load or create settings file. \nOriginal error: %v\n", fopenerr.Error())
		}
	}
	var data = AppConfig{}
	json.Unmarshal(file, &data)
	Config = data
	return data, nil
}

func createConfigFile() (error, []byte) {
	var file, err = os.Create(ConfigLocation)
	if err != nil {
		return err, nil
	}
	defer file.Close()
	var inititalData = defaultConfig
	var body, jsonerr = json.MarshalIndent(inititalData, "", "	")
	if jsonerr != nil {
		return jsonerr, nil
	}
	fmt.Fprint(file, string(body))
	return nil, body
}

func RelaoadConfigEndpoint(w http.ResponseWriter, r *http.Request) {
	var newConfig, err = readConfig()
	if err != nil {
		fmt.Printf("Unable to load or create settings file. \nOriginal error: %v\n", err.Error())
		w.WriteHeader(500)
		return
	}
	Config = newConfig
}
