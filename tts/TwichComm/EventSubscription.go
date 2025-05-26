package twichcomm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	ev "StreamTTS/EnvVariables"
)

const EventSubURL string = "https://api.twitch.tv/helix/eventsub/subscriptions"

type subscriptionRequest struct {
	Type string `json:"type"`
	Version string `json:"version"`
	Condition any  `json:"condition"`
	Transport WebhookInfo `json:"transport"`
}

type ChatMessageCondition struct {
	Broadcaster_user_ID string `json:"broadcaster_user_id"`
	User_ID string `json:"user_id"`
}

type WebhookInfo struct {
	Method string `json:"method"`
	Session_ID string `json:"session_id"`
}

func SubscribeToChat(sessionInfo *ConnectionInfo) bool {
	var body = subscriptionRequest {
		Type: "channel.chat.message",
		Version: "1",
		Condition: ChatMessageCondition{
			Broadcaster_user_ID: ev.Enviroment.BroadcasterId,
			User_ID: ev.Enviroment.UserId,
		},
		Transport: WebhookInfo{
			Method: "websocket",
			Session_ID: sessionInfo.SessionId,
		},
	}
	var bodyJSON, err = json.Marshal(body)
	if err != nil {
		return false
	}
	fmt.Printf("Sub request: %v\n", string(bodyJSON))
	var client = &http.Client{}
	var req, reqErr = http.NewRequest("POST", EventSubURL, bytes.NewReader(bodyJSON))
	if reqErr != nil {
		return false
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ev.Enviroment.UserToken))
	req.Header.Set("Client-Id", ev.Enviroment.TwichAPIKey)
	req.Header.Set("Content-Type", "application/json")
	var resp, doErr = client.Do(req)
	if doErr != nil {
		return false
	}

	fmt.Println(req.Header.Get("Authorization"))

	var respBody = parseResponce(resp)
	fmt.Printf("Subscription request: %v\n", string(respBody))
	
	return true
}
