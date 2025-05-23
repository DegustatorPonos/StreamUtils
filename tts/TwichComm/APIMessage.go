package twichcomm

// This file contains structs only
type WelcomeMessage struct {
	Metadata APIMessageMetadata `json:"metadata"`
	Payload WelcomeMessagePayload  `json:"payload"`
}

type APIMessageMetadata struct {
	Message_Id string  `json:"message_id"`
	Message_type string `json:"message_type"`
}

type WelcomeMessagePayload struct {
	Session WelcomeMessageSession  `json:"session"`
}

type WelcomeMessageSession struct {
	Id string  `json:"id"`
	Status string `json:"status"`
	Connected_at string  `json:"connected_at"`
	Keepalive_timeout_seconds int  `json:"keepalive_timeout_seconds"`
	Reconnect_url string `json:"reconnect_url"`
	Recovery_url string `json:"recovery_url"`
}
