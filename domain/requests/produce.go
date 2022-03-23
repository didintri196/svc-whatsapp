package requests

type ProduceReq struct {
	Topic string      `json:"topic"`
	Body  interface{} `json:"body"`
}
