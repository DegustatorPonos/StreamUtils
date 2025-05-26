package twichcomm

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
	ev "StreamTTS/EnvVariables"
	msgs "StreamTTS/MessageHandling"
)

const TwichWS_URI string = "wss://eventsub.wss.twitch.tv/ws"
const Origin string = "http://localhost"

const ShowMessages bool = true

type ConnectionInfo struct {
	SessionId string
}

func ConnectToWs(ApiKey string) (*ConnectionInfo, error) {
	fmt.Println("Connecting to twich API")
	var wsConn, err = websocket.Dial(TwichWS_URI, "", Origin)
	if err != nil {
		return nil, fmt.Errorf("An error occured while connecting to twich API. \nOriginal error: %v\n", err.Error())
	}

	// Reading welocme message
	var welcome_msg, wm_err = readWelcomeMessage(wsConn)
	if wm_err != nil {
		return nil, wm_err
	}
	ev.Enviroment.WsSessionID = welcome_msg.Payload.Session.Id
	// fmt.Printf("MessageID: '%v'\n", welcome_msg.Metadata.Message_Id)

	go ConnectionRoutine(wsConn)
	return &ConnectionInfo{SessionId: welcome_msg.Payload.Session.Id}, nil
}

func ConnectionRoutine(ws *websocket.Conn) {
	defer ws.Close()
	for {
		var buf = make([]byte, 2048)
		var i, err = ws.Read(buf) 
		if err != nil {
			continue
		}
		buf = buf[:i]
		var MsgType = GetMessageType(buf)
		// fmt.Printf("Msg type: %v\n", MsgType)
		switch MsgType {
		default:
			continue
		case "session_keepalive":
			continue
		case "channel.chat.message":
			var msg APIChatMessage
			var unmErr = json.Unmarshal(buf, &msg)
			if unmErr != nil {
				fmt.Printf("Error: %v\n", unmErr.Error())
				continue
			}
			// fmt.Println(msg.Payload.Event.Message.Text)
			msgs.HandleMessage(msg.Payload.Event.Chatter_User_Name, msg.Payload.Event.Message.Text)
			// fmt.Printf("Message: %v \n", string(buf))
		}
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
	if ShowMessages {
		fmt.Println(string(buf))
	}
	var unmErr = json.Unmarshal([]byte(string(buf)), &welcome_msg)
	if unmErr != nil {
		return nil, unmErr
	}
	return &welcome_msg, nil
}

func GetMessageType(data []byte) string {
	var bbm = BareBonesMessage{}
	var err = json.Unmarshal(data, &bbm)
	if err != nil {
		return ""
	}
	if bbm.Metadata.Subscription_Type != "" {
		return bbm.Metadata.Subscription_Type
	}
	return bbm.Metadata.Message_type
}
