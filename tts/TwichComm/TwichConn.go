package twichcomm

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

const TwichWS_URI string = "wss://eventsub.wss.twitch.tv/ws"
const Origin string = "http://localhost"

// struct

func ConnectToWs(ApiKey string) error {
	fmt.Println("Connecting to twich API")
	var wsConn, err = websocket.Dial(TwichWS_URI, "", Origin)
	if err != nil {
		return fmt.Errorf("An error occured while connecting to twich API. \nOriginal error: %v\n", err.Error())
	}
	// Reading welocme message
	var welcome_msg, wm_err = readWelcomeMessage(wsConn)
	if wm_err != nil {
		return wm_err
	}
	fmt.Printf("MessageID: '%v'\n", welcome_msg.Metadata.Message_Id)

	go ConnectionRoutine(wsConn)
	return nil
}

func ConnectionRoutine(ws *websocket.Conn) {
	defer ws.Close()
	for {
		var buf = make([]byte, 1024)
		var i, err = ws.Read(buf) 
		if err != nil {
			continue
		}
		buf = buf[:i]
		fmt.Printf("Welcome message: %v \n", string(buf))
	}
}

func readWelcomeMessage(ws *websocket.Conn) (*WelcomeMessage, error) {
	var buf = make([]byte, 1024)
	var i int
	var err error 
	if i, err = ws.Read(buf); err != nil {
		return nil, err
	}
	buf = buf[:i]
	var welcome_msg WelcomeMessage
	var unmErr = json.Unmarshal([]byte(string(buf)), &welcome_msg)
	if unmErr != nil {
		return nil, unmErr
	}
	return &welcome_msg, nil
}
