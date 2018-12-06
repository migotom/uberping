package driver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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

// APILoadHosts load list of hosts using external API service, see README.md for more details.
func APILoadHosts(hostParser schema.HostParser, apiConfig *schema.APIConfig) ([]schema.Host, error) {
	// authorize
	client := http.Client{}

	authReq := authorizeRequest{Name: apiConfig.Name, Secret: apiConfig.Secret}
	authJSON, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}

	res, err := client.Post(apiConfig.URL+apiConfig.Endpoints.Auth, "application/json", bytes.NewBuffer(authJSON))

	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var authRes authorizeResponse
	err = json.Unmarshal(body, &authRes)
	if err != nil {
		return nil, err
	}
	apiConfig.AuthData = authRes

	_, ok2 := apiConfig.AuthData.(authorizeResponse)
	if !ok2 {
		fmt.Println("ups1")
	}

	// load hosts to test
	req, err := http.NewRequest("GET", apiConfig.URL+fmt.Sprintf(apiConfig.Endpoints.Devices, authRes.IDServer), nil)
	req.Header.Add("X-Auth-Token", authRes.Token)
	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apiDevices []schema.Host
	err = json.Unmarshal(body, &apiDevices)
	if err != nil {
		return nil, err
	}

	for i, device := range apiDevices {
		ip, err := hostParser(device.IP)
		if err != nil {
			return nil, err
		}
		apiDevices[i].IP = ip
	}

	return apiDevices, nil
}

// APISavePingResult ...
func APISavePingResult(result schema.PingResult, apiConfig *schema.APIConfig) error {
	authRes, ok := apiConfig.AuthData.(authorizeResponse)
	if !ok {
		fmt.Println("ups", apiConfig.AuthData)
		return errors.New("Can't load config")
	}
	fmt.Println(authRes)

	// TODO implement

	return nil
}
