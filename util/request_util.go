package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	ContentTypeHeader string = "Content-Type"
	JsonContentType   string = "application/json"
	XOkapiTenant      string = "x-okapi-tenant"
	XOkapiToken       string = "x-okapi-token"
)

func DoGetDecodeReturnMapStringInteface(commandName string, url string, enableDebug bool, panicOnError bool, headers map[string]string) map[string]interface{} {
	var respMap map[string]interface{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error(commandName, "http.NewRequest error", "")
		panic(err)
	}
	AddRequestHeaders(req, headers)
	DumpHttpRequest(commandName, req, enableDebug)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if panicOnError {
			slog.Error(commandName, "http.DefaultClient.Do error", "")
			panic(err)
		} else {
			LogWarn(commandName, enableDebug, fmt.Sprintf("http.DefaultClient.Do warn - Endpoint is unreachable: %s", url))
			return nil
		}
	}
	defer func() {
		CheckStatusCodes(commandName, resp)
		resp.Body.Close()
	}()
	DumpHttpResponse(commandName, resp, enableDebug)
	err = json.NewDecoder(resp.Body).Decode(&respMap)
	if err != nil {
		if panicOnError {
			slog.Error(commandName, "json.NewDecoder error", "")
			panic(err)
		} else {
			LogWarn(commandName, enableDebug, fmt.Sprintf("json.NewDecoder warn - Cannot decode response from url: %s", url))
			return nil
		}
	}
	return respMap
}

func DoPostReturnMapStringInteface(commandName string, url string, enableDebug bool, bodyBytes []byte, headers map[string]string) map[string]interface{} {
	var respMap map[string]interface{}
	DumpHttpBody(commandName, enableDebug, bodyBytes)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		slog.Error(commandName, "http.NewRequest error", "")
		panic(err)
	}
	AddRequestHeaders(req, headers)
	DumpHttpRequest(commandName, req, enableDebug)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(commandName, "http.DefaultClient.Do error", "")
		panic(err)
	}
	defer func() {
		CheckStatusCodes(commandName, resp)
		resp.Body.Close()
	}()
	DumpHttpResponse(commandName, resp, enableDebug)
	err = json.NewDecoder(resp.Body).Decode(&respMap)
	if err != nil {
		slog.Error(commandName, "json.NewDecoder error", "")
		panic(err)
	}

	return respMap
}

func DoPutReturnNoContent(commandName string, url string, enableDebug bool, bodyBytes []byte, headers map[string]string) {
	DumpHttpBody(commandName, enableDebug, bodyBytes)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		slog.Error(commandName, "http.NewRequest error", "")
		panic(err)
	}
	AddRequestHeaders(req, headers)
	DumpHttpRequest(commandName, req, enableDebug)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(commandName, "http.DefaultClient.Do error", "")
		panic(err)
	}
	defer func() {
		CheckStatusCodes(commandName, resp)
		resp.Body.Close()
	}()
	DumpHttpResponse(commandName, resp, enableDebug)
}

func AddRequestHeaders(req *http.Request, headers map[string]string) {
	if len(headers) == 0 {
		req.Header.Add(ContentTypeHeader, JsonContentType)
		return
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}
