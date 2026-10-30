package cli

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"cn.qfei/contract-cli/internal/config"
	"cn.qfei/contract-cli/internal/credential"
)

func fixedEnvironmentProfile(t *testing.T) config.Profile {
	t.Helper()
	profile := legacyNonProductionProfile(t, productionEnvironment, productionOpenPlatformOrigin, productionAccountOrigin)
	profile.Name = "contract"
	profile.DefaultIdentity = config.IdentityUser
	profile.Identities.App = config.AppIdentity{}
	return profile
}

func TestProfileEndpointsStayWithinFixedEnvironment(t *testing.T) {
	profile := fixedEnvironmentProfile(t)
	if err := validateProductionProfile(profile); err != nil {
		t.Fatalf("fixed environment profile rejected: %v", err)
	}
	for _, change := range []func(*config.Profile){
		func(p *config.Profile) { p.Environment = "blue" },
		func(p *config.Profile) { p.Resource = "https://other-api.example" },
		func(p *config.Profile) { p.AppTokenEndpoint = "https://other-api.example/token" },
		func(p *config.Profile) { p.Identities.User.TokenEndpoint = "https://other-account.example/token" },
		func(p *config.Profile) {
			p.OpenPlatformBaseURL = "http://" + strings.TrimPrefix(productionOpenPlatformOrigin, "https://")
		},
	} {
		changed := profile
		change(&changed)
		if err := validateProductionProfile(changed); err == nil {
			t.Fatal("accepted mismatched or insecure profile")
		}
	}
}

func TestConfigRejectsExternalOriginOverridesBeforeRequests(t *testing.T) {
	for _, item := range []struct{ flag, value string }{
		{"--open-platform-base-url", "https://other-api.example"},
		{"--authorization-server-metadata-url", "https://other-account.example/.well-known/oauth-authorization-server/contract"},
	} {
		t.Run(item.flag, func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			requests := 0
			app := New(Options{Store: store, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{},
				HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
					requests++
					return productionResponse(http.StatusOK, `{}`), nil
				})},
			})
			err := app.Run(context.Background(), []string{"config", "add", "--env", productionEnvironment, "--name", "blocked", item.flag, item.value})
			if err == nil || requests != 0 {
				t.Fatalf("external origin was not rejected before HTTP: requests=%d err=%v", requests, err)
			}
			if _, found, err := store.LookupProfile("blocked"); err != nil || found {
				t.Fatalf("rejected override changed configuration: found=%t err=%v", found, err)
			}
		})
	}
}

