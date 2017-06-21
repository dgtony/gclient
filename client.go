package gclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ERR_CODE_VALUE_NOT_FOUND = 21
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Smart client constructor
func NewClient(cacheBaseURL string) *Client {
	client := Client{HTTPClient: &http.Client{}}
	if len(cacheBaseURL) > 7 {
		client.BaseURL = cacheBaseURL
	} else {
		// default base
		client.BaseURL = "http://localhost:8080"
	}
	return &client
}

/* client interface */

// Get value from cache with given key and operation timeout in seconds.
// Timeout could be disabled by setting to zero
// Return tuple: (value, foundFlag, error)
func (c *Client) Get(key string, timeout int) (interface{}, bool, error) {
	req := GetKeyRequest{Key: key}
	respStatus, respBody, err := c.makeRequest("GetKey", req, timeout)
	if err != nil {
		return nil, false, err
	}
	var succRespBox GetKeyResponse
	return processResponseValue(respStatus, respBody, &succRespBox, ERR_CODE_VALUE_NOT_FOUND)
}

// Get value from map stored with key using sub-key
func (c *Client) GetSubKey(key string, subKey string, timeout int) (interface{}, bool, error) {
	req := GetKeySubKeyRequest{Key: key, SubKey: subKey}
	respStatus, respBody, err := c.makeRequest("GetKeySubKey", req, timeout)
	if err != nil {
		return nil, false, err
	}
	var succRespBox GetKeySubResponse
	return processResponseValue(respStatus, respBody, &succRespBox, ERR_CODE_VALUE_NOT_FOUND)
}

// Get value from list stored with key using sub-index.
// Indexing in stored list starts from 1.
func (c *Client) GetSubIndex(key string, subIndex int, timeout int) (interface{}, bool, error) {
	req := GetKeySubIndexRequest{Key: key, SubIndex: subIndex}
	respStatus, respBody, err := c.makeRequest("GetKeySubIndex", req, timeout)
	if err != nil {
		return nil, false, err
	}
	var succRespBox GetKeySubResponse
	return processResponseValue(respStatus, respBody, &succRespBox, ERR_CODE_VALUE_NOT_FOUND)
}

// Save new item with given key and TTL.
// Item value could be a string, map with string keys or arbitrary slice.
// Parameter key TTL sets lifetime for the key in seconds.
func (c *Client) Set(key string, value interface{}, keyttl int, timeout int) error {
	req := SetKeyRequest{
		Key:   key,
		Value: value,
		TTL:   keyttl}
	respStatus, respBody, err := c.makeRequest("SetKey", req, timeout)
	if err != nil {
		return err
	}
	if requestSucceed(respStatus) {
		return nil
	}
	errResp, err := decodeError(respBody)
	if err != nil {
		return errors.New(fmt.Sprintf("bad error model => %s", err))
	}
	return errors.New(fmt.Sprintf("cache error => %s", errResp.Reason))
}

// Remove stored key from cache.
func (c *Client) Remove(key string, timeout int) error {
	req := RemoveKeyRequest{Key: key}
	respStatus, respBody, err := c.makeRequest("RemoveKey", req, timeout)
	if err != nil {
		return err
	}
	if requestSucceed(respStatus) {
		return nil
	}
	errResp, err := decodeError(respBody)
	if err != nil {
		return errors.New(fmt.Sprintf("bad error model => %s", err))
	}
	return errors.New(fmt.Sprintf("cache error => %s", errResp.Reason))
}

// Return slice of all keys currently stored in cache.
func (c *Client) Keys(timeout int) ([]string, error) {
	respStatus, respBody, err := c.makeRequest("GetStoredKeys", nil, timeout)
	if err != nil {
		return nil, err
	}
	var succRespBox GetStoredKeysResponse
	return processResponseKeys(respStatus, respBody, &succRespBox)
}

// Return slice of stored keys according to given mask.
// Mask use glob pattern matching rules.
func (c *Client) KeysMask(mask string, timeout int) ([]string, error) {
	req := GetStoredKeysRequest{Mask: mask}
	respStatus, respBody, err := c.makeRequest("GetStoredKeysMask", req, timeout)
	if err != nil {
		return nil, err
	}
	var succRespBox GetStoredKeysResponse
	return processResponseKeys(respStatus, respBody, &succRespBox)
}

/* internals */

func (c *Client) String() string {
	return fmt.Sprintf("client => cache location: %s", c.BaseURL)
}

func (c *Client) makeRequest(opName string, reqBox interface{}, timeout int) (int, []byte, error) {
	var req *http.Request
	var buff bytes.Buffer

	// get endpoint location
	opInfo, ok := Operations[opName]
	if !ok {
		return 0, nil, errors.New(fmt.Sprintf("no operation info => %s", opName))
	}
	url := buildURL(c.BaseURL, opInfo.Endpoint)

	// prepare request
	if reqBox != nil {
		if err := json.NewEncoder(&buff).Encode(reqBox); err != nil {
			return 0, nil, errors.New(fmt.Sprintf("cannot encode request => %s", err))
		}
		req, _ = http.NewRequest(opInfo.Method, url, &buff)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(opInfo.Method, url, nil)
	}
	req.Header.Set("Accept", "application/json")

	// add cancellation timeout
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
	}

	// dial
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("cannot make request => %s", err))
	}

	// drain buffer and close to reuse connection
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("cannot read response body => %s", err))
	}
	return resp.StatusCode, body, nil
}

// return tuple: (value, succeedFlag, error)
func processResponseValue(respStatus int, respBody []byte, succRespBox ValueReader, failCode int) (interface{}, bool, error) {
	if requestSucceed(respStatus) {
		if err := json.Unmarshal(respBody, &succRespBox); err != nil {
			return nil, false, errors.New(fmt.Sprintf("bad response model => %s", err))
		}
		return succRespBox.GetValue(), true, nil

	} else {
		errResp, err := decodeError(respBody)
		if err != nil {
			return nil, false, errors.New(fmt.Sprintf("bad error model => %s", err))
		}
		if errResp.Code == failCode {
			return nil, false, nil
		} else {
			return nil, false, errors.New(fmt.Sprintf("cache error => %s", errResp.Reason))
		}
	}
}

// return tuple: (keyList, error)
func processResponseKeys(respStatus int, respBody []byte, succRespBox KeyReader) ([]string, error) {
	if requestSucceed(respStatus) {
		if err := json.Unmarshal(respBody, &succRespBox); err != nil {
			return nil, errors.New(fmt.Sprintf("bad response model => %s", err))
		}
		return succRespBox.GetKeys(), nil
	} else {
		errResp, err := decodeError(respBody)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("bad error model => %s", err))
		}
		return nil, errors.New(fmt.Sprintf("cache error => %s", errResp.Reason))
	}
}
