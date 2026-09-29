package twichcomm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	ev "StreamTTS/EnvVariables"
)

const sendChatURL = "https://api.twitch.tv/helix/chat/messages"

type chatMessageReq struct {
	BroadcasterId string `json:"broadcaster_id"`
	SenderId string `json:"sender_id"`
	Message string `json:"message"`
	ReplyParentMessageId string `json:"reply_parent_message_id,omitempty"`
}

type chatMessageResult struct {
	Data struct {
		MessageId string `json:"message_id"`
		IsSent bool `json:"is_sent"`
		DropReason *struct {
			Code string `json:"code"`
			Message string `json:"message"`
		} `json:"dropreason"`
	} `json:"data"`
}

func SendMessage(message string, replyMessageId string) error {
	var data = chatMessageReq {
		BroadcasterId: ev.Enviroment.BroadcasterId,
		SenderId: ev.Enviroment.UserId,
		Message: message,
		ReplyParentMessageId: replyMessageId,
	}
	var body, jsonErr = json.Marshal(data)
	if jsonErr != nil {
		return fmt.Errorf("Failed to form a message body: %s", jsonErr.Error())
	}

	fmt.Printf("Sending %s\n", body)

	var client = &http.Client{}
	var req, _ = http.NewRequest("POST", sendChatURL, bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ev.Enviroment.UserToken))
	req.Header.Set("Client-Id", ev.Enviroment.TwichAPIKey)
	req.Header.Set("Content-Type", "application/json")

	var resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send a message request: %s", err.Error()) 
	}
	var respBody = parseResponce(resp)

	fmt.Printf("Got %s\n", string(respBody))

	var res chatMessageResult
	jsonErr = json.Unmarshal(respBody, &res)
	if jsonErr != nil {
		return fmt.Errorf("Failed to parse a message responce body: %s", jsonErr.Error())
	}

	if (!res.Data.IsSent && res.Data.DropReason != nil) {
		return fmt.Errorf("Failed to send a message. Code: '%s'. Reason: %s", 
		res.Data.DropReason.Code, 
		res.Data.DropReason.Message)
	}

	return nil
}
