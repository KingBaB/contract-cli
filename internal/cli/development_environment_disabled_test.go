//go:build !dev

package cli

import (
	"net/http"
	"testing"

	"cn.qfei/contract-cli/internal/config"
)

func TestDefaultBuildStillRejectsDevelopment(t *testing.T) {
	if _, err := resolveEnvironment("dev"); err == nil {
		t.Fatal("default builds must reject dev")
	}
	requests := 0
	app := New(Options{Store: config.NewStore(t.TempDir()), HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
		requests++
		return productionResponse(http.StatusOK, `{}`), nil
	})}})
	response, err := app.httpClient.Get("https://dev-open.qtech.cn/open-apis/test")
	if response != nil {
		response.Body.Close()
	}
	if err == nil || requests != 0 {
		t.Fatalf("default development guard: error=%v requests=%d", err, requests)
	}
}
