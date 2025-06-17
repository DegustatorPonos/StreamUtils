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

type ActiveSubscriptionsList struct {
	Data []struct{ 
		Id string `json:"id"`
	} `json:"data"`
}

func OnExit() {
	fmt.Println("xdx")
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
	if ShowMessages {
		fmt.Printf("Subscription request: %v\n", string(bodyJSON))
	}
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
	var respBody = parseResponce(resp)
	fmt.Printf("Subscription responce: %v\n", string(respBody))
	return true
}

func ClearSubscriptions() {
	var subs = GetActiveSubscriptions()
	var client = &http.Client{}
	for _, v := range subs {
		var req, _ = http.NewRequest("DELETE", fmt.Sprintf("%v?id=%v", EventSubURL, v), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ev.Enviroment.UserToken))
		req.Header.Set("Client-Id", ev.Enviroment.TwichAPIKey)
		var resp, err = client.Do(req)
		if err != nil || resp.StatusCode != 204 {
			fmt.Printf("Error while unsubbing from event %v. Error code: %v\n", v, resp.StatusCode)
		}
	}
}

func GetActiveSubscriptions() []string {
	var client = &http.Client{}
	var req, _ = http.NewRequest("GET", EventSubURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ev.Enviroment.UserToken))
	req.Header.Set("Client-Id", ev.Enviroment.TwichAPIKey)
	var resp, err = client.Do(req)
	if err != nil {
		return nil
	}
	var body = parseResponce(resp)
	if ShowMessages {
		fmt.Println(string(body))
	}
	var subList = ActiveSubscriptionsList{}
	json.Unmarshal(body, &subList)
	var outp = make([]string, 0)
	for _, v := range subList.Data {
		outp = append(outp, v.Id)
	}
	return outp
}
