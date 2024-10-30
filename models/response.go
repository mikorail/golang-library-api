package models

type ApiResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count,omitempty"` // Optional for lists
}
