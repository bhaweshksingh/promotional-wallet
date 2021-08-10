package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseRequest(req *http.Request, data interface{}) (err error) {
	if req == nil {
		fmt.Errorf("ParseRequest. Error: request is nil")
	}

	if req.Body == nil {
		fmt.Errorf("ParseRequest. Error: request body  is nil")
	}

	decoder := json.NewDecoder(req.Body)

	err = decoder.Decode(&data)
	if err != nil {
		fmt.Errorf("ParseRequest.decoder.Decode. Error: %v", err)
	}

	return
}
