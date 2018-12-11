package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/migotom/uberping/internal/schema"
)

type authorizeRequest struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type authorizeResponse struct {
	Token    string `json:"token"`
	IDServer int    `json:"id_server"`
}

type updateDeviceRequest struct {
	Loss    int     `json:"loss"`
	AvgTime float64 `json:"average_time"`
}

type apiClient struct {
	httpClient http.Client
	authData   authorizeResponse
	apiConfig  *schema.APIConfig
}

// authorize obtains one-time token
func (c *apiClient) authorize() error {
	authReq := authorizeRequest{Name: c.apiConfig.Name, Secret: c.apiConfig.Secret}

	authJSON, err := json.Marshal(authReq)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Post(c.apiConfig.URL+c.apiConfig.Endpoints.Authenticate, "application/json", bytes.NewBuffer(authJSON))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf(res.Status)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &c.authData)
	if err != nil {
		return fmt.Errorf("Invalid JSON response: %v", err)
	}

	return nil
}

// do request to external API using auth token
func (c *apiClient) request(method, url string, requestBody []byte) ([]byte, error) {
	var lastError string

	for retries := 0; retries < 3; retries++ {

		req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, err
		}

		req.Header.Add("X-Auth-Token", c.authData.Token)
		req.Header.Set("Content-Type", "application/json")

		res, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode == 401 {
			// reauthorize and retry with new token
			c.authorize()
			continue
		} else if res.StatusCode != 200 {
			// note last error and retry
			lastError = res.Status
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		defer res.Body.Close()

		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return responseBody, nil
	}
	return nil, fmt.Errorf("request retry limit exceeded, last error: %v", lastError)
}

func getAPIClient(apiConfig *schema.APIConfig) *apiClient {
	client, ok := apiConfig.Client.(*apiClient)
	if !ok {
		client = &apiClient{}
		client.apiConfig = apiConfig
		apiConfig.Client = client
	}
	return client
}

// APILoadHosts load list of hosts using external API service, see README.md for more details.
func APILoadHosts(hostParser schema.HostParser, apiConfig *schema.APIConfig) ([]schema.Host, error) {
	client := getAPIClient(apiConfig)

	if err := client.authorize(); err != nil {
		return nil, err
	}

	body, err := client.request("GET", apiConfig.URL+fmt.Sprintf(apiConfig.Endpoints.GetDevices, client.authData.IDServer), nil)
	if err != nil {
		return nil, err
	}

	var apiDevices []schema.Host
	err = json.Unmarshal(body, &apiDevices)
	if err != nil {
		return nil, err
	}

	for i, device := range apiDevices {
		ip, port, err := hostParser(device.IP)
		if err != nil {
			return nil, err
		}
		apiDevices[i].IP = ip
		apiDevices[i].Port = port
	}

	return apiDevices, nil
}

// APISavePingResult save probe results using external API
func APISavePingResult(result schema.ProbeResult, apiConfig *schema.APIConfig) error {
	client := getAPIClient(apiConfig)

	// consider this, should skip update if host doesn't have ID, or should search in API using IP?
	if result.Host.ID == 0 {
		return nil
	}

	apiDevResult := updateDeviceRequest{Loss: int(result.Loss), AvgTime: result.AvgTime}

	apiDevResultJSON, err := json.Marshal(apiDevResult)
	if err != nil {
		return err
	}

	_, err = client.request("POST", apiConfig.URL+fmt.Sprintf(apiConfig.Endpoints.UpdateDevice, result.Host.ID), apiDevResultJSON)
	if err != nil {
		return err
	}

	return nil
}
