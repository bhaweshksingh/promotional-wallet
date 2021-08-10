package utils

import (
	"account/pkg/http/contract"
	"account/pkg/http/internal/resperr"
	"encoding/json"
	"net/http"
)

func WriteSuccessResponse(resp http.ResponseWriter, statusCode int, data interface{}) {
	writeAPIResponse(resp, statusCode, contract.NewSuccessResponse(data))
}

func WriteFailureResponse(resp http.ResponseWriter, err resperr.ResponseError) {
	writeAPIResponse(resp, err.StatusCode(), contract.NewFailureResponse(err.Description()))
}

func writeAPIResponse(resp http.ResponseWriter, code int, ar contract.APIResponse) {
	b, err := json.Marshal(&ar)
	if err != nil {
		writeResponse(resp, http.StatusInternalServerError, []byte("server error"))
		return
	}

	writeResponse(resp, code, b)
}

func writeResponse(response http.ResponseWriter, statusCode int, body []byte) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	_, _ = response.Write(body)
}
