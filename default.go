package go_simple_router

import (
	"encoding/json"
	"errors"
	"net/http"
)

type defaultResponse struct {
	Header defaultResponseHeader `json:"header"`
	Data   interface{}           `json:"data"`
}

type defaultResponseHeader struct {
	IsSuccess bool   `json:"is_success"`
	Reason    string `json:"reason,omitempty"`
}

func defaultHandlerFunc(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, errors.New("nil handler")
}

func defaultHandlerErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}

	resp := &defaultResponse{
		Header: defaultResponseHeader{
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

func defaultHandlerSuccessFunc(w http.ResponseWriter, r *http.Request, data interface{}) {
	resp := &defaultResponse{
		Header: defaultResponseHeader{
			IsSuccess: true,
		},
		Data: data,
	}

	payload, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}
