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

func doRequest(t *testing.T, fnAppName, fnAppRoute string, contentType string, requestBody interface{}) (*bytes.Buffer, *http.Response, error) {
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

	return output, response, nil
}

func callMultiple(wg *sync.WaitGroup, times int, t *testing.T, s *SuiteSetup, fnRoute, fnImage,
	fnFormat string) {

	timeout := int32(30)
	idleTimeout := int32(10)
	CreateApp(t, s.Context, s.Client, s.AppName, map[string]string{})
	CreateRoute(t, s.Context, s.Client, s.AppName, fnRoute, fnImage, "sync",
		fnFormat, timeout, idleTimeout, s.RouteConfig, s.RouteHeaders)

	wg.Add(times)

	for i := 0; i < times; i++ {
		go func() {
			defer wg.Done()

			requestBody := fmt.Sprintf(`{"name":"%v"}`, RandStringBytes(100))
			output, response, err := doRequest(t, s.AppName, fnRoute, "text/plain", requestBody)
			if err != nil {
				t.Errorf("Got unexpected error: %v", err)
			}
			if response.StatusCode != 200 {
				t.Log(output.String())
				t.Errorf("Status code assertion error.\n\tExpected: %v\n\tActual: %v",
					200, response.StatusCode)
			}
		}()
	}
	wg.Wait()

	DeleteApp(t, s.Context, s.Client, s.AppName)
}

func callOnce(t *testing.T, s *SuiteSetup, fnRoute, fnImage,
	fnFormat string, requestBody interface{}) (*bytes.Buffer, *http.Response, error) {

	timeout := int32(30)
	idleTimeout := int32(10)
	CreateApp(t, s.Context, s.Client, s.AppName, map[string]string{})
	CreateRoute(t, s.Context, s.Client, s.AppName, fnRoute, fnImage, "sync",
		fnFormat, timeout, idleTimeout, s.RouteConfig, s.RouteHeaders)

	output, response, err := doRequest(t, s.AppName, fnRoute, "application/json", requestBody)
	if err != nil {
		return nil, response, err
	}

	DeleteApp(t, s.Context, s.Client, s.AppName)

	return output, response, nil
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
		//        "name": "Jimmy"
		//    }
		// response:
		//    "Hello Jimmy"
		t.Run(fmt.Sprintf("test-fdk-%v-small-body", format), func(t *testing.T) {

			t.Parallel()
			s := SetupDefaultSuite()
			route := fmt.Sprintf("/test-fdk-%v-format-small-body", format)

			output, response, err := callOnce(t, s, route, FDKImage, format, helloJohnPayload)
			if err != nil {
				t.Errorf("unexpected error: %v", err.Error())
			}
			if !strings.Contains(output.String(), helloJohnExpectedOutput) {
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
	var wg sync.WaitGroup

	for _, format := range formats {
		// this test attempts to send 10 concurrent requests
		// to a function in order to see if it's capable to handle more than 1 event
		// the only thing that matters in this test is response code, it should be 200 OK for all requests,
		// if one assertion fails means that FDK or Fn failed to dispatch necessary number of calls
		t.Run(fmt.Sprintf("test-fdk-%v-multiple-events", format), func(t *testing.T) {

			t.Parallel()
			s := SetupDefaultSuite()
			route := fmt.Sprintf("/test-fdk-%v-multiple-events", format)

			callMultiple(&wg, 100, t, s, route, FDKImage, format)
		})
	}
}
