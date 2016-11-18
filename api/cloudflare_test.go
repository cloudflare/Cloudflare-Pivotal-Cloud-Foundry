package api_test

import (
	"reflect"
	"testing"

	"github.com/cloudflare/Cloudflare-Pivotal-Cloud-Foundry/api"
)

func TestGetAuthHeaders(t *testing.T) {
	email := "my@email.com"
	key := "myKey"
	auth := api.AuthHeaders{XAuthEmail: email, XAuthKey: key}
	testApi := api.CloudflareAPI{Auth: auth}
	testApi.Auth = auth

	result := testApi.GetAuthHeaders()
	if !reflect.DeepEqual(result, auth) {
		t.Errorf("GetAuthHeaders failed")
	}
}

func TestSetAuthHeaders(t *testing.T) {
	email := "my@email.com"
	key := "myKey"
	auth := api.AuthHeaders{XAuthEmail: email, XAuthKey: key}
	testApi := api.CloudflareAPI{}

	testApi.SetAuthHeaders(auth)
	if !reflect.DeepEqual(testApi.Auth, auth) {
		t.Errorf("TestSetAuthHeaders failed")
	}
}
