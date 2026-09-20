package cli

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"cn.qfei/contract-cli/internal/config"
)

func TestTestOnlyEnvironmentDefaultConfig(t *testing.T) {
	for _, explicit := range []bool{false, true} {
		t.Run(map[bool]string{false: "default", true: "explicit"}[explicit], func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			requests := 0
			app := New(Options{Store: store, Stdout: io.Discard, Stderr: io.Discard,
				HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(r *http.Request) (*http.Response, error) {
					requests++
					if r.URL.String() != "https://test-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract" {
						t.Errorf("unexpected discovery destination: %s", r.URL)
					}
					return productionResponse(200, `{"issuer":"https://test-myaccount.qtech.cn/contract","authorization_endpoint":"https://test-myaccount.qtech.cn/api/public/oauth/authorize/contract","token_endpoint":"https://test-myaccount.qtech.cn/api/public/oauth/token/contract","registration_endpoint":"https://test-myaccount.qtech.cn/api/public/oauth/register/contract"}`), nil
				})},
			})
			args := []string{"config", "add", "--name", "test-only"}
			if explicit {
				args = append(args, "--env", "test")
			}
			if err := app.Run(context.Background(), args); err != nil {
				t.Fatalf("test configuration rejected: %v", err)
			}
			p, err := store.GetProfile("test-only")
			if err != nil {
				t.Fatal(err)
			}
			if p.Environment != "test" || p.OpenPlatformBaseURL != "https://test-open.qtech.cn" || p.Resource != "https://test-open.qtech.cn" || requests != 1 || p.Identities.User.DeviceClientID != "" {
				t.Fatal("configuration must use test without a production Device client")
			}
		})
	}
}

func TestTestOnlyEnvironmentRejectsDisabledEntrypoints(t *testing.T) {
	for _, env := range []string{"prod", "dev"} {
		for _, args := range [][]string{
			{"config", "add", "--env", env},
			{"auth", "status", "--as", "user"},
			{"auth", "login", "--as", "app", "--app-id", "fixture", "--app-secret", "fixture-secret"},
			{"auth", "init"},
			{"contract", "get", "123", "--as", "user"},
			{"contract", "get", "123", "--as", "app"},
			{"review", "cooperation", "get", "--cooperation-id", "123"},
		} {
			t.Run(env+"/"+strings.Join(args[:2], "-"), func(t *testing.T) {
				p := config.Profile{Name: "blocked", Environment: env, OpenPlatformBaseURL: "https://open.qfei.cn", Resource: "https://open.qfei.cn"}
				if env == "dev" {
					p.OpenPlatformBaseURL = "https://dev-open.qtech.cn"
					p.Resource = p.OpenPlatformBaseURL
				}
				p.Identities.User.Token = &config.Token{AccessToken: "fixture-token", Expiry: time.Now().Add(time.Hour)}
				p.Identities.App = config.AppIdentity{AppID: "fixture", Token: &config.Token{AccessToken: "fixture-token", Expiry: time.Now().Add(time.Hour)}}
				store := config.NewStore(t.TempDir())
				if err := store.UpsertProfile(p, true); err != nil {
					t.Fatal(err)
				}
				requests := 0
				app := New(Options{Store: store, Secrets: config.NewSecretsStore(t.TempDir()), Stdout: io.Discard, Stderr: io.Discard,
					HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) { requests++; return productionResponse(200, `{}`), nil })},
				})
				err := app.Run(context.Background(), args)
				if err == nil || !strings.Contains(err.Error(), "test") || requests != 0 {
					t.Fatalf("disabled environment must fail before HTTP: err=%v requests=%d", err, requests)
				}
				saved, err := store.GetProfile("blocked")
				if err != nil || saved.Environment != env || saved.Identities.User.Token == nil {
					t.Fatal("blocked profile must be preserved")
				}
			})
		}
	}
}

func TestTestOnlyEnvironmentNetworkDestinations(t *testing.T) {
	for _, host := range []string{"test-open.qtech.cn", "test-myaccount.qtech.cn", "open.qfei.cn", "myaccount.qfei.cn", "dev-open.qtech.cn", "dev-myaccount.qtech.cn."} {
		t.Run(host, func(t *testing.T) {
			requests := 0
			app := New(Options{Store: config.NewStore(t.TempDir()), Stderr: io.Discard, HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) { requests++; return productionResponse(200, `{}`), nil })}})
			res, err := app.httpClient.Get("https://" + host + "/check")
			if res != nil {
				res.Body.Close()
			}
			if strings.HasPrefix(host, "test-") {
				if err != nil || requests != 1 {
					t.Fatalf("test must be allowed: %v requests=%d", err, requests)
				}
			} else if err == nil || requests != 0 {
				t.Fatalf("disabled destination must be blocked: %v requests=%d", err, requests)
			}
		})
	}
}
