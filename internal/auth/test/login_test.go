package test

import (
	"enuma-elish/internal/auth/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	server := httptest.NewServer(testApi)
	defer server.Close()

	reqBody := request.LoginRequest{
		Email:    "admin@gmail.com",
		Password: "12345678",
	}

	response := commonHttp.NewResponse()
	httpClient := commonHttp.NewHttpClient().
		SetUrl(server.URL + "/api/v1/auth/login").
		SetMethod(http.MethodPost).
		SetJsonHeader().
		SetRequestBody(&reqBody).
		Do().
		UnmarshalResponse(response)

	if httpClient.Error() != nil {
		t.Errorf("http response error: %v", httpClient.Error())
	}

	if httpClient.Status() != http.StatusOK {
		t.Fatalf("expected status 200 Created, got %d", httpClient.Status())
	}
}
