package handlers

import "svc-whatsapp/usecase"

type Handler struct {
	Contract *usecase.Contract
}

func NewHandler(ucContract *usecase.Contract) Handler {
	return Handler{Contract: ucContract}
}
