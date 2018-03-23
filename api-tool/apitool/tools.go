package apitool

import (
	"bytes"
	"encoding/json"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/apihelper"
	"github.com/gin-contrib/cors"
	"github.com/kr/pretty"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const Week = time.Hour * 24 * 7

// RequestHeader is the object to send to the HttpRequestHandlers
type RequestHeader struct {
	URL           string
	Method        string
	ContentType   string
	Authorization string
}

// WaitForServerToStart return true only when API is ready
func WaitForServerToStart(url string) bool {
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", url, nil)
		client := &http.Client{}
		if _, err := client.Do(req); err == nil {
			return true
		}
		time.Sleep(2500)
	}
	return false
}

// Generate CORS config for router
func DefaultCORSConfig() cors.Config {
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowCredentials = true
	CORSConfig.AllowAllOrigins = true
	CORSConfig.AllowHeaders = []string{"*"}
	CORSConfig.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}
	return CORSConfig
}

// HttpRequestHandler is a handler for easy http requests
func HttpRequestHandler(requestHeader RequestHeader, requestBody interface{}, responseDTO interface{}) (*http.Response, apihelper.ApiErrors) {
	req, err := http.NewRequest(requestHeader.Method, requestHeader.URL, encodeRequestBody(requestBody))
	req.Header.Set("Content-Type", requestHeader.ContentType)
	req.Header.Set("Authorization", requestHeader.Authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	var apiErrors apihelper.ApiErrors
	json.Unmarshal(body, &responseDTO)
	json.Unmarshal(body, &apiErrors)

	if len(apiErrors.Errors) > 0 {
		pretty.Println(apiErrors.Errors)
	}

	return resp, apiErrors
}

// HttpRequestHandlerForUnitTesting is the same as HttpRequestHandler but add logging for tests
func HttpRequestHandlerForUnitTesting(t *testing.T, requestHeader RequestHeader, requestBody interface{}, responseDTO interface{}) (*http.Response, apihelper.ApiErrors) {

	// Log header
	t.Log(requestHeader)

	// Create request
	req, err := http.NewRequest(requestHeader.Method, requestHeader.URL, encodeRequestBodyAndLog(t, requestBody))
	req.Header.Set("Content-Type", requestHeader.ContentType)
	req.Header.Set("Authorization", requestHeader.Authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read response
	body, _ := ioutil.ReadAll(resp.Body)

	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &responseDTO)
	json.Unmarshal(body, &apiErrors)

	// Log response
	t.Log(responseDTO)
	t.Log(apiErrors)
	if len(apiErrors.Errors) > 0 {
		pretty.Println(apiErrors.Errors)
	}

	return resp, apiErrors
}

type ChatfuelError struct {
	SetAttributes apihelper.Error `json:"set_attributes"`
}

// HttpRequestHandlerForChatfuelUnitTesting is the same as HttpRequestHandler but add logging for tests
func HttpRequestHandlerForChatfuelUnitTesting(t *testing.T, requestHeader RequestHeader, requestBody interface{}, responseDTO interface{}) (*http.Response, ChatfuelError) {

	// Log header
	t.Log(requestHeader)

	// Create request
	req, err := http.NewRequest(requestHeader.Method, requestHeader.URL, encodeRequestBodyAndLog(t, requestBody))
	req.Header.Set("Content-Type", requestHeader.ContentType)
	req.Header.Set("Authorization", requestHeader.Authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read response
	body, _ := ioutil.ReadAll(resp.Body)

	var chatfuelError ChatfuelError
	json.Unmarshal(body, &responseDTO)
	json.Unmarshal(body, &chatfuelError)

	// Log response
	t.Log(responseDTO)
	t.Log(chatfuelError)

	return resp, chatfuelError
}

func encodeRequestBody(reqBody interface{}) io.Reader {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	return b
}

func encodeRequestBodyAndLog(t *testing.T, reqBody interface{}) io.Reader {
	t.Log("testing with following parameters:")
	t.Log(reqBody)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	return b
}

// ExtractQueryParams format arguments for gorm from a map[string]interface{}
func ExtractQueryParams(queryParams map[string]interface{}) (query string, args []interface{}) {
	for key, value := range queryParams {
		if query == "" {
			query += key
		} else {
			query += " AND " + key
		}
		args = append(args, value)
	}
	return
}
