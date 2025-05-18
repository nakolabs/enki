package test

import (
	"enuma-elish/internal/ppdb/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSchoolSuccess(t *testing.T) {
	server := httptest.NewServer(testApi)
	defer server.Close()

	reqBody := request.CreateSchoolRequest{
		Name:  "test",
		Level: "elementary",
	}

	response := commonHttp.NewResponse()
	httpClient := commonHttp.NewHttpClient().
		SetUrl(server.URL + "/api/v1/school/create").
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
