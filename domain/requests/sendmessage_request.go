package requests

type SendMessageRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}
