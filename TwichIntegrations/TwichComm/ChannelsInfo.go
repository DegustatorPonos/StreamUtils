package twichcomm

import (
	"encoding/json"
	"fmt"
	"net/http"

	ev "StreamTTS/EnvVariables"
)

const UserIdentURL string = "https://api.twitch.tv/helix/users"
const ViewerListURL string = "https://api.twitch.tv/helix/chat/chatters"

type ChannelInfoResponce struct {
	Data []ChannelInfo `json:"data"`
}

type ChannelInfo struct {
	Id string `json:"id"`
	Login string `json:"login"`
	DisplayName string `json:"display_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type UserInfo struct {
	UserID string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName string `json:"user_name"`
}

type ViewerList struct {
	Data []UserInfo  `json:"data"`
	Total int `json:"total"`
}

func GetChannelId(login string) string {
	var ApiResp, ApiErr = GetChannelInfo(login)
	if ApiErr != nil || len(ApiResp.Data) < 1 {
		return ""
	}
	return ApiResp.Data[0].Id
}

func GetChannelInfo(channelName string) (*ChannelInfoResponce, error) {
	var client = http.Client{}
	var req, _ = http.NewRequest("GET", fmt.Sprintf("%v?login=%v", UserIdentURL, channelName), nil)
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
	var ApiResp = ChannelInfoResponce{}
	var jsonErr = json.Unmarshal(body, &ApiResp)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &ApiResp, nil
}

func GetStreamViewers(streamerID string, adminId string) (*ViewerList, error) {
	var client = http.Client{}
	var req, _ = http.NewRequest("GET", fmt.Sprintf("%v?broadcaster_id=%v&moderator_id=%v", ViewerListURL, streamerID, adminId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ev.Enviroment.UserToken))
	req.Header.Set("Client-Id", ev.Enviroment.TwichAPIKey)
	var resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	var body = parseResponce(resp)
	if ShowMessages {
		fmt.Print("Users: ")
		fmt.Println(string(body))
	}
	var ApiResp = ViewerList{}
	var jsonErr = json.Unmarshal(body, &ApiResp)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &ApiResp, nil
}
