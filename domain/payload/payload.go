package payload

type SendMessagePayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type ConnectMessagePayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}
