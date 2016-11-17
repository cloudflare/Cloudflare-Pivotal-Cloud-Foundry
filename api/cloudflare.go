package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const CLOUDFLARE_CLIENT_API_ENDPOINT = "https://api.cloudflare.com/client/v4/"
const CLOUDFLARE_CLIENT_API_ZONES = CLOUDFLARE_CLIENT_API_ENDPOINT + "zones/"
const X_AUTH_EMAIL_HEADER = "X-Auth-Email"
const X_AUTH_KEY_HEADER = "X-Auth-Key"

type CloudflareAPIInterface interface {
	AddZone(domain string) ([]byte, error)
	DeleteZone(zoneId string) error
	SetAuthHeaders(authHeaders AuthHeaders)
}

type CloudflareAPI struct {
	Auth AuthHeaders
}

type AuthHeaders struct {
	XAuthEmail string `json:"x-auth-email"`
	XAuthKey   string `json:"x-auth-key"`
}

func (api CloudflareAPI) GetAuthHeaders() AuthHeaders {
	return api.Auth
}

func (api *CloudflareAPI) SetAuthHeaders(authHeaders AuthHeaders) {
	api.Auth = authHeaders
}

func (api CloudflareAPI) AddZone(domain string) ([]byte, error) {
	var jsonBody = []byte(`{"name" : "` + domain + `"}`)
	request, err := http.NewRequest("POST", CLOUDFLARE_CLIENT_API_ZONES, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	var authHeaders = api.GetAuthHeaders()
	request.Header.Add(X_AUTH_EMAIL_HEADER, authHeaders.XAuthEmail)
	request.Header.Add(X_AUTH_KEY_HEADER, authHeaders.XAuthKey)

	var client = http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (api CloudflareAPI) DeleteZone(zoneId string) error {
	request, err := http.NewRequest("DELETE", CLOUDFLARE_CLIENT_API_ZONES+zoneId, nil)
	if err != nil {
		return err
	}

	var authHeaders = api.GetAuthHeaders()
	request.Header.Add(X_AUTH_EMAIL_HEADER, authHeaders.XAuthEmail)
	request.Header.Add(X_AUTH_KEY_HEADER, authHeaders.XAuthKey)

	var client = http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
