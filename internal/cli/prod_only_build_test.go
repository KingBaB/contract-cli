package cli

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"cn.qfei/contract-cli/internal/config"
)

// This policy must also pass with -tags dev: build tags cannot enable another environment.
func TestAllBuildsRejectNonProductionPresets(t *testing.T) {
	t.Parallel()
	for _, environment := range []string{"dev", "test", "blue"} {
		t.Run(environment, func(t *testing.T) {
			t.Parallel()
			store := config.NewStore(t.TempDir())
			calls := 0
			app := New(Options{Store: store, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{},
				LookupEnv: func(key string) (string, bool) { return "1", key == "CONTRACT_CLI_NO_UPDATE_CHECK" },
				HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
					calls++
					return nil, errors.New("unsupported environment reached the transport")
				})},
			})
			err := app.Run(context.Background(), []string{"config", "add", "--env", environment, "--name", "blocked"})
			if err == nil || !strings.Contains(err.Error(), "unsupported environment") || calls != 0 {
				t.Fatalf("build allowed non-production preset: environment=%s calls=%d error=%v", environment, calls, err)
			}
			if _, found, err := store.LookupProfile("blocked"); err != nil || found {
				t.Fatalf("unsupported profile was persisted: found=%t error=%v", found, err)
			}
		})
	}
}

func TestAllBuildsBlockDevelopmentOrigins(t *testing.T) {
	t.Parallel()
	for _, rawURL := range []string{
		"https://dev-open.qtech.cn/probe", "https://dev-myaccount.qtech.cn/probe",
		"http://dev-open.qtech.cn/probe", "http://dev-myaccount.qtech.cn/probe",
		"https://dev-open.qtech.cn./probe", "https://dev-myaccount.qtech.cn./probe",
		"https://dev-open.qtech.cn:8443/probe", "https://dev-myaccount.qtech.cn:8443/probe",
	} {
		t.Run(rawURL, func(t *testing.T) {
			t.Parallel()
			calls := 0
			app := New(Options{Store: config.NewStore(t.TempDir()), Stderr: &bytes.Buffer{},
				HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
					calls++
					return productionResponse(http.StatusOK, `{}`), nil
				})},
			})
			response, err := app.httpClient.Get(rawURL)
			if response != nil {
				response.Body.Close()
			}
			if err == nil || !strings.Contains(err.Error(), "blocks non-production host") || calls != 0 {
				t.Fatalf("development origin reached transport: calls=%d error=%v", calls, err)
			}
		})
	}
}

func TestAllBuildsBlockProductionRedirectToDevelopment(t *testing.T) {
	t.Parallel()
	for _, origin := range []string{"https://dev-open.qtech.cn", "https://dev-myaccount.qtech.cn"} {
		t.Run(origin, func(t *testing.T) {
			t.Parallel()
			calls := 0
			app := New(Options{Store: config.NewStore(t.TempDir()), Stderr: &bytes.Buffer{},
				HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
					calls++
					response := productionResponse(http.StatusTemporaryRedirect, "")
					response.Header.Set("Location", origin+"/finish")
					return response, nil
				})},
			})
			response, err := app.httpClient.Get(productionOpenPlatformOrigin + "/start")
			if response != nil {
				response.Body.Close()
			}
			if err == nil || !strings.Contains(err.Error(), "blocks non-production host") || calls != 1 {
				t.Fatalf("production redirect reached development origin: calls=%d error=%v", calls, err)
			}
		})
	}
}
