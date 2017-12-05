package test

import (
	"bytes"
	"encoding/json"
	"net/url"
	"path"
	"strconv"
	"strings"
	"testing"

	"fmt"
	"net/http"
	"os"
	"sync"
)

type JSONResponse struct {
	Message string `json:"message"`
}


func doRequest(fnAppName, fnAppRoute string, contentType string, requestBody, responseBody interface{})  (*bytes.Buffer, *http.Response, error) {
	u := url.URL{
		Scheme: "http",
		Host:   Host(),
	}
	u.Path = path.Join(u.Path, "r", fnAppName, fnAppRoute)

	b, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, err
	}
	content := bytes.NewBuffer(b)
	output := &bytes.Buffer{}

	response, err := CallFN(u.String(), contentType, content, output, "POST", []string{})

	if err != nil {
		return nil, response, err
	}
	err = json.Unmarshal(output.Bytes(), responseBody)
	if err != nil {
		return nil, response, err
	}

	return output, response, nil
}

func callMultiple(times int, t *testing.T, s *SuiteSetup, fnRoute, fnImage,
	fnFormat string) {

	CreateApp(t, s.Context, s.Client, s.AppName, map[string]string{})
	CreateRoute(t, s.Context, s.Client, s.AppName, fnRoute, fnImage, "sync",
		fnFormat, s.RouteConfig, s.RouteHeaders)

	var wg *sync.WaitGroup
	wg.Add(times)

	go func() {
		defer wg.Done()

		requestBody := RandStringBytes(100)
		responseBody := &JSONResponse{}
		_, response, err := doRequest(s.AppName, fnRoute, "text/plain", requestBody, responseBody)
		if err != nil {
			t.Errorf("Got unexpected error: %v", err)
		}
		if response.StatusCode != 200 {
			t.Errorf("Status code assertion error.\n\tExpected: %v\n\tActual: %v",
				200, response.StatusCode)
		}
	}()

	wg.Wait()
	DeleteApp(t, s.Context, s.Client, s.AppName)
}

func callOnce(t *testing.T, s *SuiteSetup, fnRoute, fnImage,
	fnFormat string, requestBody interface{}, responseBody interface{}) (*bytes.Buffer, *http.Response) {

	CreateApp(t, s.Context, s.Client, s.AppName, map[string]string{})
	CreateRoute(t, s.Context, s.Client, s.AppName, fnRoute, fnImage, "sync",
		fnFormat, s.RouteConfig, s.RouteHeaders)

	output, response, err := doRequest(s.AppName, fnRoute, "application/json", requestBody, responseBody)
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	DeleteApp(t, s.Context, s.Client, s.AppName)

	return output, response
}

func TestFDKFormatSmallBody(t *testing.T) {

	FDKImage := os.Getenv("FDK_FUNCTION_IMAGE")
	if FDKImage == "" {
		t.Error("Please set FDK-based function image to test")
	}
	formats := []string{"http", "json"}

	helloJohnPayload := &struct {
		Name string `json:"name"`
	}{
		Name: "Jimmy",
	}
	helloJohnExpectedOutput := "Hello Jimmy"
	for _, format := range formats {

		// echo function:
		// payload:
		//    {
		//        "name": "John"
		//    }
		// response:
		//    "Hello John"
		t.Run(fmt.Sprintf("test-fdk-%v-small-body", format), func(t *testing.T) {

			t.Parallel()
			s := SetupDefaultSuite()
			route := fmt.Sprintf("/test-fdk-%v-format-small-body", format)

			responsePayload := &JSONResponse{}
			output, response := callOnce(t, s, route, FDKImage, format, helloJohnPayload, responsePayload)

			if !strings.Contains(helloJohnExpectedOutput, responsePayload.Message) {
				t.Errorf("Output assertion error.\n\tExpected: %v\n\tActual: %v", helloJohnExpectedOutput, output.String())
			}
			if response.StatusCode != 200 {
				t.Errorf("Status code assertion error.\n\tExpected: %v\n\tActual: %v", 200, response.StatusCode)
			}

			expectedHeaderNames := []string{"Content-Type", "Content-Length"}
			expectedHeaderValues := []string{"text/plain; charset=utf-8", strconv.Itoa(output.Len())}
			for i, name := range expectedHeaderNames {
				actual := response.Header.Get(name)
				expected := expectedHeaderValues[i]
				if !strings.Contains(expected, actual) {
					t.Errorf("HTTP header assertion error for %v."+
						"\n\tExpected: %v\n\tActual: %v", name, expected, actual)
				}
			}

		})
	}
}

func TestFDKMultipleEvents(t *testing.T) {

	FDKImage := os.Getenv("FDK_FUNCTION_IMAGE")
	if FDKImage == "" {
		t.Error("Please set FDK-based function image to test")
	}
	formats := []string{"http", "json"}

	for _, format := range formats {
		// this test attempts to send 100 concurrent requests
		// to a function in order to see if it's capable to handle more than 1 event
		// the only thing that matters in this test is response code, it should be 200 OK for all requests,
		// if one assertion fails means that FDK or Fn failed to dispatch necessary number of calls
		t.Run(fmt.Sprintf("test-fdk-%v-multiple-events", format), func(t *testing.T) {

			t.Parallel()
			s := SetupDefaultSuite()
			route := fmt.Sprintf("/test-fdk-%v-multiple-events", format)

			callMultiple(10, t, s, route, FDKImage, format)
		})
	}
}
