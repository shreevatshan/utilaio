package utilaio

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	PROTOCOL_HTTP  = "http"
	PROTOCOL_HTTPS = "https"

	REQUEST_TYPE_GET  = "GET"
	REQUEST_TYPE_POST = "POST"

	CONTENT_TYPE_APPLICATION_JSON = "application/json"

	DEFAULT_REQUEST_TIMEOUT_SECONDS = 10
	HTTP_STATUS_OK                  = 200
)

type HRequest struct {
	Time_out             int
	Enable_proxy         bool
	Download_response    bool
	Body                 []byte
	Url                  *url.URL
	Proxy                *url.URL
	Insecure_skip_verify bool
	Api                  string
	Request_type         string
	Content_type         string
	Download_location    string
	Headers              map[string]string
}

type HResponse struct {
	Code int
	Body []byte
	Err  error
}

func formAPI(Url *url.URL) (string, error) {
	api := Url.String()

	parsed_url, err := url.ParseRequestURI(api)
	if err != nil {
		return "", err
	} else if len(parsed_url.Hostname()) == 0 {
		return "", fmt.Errorf("not a valid HTTP request uri (%s)", api)
	} else {
		return api, nil
	}

}

func (request *HRequest) validateRequest() error {
	if len(request.Request_type) == 0 {
		return errors.New("request type not set")
	}

	if len(request.Api) == 0 {
		var err error
		request.Api, err = formAPI(request.Url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (response *HResponse) downloadResponse(resp *http.Response, request HRequest) {

	var err error

	defer resp.Body.Close()
	out, err := os.Create(request.Download_location)
	if err != nil {
		response.Err = fmt.Errorf("error creating file for package [%v]", err)
		return
	}

	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		response.Err = fmt.Errorf("error saving package [%v]", err)
		return
	}
}

func SendRequest(request HRequest) HResponse {
	var response HResponse
	var req *http.Request
	var err error

	if request.Time_out == 0 {
		request.Time_out = DEFAULT_REQUEST_TIMEOUT_SECONDS
	}

	request.validateRequest()
	if err != nil {
		response.Err = err
		return response
	}

	transport := http.Transport{
		Proxy:           http.ProxyURL(request.Proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: request.Insecure_skip_verify}}

	client := http.Client{
		Timeout:   time.Duration(request.Time_out) * time.Second,
		Transport: &transport}

	req, err = http.NewRequest(request.Request_type, request.Api, bytes.NewBuffer(request.Body))
	if err != nil {
		response.Err = err
		return response
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	http_response, err := client.Do(req)
	if err != nil {
		response.Err = err
		return response
	}

	response.Code = http_response.StatusCode

	if http_response.StatusCode == HTTP_STATUS_OK {
		if request.Download_response {
			response.downloadResponse(http_response, request)
		} else {
			defer http_response.Body.Close()
			response.Body, err = io.ReadAll(http_response.Body)
			if err != nil {
				response.Err = err
				return response
			}
		}
	} else {
		response.Err = fmt.Errorf("response code not OK (%d)", response.Code)
	}

	return response
}
