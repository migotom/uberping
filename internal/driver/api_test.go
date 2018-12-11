package driver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/migotom/uberping/internal/schema"
)

type ErrorTestAuthorize func(apiClient, error)
type ErrorTestRequest func([]byte, error)

func TestAuthorize(t *testing.T) {
	cases := []struct {
		Name    string
		Handler http.HandlerFunc
		Test    ErrorTestAuthorize
	}{
		{
			Name: "InvalidCredentials",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			Test: func(client apiClient, err error) {
				if !strings.HasPrefix(err.Error(), "401 Unauthorized") {
					t.Errorf("Expecting 401 Unauthorized, got: %v", err)
				}
			},
		},
		{
			Name: "InvalidAPIResponse",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "some invalid API response")
			},
			Test: func(client apiClient, err error) {
				if !strings.HasPrefix(err.Error(), "Invalid JSON response") {
					t.Errorf("Expecting JSON parsing failure, got: %v", err)
				}
			},
		},
		{
			Name: "ValidAPIResponse",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{
					"token": "some_token",
					"id_server": 1
				}`)
			},
			Test: func(client apiClient, err error) {
				if err != nil {
					t.Errorf("Expecting correct response, got: %v", err)
				}
				if client.authData.Token != "some_token" {
					t.Errorf("Expecting token, got: %s", client.authData.Token)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tc.Handler))
			defer ts.Close()
			var client = apiClient{
				apiConfig: &schema.APIConfig{
					Name:   "test",
					Secret: "secret",
				},
			}
			client.apiConfig.URL = ts.URL

			tc.Test(client, client.authorize())
		})
	}
}

func TestRequest(t *testing.T) {
	var tries int

	cases := []struct {
		Name    string
		Handler http.HandlerFunc
		Test    ErrorTestRequest
	}{
		{
			Name: "WithInvalidToken",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				tries++
				w.WriteHeader(http.StatusUnauthorized)
			},
			Test: func(body []byte, err error) {
				if !strings.HasPrefix(err.Error(), "request retry limit exceeded") {
					t.Errorf("Expecting request retry limit exceeded, got: %v", err)
				}

				// 3 tries, each one try also fetch auth token
				if tries != 6 {
					t.Errorf("Expecting 3 tries, got: %v", tries)
				}
			},
		},
		{
			Name: "WithValidToken",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				tries++
				fmt.Fprintln(w, "some API response")
			},
			Test: func(body []byte, err error) {
				if err != nil {
					t.Errorf("Expecting valid response, got: %v", err)
				}

				if tries != 1 {
					t.Errorf("Expecting 1 try, got: %v", tries)
				}

				if string(body) != "some API response\n" {
					t.Errorf("Expecting correct response, got: %v", string(body))
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tries = 0
			ts := httptest.NewServer(http.HandlerFunc(tc.Handler))
			defer ts.Close()
			var client = apiClient{
				apiConfig: &schema.APIConfig{
					Name:   "test",
					Secret: "secret",
				},
			}
			client.apiConfig.URL = ts.URL

			tc.Test(client.request("GET", ts.URL, nil))
		})
	}
}

func TestAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.String() {
		case "/auth":
			fmt.Fprintln(w, `
			{
				"token": "valid_token",
				"id_server": 1
			}`)
		case "/devices/1":
			fmt.Fprintf(w, `[
			{
				"id": 10,
				"name": "Dev",
				"id_server": 1,
				"ip": "192.168.12.150/24",
				"average_time": 0,
				"loss": 100,
				"test_date": "2018-12-07T21:59:58.312Z",
				"inactive_since": null
			}]`)
		case "/update/10":
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Update request body invalid, error: %v", err)
			}

			var req updateDeviceRequest
			err = json.Unmarshal(reqBody, &req)
			if err != nil {
				t.Errorf("Update request body failed to unmarshal, error: %v", err)
			}

			if req.Loss != 50.0 || req.AvgTime != 0.11 {
				t.Errorf("Update request invalid body, got: %v", req)
			}

			fmt.Fprintln(w, "some API response")
		default:
			t.Error("Unexpected request")
		}
	}))
	defer ts.Close()

	var client = apiClient{
		apiConfig: &schema.APIConfig{
			Name:   "test",
			Secret: "secret",
		},
	}
	client.apiConfig.URL = ts.URL
	client.apiConfig.Endpoints = schema.APIEndpoints{
		Authenticate: "/auth",
		GetDevices:   "/devices/%d",
		UpdateDevice: "/update/%d",
	}

	hosts, err := APILoadHosts(trueParser, client.apiConfig)
	if err != nil {
		t.Errorf("Fail to load hosts, got error: %v", err)
	}

	if hosts[0].IP != "192.168.12.150/24" {
		t.Errorf("Expected to read correct IP, got: %v", hosts)
	}

	err = APISavePingResult(schema.ProbeResult{Host: hosts[0], Loss: 50.0, AvgTime: 0.11}, client.apiConfig)
	if err != nil {
		t.Errorf("Fail to save ping result, got error: %v", err)
	}
}
