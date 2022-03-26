package payload

type SendMessagePayload struct {
	ApiKey  string `json:"hex"`
	To      string `json:"to"`
	Message string `json:"message"`
}
