package driver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/migotom/uberping/internal/schema"
)

var client = apiClient{
	apiConfig: &schema.APIConfig{
		Name:   "test",
		Secret: "secret",
	},
}

func TestAuthorizeInvalidCredentials(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()
	client.apiConfig.URL = ts.URL

	if err := client.authorize(); !strings.HasPrefix(err.Error(), "401 Unauthorized") {
		t.Errorf("Expecting 401 Unauthorized, got: %v", err)
	}
}

func TestAuthorizeInvalidAPIResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "some invalid API response")
	}))
	defer ts.Close()
	client.apiConfig.URL = ts.URL

	if err := client.authorize(); !strings.HasPrefix(err.Error(), "Invalid JSON response") {
		t.Errorf("Expecting JSON parsing failure, got: %v", err)
	}
}

func TestAuthorizeValidAPIResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"token": "some_token",
			"id_server": 1
		}`)
	}))
	defer ts.Close()
	client.apiConfig.URL = ts.URL

	if err := client.authorize(); err != nil {
		t.Errorf("Expecting correct response, got: %v", err)
	}
	if client.authData.Token != "some_token" {
		t.Errorf("Expecting token, got: %s", client.authData.Token)
	}
}

func TestRequestWithInvalidToken(t *testing.T) {
	tries := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tries++
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()
	client.apiConfig.URL = ts.URL

	if _, err := client.request("GET", ts.URL, nil); !strings.HasPrefix(err.Error(), "request retry limit exceeded") {
		t.Errorf("Expecting request retry limit exceeded, got: %v", err)
	}

	// 3 tries, each one try also fetch auth token
	if tries != 6 {
		t.Errorf("Expecting 3 tries, got: %v", tries)
	}
}

func TestRequestWithValidToken(t *testing.T) {
	tries := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tries++
		fmt.Fprintln(w, "some API response")
	}))
	defer ts.Close()
	client.apiConfig.URL = ts.URL

	body, err := client.request("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Expecting valid response, got: %v", err)
	}

	if tries != 1 {
		t.Errorf("Expecting 1 try, got: %v", tries)
	}

	if string(body) != "some API response\n" {
		t.Errorf("Expecting correct response, got: %v", string(body))
	}
}