func TestConfigFixedAddressesAndDeviceQuery(t *testing.T) {
	store := config.NewStore(t.TempDir())
	workspace := t.TempDir()
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{}}
	output := &bytes.Buffer{}
	requests := 0
	app := New(Options{
		Store: store, Secrets: config.NewSecretsStore(t.TempDir()), CredentialStore: credentials,
		Stdout: output, Stderr: &bytes.Buffer{},
		LookupEnv: func(name string) (string, bool) { return workspace, name == "SKILL_SESSION_WORKSPACE" },
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			if req.URL.Scheme+"://"+req.URL.Host == productionAccountOrigin {
				switch req.URL.Path {
				case "/.well-known/oauth-authorization-server/contract":
					return productionResponse(http.StatusOK, strings.ReplaceAll(testAccountMetadata, "https://test-myaccount.qtech.cn", productionAccountOrigin)), nil
				case "/api/public/oauth/device-authorization/contract":
					return productionResponse(http.StatusOK, `{"device_code":"fixture-code","user_code":"fixture-user","verification_uri":"`+productionAccountOrigin+`/device","verification_uri_complete":"`+productionAccountOrigin+`/device?user_code=fixture-user","expires_in":600}`), nil
				case "/api/public/oauth/token/contract":
					return productionResponse(http.StatusOK, `{"access_token":"fixture-token","refresh_token":"fixture-refresh","token_type":"Bearer","expires_in":3600}`), nil
				}
			}
			if req.URL.Scheme+"://"+req.URL.Host != productionOpenPlatformOrigin || req.Header.Get("Authorization") != "Bearer fixture-token" {
				t.Fatalf("request escaped configured target: %s", req.URL.Host)
			}
			return productionResponse(http.StatusOK, `{"code":0,"data":{"items":[]}}`), nil
		})},
	})
	commands := [][]string{
		{"config", "add", "--env", productionEnvironment, "--name", "contract", "--open-platform-base-url", productionOpenPlatformOrigin, "--authorization-server-metadata-url", productionAccountOrigin + "/.well-known/oauth-authorization-server/contract"},
		{"auth", "init", "--profile", "contract"},
		{"auth", "complete", "--profile", "contract"},
		{"contract", "search", "--profile", "contract", "--data", `{"page_size":1}`},
		{"api", "call", "--profile", "contract", "--channel", "web", "GET", "/open-web/web/clm/contract/info"},
	}
	for _, args := range commands {
		output.Reset()
		if err := app.Run(context.Background(), args); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		if args[0] == "auth" && args[1] == "complete" && !strings.Contains(output.String(), `"status":"succeeded"`) {
			t.Fatalf("authorization did not complete: %s", output.String())
		}
		if strings.Contains(output.String(), "fixture-token") || strings.Contains(output.String(), "fixture-code") {
			t.Fatal("credential leaked to command output")
		}
	}
	if requests != 5 {
		t.Fatalf("requests=%d, want metadata, init, complete and two queries", requests)
	}
	profile, err := store.GetProfile("contract")
	if err != nil || profile.Environment != productionEnvironment || profile.Resource != productionOpenPlatformOrigin {
		t.Fatalf("configured target mismatch: %v", err)
	}
}

func TestDeviceSnapshotRejectsAddressChangeWithSameEnvironmentLabel(t *testing.T) {
	old := fixedEnvironmentProfile(t)
	changed := fixedEnvironmentProfile(t)
	changed.Identities.User.TokenEndpoint = productionAccountOrigin + "/changed-token"
	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(changed, true); err != nil {
		t.Fatal(err)
	}
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		"contract": {DeviceProfile: snapshotDeviceProfile(old), Token: &config.Token{AccessToken: "fixture-token", Expiry: time.Now().Add(time.Hour)}},
	}}
	app := New(Options{Store: store, CredentialStore: credentials})
	if _, err := app.loadDeviceAwareProfile("contract"); err == nil {
		t.Fatal("accepted credential from previous target")
	}
}

func TestProfileHTTPRedirectsStayOnOriginalOrigin(t *testing.T) {
	for _, target := range []string{"https://other-api.example/finish", productionAccountOrigin + "/finish", "http://" + strings.TrimPrefix(productionOpenPlatformOrigin, "https://") + "/finish"} {
		t.Run(target, func(t *testing.T) {
			requests := 0
			app := New(Options{Store: config.NewStore(t.TempDir()), HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
				requests++
				response := productionResponse(http.StatusTemporaryRedirect, "")
				response.Header.Set("Location", target)
				return response, nil
			})}})
			response, err := app.httpClient.Get(productionOpenPlatformOrigin + "/start")
			if response != nil {
				_ = response.Body.Close()
			}
			if err == nil || requests != 1 {
				t.Fatalf("cross-origin redirect err=%v requests=%d", err, requests)
			}
		})
	}
}

const testAccountMetadata = `{"issuer":"contract","authorization_endpoint":"https://test-myaccount.qtech.cn/api/public/oauth/authorize/contract","device_authorization_endpoint":"https://test-myaccount.qtech.cn/api/public/oauth/device-authorization/contract","token_endpoint":"https://test-myaccount.qtech.cn/api/public/oauth/token/contract","revocation_endpoint":"https://test-myaccount.qtech.cn/api/public/oauth/revoke/contract","registration_endpoint":"https://test-myaccount.qtech.cn/api/public/oauth/register/contract"}`
