package go_simple_router

import (
	"encoding/json"
	"net/http"
)

type DefaultResponse struct {
	Header DefaultResponseHeader `json:"header"`
	Data   interface{}           `json:"data"`
}

type DefaultResponseHeader struct {
	IsSuccess bool   `json:"is_success"`
	Reason    string `json:"reason,omitempty"`
}

func DefaultErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}

	resp := &DefaultResponse{
		Header: DefaultResponseHeader{
			IsSuccess: false,
			Reason:    errorMessage,
		},
		Data: nil,
	}

	payload, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(payload)
}

func DefaultSuccessFunc(w http.ResponseWriter, r *http.Request, data interface{}) {
	resp := &DefaultResponse{
		Header: DefaultResponseHeader{
			IsSuccess: true,
		},
		Data: data,
	}

	payload, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}
