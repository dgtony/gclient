package gclient

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ApiOperation struct {
	Endpoint, Method string
}
type OpMap map[string]ApiOperation

// supported REST API operations
var Operations = OpMap{
	"GetKey":            ApiOperation{Endpoint: "item", Method: "GET"},
	"GetKeySubKey":      ApiOperation{Endpoint: "item", Method: "GET"},
	"GetKeySubIndex":    ApiOperation{Endpoint: "item", Method: "GET"},
	"SetKey":            ApiOperation{Endpoint: "item", Method: "POST"},
	"RemoveKey":         ApiOperation{Endpoint: "item", Method: "DELETE"},
	"GetStoredKeys":     ApiOperation{Endpoint: "keys", Method: "GET"},
	"GetStoredKeysMask": ApiOperation{Endpoint: "keys", Method: "GET"},
}

func buildURL(baseURL, endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURL, endpoint)
}

func decodeError(respBody []byte) (ErrorResponse, error) {
	var errResp ErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		return ErrorResponse{}, errors.New(fmt.Sprintf("bad error model => %s", err))
	}
	return errResp, nil
}

func requestSucceed(responseStatus int) bool {
	if responseStatus > 199 && responseStatus < 300 {
		return true
	}
	return false
}
