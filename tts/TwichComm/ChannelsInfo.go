package twichcomm

import (
	"encoding/json"
	"fmt"
	"net/http"

	ev "StreamTTS/EnvVariables"
)

const UserIdentURL string = "https://api.twitch.tv/helix/users"

type ChannelApiResp struct {
	Data []ChannelInfo `json:"data"`
}

type ChannelInfo struct {
	Id string `json:"id"`
	Login string `json:"login"`
	Display_Name string `json:"display_name"`
}

func GetChannelId(login string) string {
	var ApiResp, ApiErr = GetChannelInfo(login)
	if ApiErr != nil || len(ApiResp.Data) < 1 {
		return ""
	}
	return ApiResp.Data[0].Id
}

func GetChannelInfo(login string) (*ChannelApiResp, error) {
	var client = http.Client{}
	var req, _ = http.NewRequest("GET", fmt.Sprintf("%v?login=%v", UserIdentURL, login), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ev.Enviroment.UserToken))
	req.Header.Set("Client-Id", ev.Enviroment.TwichAPIKey)
	var resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	var body = parseResponce(resp)
	if ShowMessages {
		fmt.Print("Channel info: ")
		fmt.Println(string(body))
	}
	var ApiResp = ChannelApiResp{}
	var jsonErr = json.Unmarshal(body, &ApiResp)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &ApiResp, nil
}
