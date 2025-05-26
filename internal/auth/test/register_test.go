package test

import (
	"context"
	"encoding/json"
	"enuma-elish/internal/auth/repository"
	"enuma-elish/internal/auth/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	server := httptest.NewServer(testApi)
	defer server.Close()

	body := request.Register{
		Password: "12345678",
		Name:     "test",
		Email:    "galihwisnu18@gmail.com",
	}

	response := commonHttp.NewResponse()
	httpClient := commonHttp.NewHttpClient().
		SetRequestBody(body).
		SetUrl(server.URL + "/api/v1/auth/register").
		SetMethod("POST").
		SetJsonHeader().
		Do().
		UnmarshalResponse(&response)

	if httpClient.Error() != nil {
		t.Errorf("http response error: %v", httpClient.Error())
	}

	if httpClient.Status() != http.StatusOK {
		t.Fatalf("expected status 200 Created, got %d", httpClient.Status())
	}

	s, err := testInfra.Redis.Get(context.Background(), repository.VerifyEmailTokenKey+":"+body.Email).Result()
	if err != nil {
		t.Errorf("Failed to get userVerify")
	}

	u := repository.UserVerifyEmailToken{}
	err = json.Unmarshal([]byte(s), &u)
	if err != nil {
		t.Errorf("Failed to unmarshal userVerify")
	}

	verifyEmailBody := request.VerifyEmailRequest{
		Token: u.Token,
		Email: body.Email,
	}

	vRes := commonHttp.NewResponse()
	httpClient = commonHttp.NewHttpClient().
		SetRequestBody(verifyEmailBody).
		SetUrl(server.URL + "/api/v1/auth/register/verify-email").
		SetMethod("POST").
		SetJsonHeader().
		Do().
		UnmarshalResponse(vRes)

	if httpClient.Error() != nil {
		t.Errorf("http response error: %v", httpClient.Error())
	}

	if httpClient.Status() != http.StatusOK {
		t.Fatalf("expected status 200 Created, got %d", httpClient.Status())
	}

	_, err = testInfra.Postgres.Exec("DELETE FROM users WHERE email=$1", body.Email)
	if err != nil {
		t.Log(err.Error())
	}
}
