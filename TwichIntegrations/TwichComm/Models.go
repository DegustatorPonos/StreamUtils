package twichcomm

// This file contains structs only
type WelcomeMessage struct {
	Metadata APIMessageMetadata `json:"metadata"`
	Payload WelcomeMessagePayload  `json:"payload"`
}

type APIMessageMetadata struct {
	Message_Id string  `json:"message_id"`
	Message_type string `json:"message_type"`
	Subscription_Type string `json:"subscription_type"`
}

type WelcomeMessagePayload struct {
	Session WelcomeMessageSession  `json:"session"`
}

// used to determine the type of the message and decide on its processing later
type BareBonesMessage struct {
	Metadata APIMessageMetadata `json:"metadata"`
}

type WelcomeMessageSession struct {
	Id string  `json:"id"`
	Status string `json:"status"`
	Connected_at string  `json:"connected_at"`
	Keepalive_timeout_seconds int  `json:"keepalive_timeout_seconds"`
	Reconnect_url string `json:"reconnect_url"`
	Recovery_url string `json:"recovery_url"`
}

type WebhookInfo struct {
	Method string `json:"method"`
	Session_ID string `json:"session_id"`
}

type EventPaylaodData struct {
	Subscription SubscriptionInfo `json:"subscription"`
	Event MessageEventData `json:"event"`
}

type MessageEventData struct {
	Broadcaster_User_Id string `json:"broadcaster_user_id"`
	Broadcaster_User_Login string `json:"broadcaster_user_login"`
	Broadcaster_User_Name string `json:"broadcaster_user_name"`
	Source_Broadcaster_User_Id string  `json:"source_broadcaster_user_id"`
	Source_Broadcaster_User_Login string  `json:"source_broadcaster_user_login"`
	Source_Broadcaster_User_Name string  `json:"source_broadcaster_user_name"`
	Chatter_User_Id string   `json:"chatter_user_id"`
	Chatter_User_Login string   `json:"chatter_user_login"`
	Chatter_User_Name string   `json:"chatter_user_name"`
	MessageID string `json:"messageid"`
	Is_Source_Only string `json:"is_source_only"`
	Message ChatMessage `json:"message"`
}

type ChatMessage struct {
	Text string `json:"text"`
}

type Badge struct {
	Set_ID string  `json:"set_id"`
	Id string `json:"id"`
}

type SubscriptionInfo struct {
	Id string `json:"id"`
	Status string `json:"status"`
	Cost int `json:"cost"`
	Created_At string `json:"created_at"`
	Transport WebhookInfo `json:"transport"`
}

type APIChatMessage struct {
	Metadata APIMessageMetadata `json:"metadata"`
	Payload EventPaylaodData `json:"payload"`
}

type OAuthResponce struct {
	Client_ID string `json:"client_id"`
	Login string `json:"login"`
	Scopes []string `json:"scopes"`
	User_Id string `json:"user_id"`
	Expires_In int `json:"expires_in"`
}
