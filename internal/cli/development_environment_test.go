//go:build dev

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"cn.qfei/contract-cli/internal/config"
	"cn.qfei/contract-cli/internal/credential"
)

func developmentProfile(t *testing.T) config.Profile {
	t.Helper()
	return legacyNonProductionProfile(t, "dev", "https://dev-open.qtech.cn", "https://dev-myaccount.qtech.cn")
}

func TestDevelopmentBuildConfigAndTraceRequest(t *testing.T) {
	store := config.NewStore(t.TempDir())
	output := &bytes.Buffer{}
	requests := 0
	app := New(Options{
		Store: store, Stdout: output, Stderr: io.Discard,
		Logger:          slog.New(slog.NewTextHandler(io.Discard, nil)),
		CredentialStore: &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{}},
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(request *http.Request) (*http.Response, error) {
			requests++
			if request.URL.String() == "https://dev-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract" {
				user := developmentProfile(t).Identities.User
				metadata, err := json.Marshal(map[string]string{
					"issuer":                        "https://dev-myaccount.qtech.cn/contract",
					"authorization_endpoint":        user.AuthorizationEndpoint,
					"device_authorization_endpoint": user.DeviceAuthorizationEndpoint,
					"token_endpoint":                user.TokenEndpoint,
					"revocation_endpoint":           user.RevocationEndpoint,
					"registration_endpoint":         user.RegistrationEndpoint,
				})
				if err != nil {
					t.Fatal(err)
				}
				return productionResponse(http.StatusOK, string(metadata)), nil
			}
			if request.URL.String() != "https://dev-open.qtech.cn/open-apis/contract/v1/mcp/contracts/1?user_id_type=user_id" {
				t.Fatalf("unexpected request URL: %s", request.URL)
			}
			parts := regexp.MustCompile(`^00-([0-9a-f]{32})-([0-9a-f]{16})-01$`).FindStringSubmatch(request.Header.Get("traceparent"))
			if parts == nil || request.Header.Get("X-Log-Id") != parts[1] {
				t.Fatal("development request must preserve matching Trace headers")
			}
			if request.Header.Get("Authorization") != "Bearer dev-test-token" {
				t.Fatal("development request did not use its profile token")
			}
			return productionResponse(http.StatusOK, `{"code":0,"data":{"id":"1"}}`), nil
		})},
	})
	if err := app.Run(context.Background(), []string{"config", "add", "--env", "dev", "--name", "trace-dev"}); err != nil {
		t.Fatalf("configure development profile: %v", err)
	}
	profile, err := store.GetProfile("trace-dev")
	if err != nil {
		t.Fatal(err)
	}
	if profile.Environment != "dev" || profile.OpenPlatformBaseURL != "https://dev-open.qtech.cn" || profile.Resource != profile.OpenPlatformBaseURL {
		t.Fatalf("incorrect development profile environment: %s", profile.Environment)
	}
	profile.Identities.User.AuthMode = config.UserAuthModeAuthorizationCode
	profile.Identities.User.Token = &config.Token{AccessToken: "dev-test-token", Expiry: time.Now().Add(time.Hour)}
	if err := store.SaveProfile(profile); err != nil {
		t.Fatal(err)
	}
	if err := app.Run(context.Background(), []string{"contract", "get", "1", "--profile", "trace-dev", "--as", "user"}); err != nil {
		t.Fatalf("read development contract: %v", err)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want discovery and one read", requests)
	}
}

func TestDevelopmentBuildRejectsMixedEnvironmentEndpoints(t *testing.T) {
	for _, mutate := range []func(*config.Profile){
		func(profile *config.Profile) { profile.Resource = "https://open.qfei.cn" },
		func(profile *config.Profile) {
			profile.Identities.User.TokenEndpoint = "https://myaccount.qfei.cn/api/public/oauth/token/contract"
		},
		func(profile *config.Profile) { profile.Environment = "prod" },
	} {
		profile := developmentProfile(t)
		mutate(&profile)
		if err := validateProductionProfile(profile); err == nil {
			t.Fatal("cross-environment profile was accepted")
		}
	}
}

func TestDevelopmentBuildRestrictsDevelopmentTransportOrigins(t *testing.T) {
	for _, test := range []struct {
		url     string
		allowed bool
	}{
		{"https://dev-open.qtech.cn/open-apis/contract/v1/contracts", true},
		{"https://dev-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract", true},
		{"http://dev-open.qtech.cn/open-apis/test", false},
		{"https://dev-open.qtech.cn./open-apis/test", false},
		{"https://dev-open.qtech.cn:8443/open-apis/test", false},
		{"https://someone@dev-open.qtech.cn/open-apis/test", false},
		{"https://test-open.qtech.cn/open-apis/test", false},
	} {
		t.Run(test.url, func(t *testing.T) {
			requests := 0
			app := New(Options{Store: config.NewStore(t.TempDir()), HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
				requests++
				return productionResponse(http.StatusOK, `{}`), nil
			})}})
			response, err := app.httpClient.Get(test.url)
			if response != nil {
				response.Body.Close()
			}
			if test.allowed && (err != nil || requests != 1) {
				t.Fatalf("allowed development request: error=%v requests=%d", err, requests)
			}
			if !test.allowed && (err == nil || requests != 0) {
				t.Fatalf("invalid development request: error=%v requests=%d", err, requests)
			}
		})
	}
}

func TestDevelopmentBuildHelpAndCredentialIsolation(t *testing.T) {
	output := &bytes.Buffer{}
	app := New(Options{Store: config.NewStore(t.TempDir()), Stdout: output})
	if err := app.Run(context.Background(), []string{"config", "add", "--help"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "--env <prod|dev>") {
		t.Fatal("development build help must expose the dev preset")
	}
	profile := developmentProfile(t)
	stored := credential.DeviceCredential{DeviceProfile: snapshotDeviceProfile(profile)}
	if err := validateProductionDeviceCredential(profile.Name, stored, "dev"); err != nil {
		t.Fatalf("valid development snapshot rejected: %v", err)
	}
	if err := validateProductionDeviceCredential(profile.Name, stored, "prod"); err == nil {
		t.Fatal("development credentials must not be reused in production")
	}
}

func TestDevelopmentBuildPreservesSameEnvironmentLoginAndSupportsLogout(t *testing.T) {
	store := config.NewStore(t.TempDir())
	profile := developmentProfile(t)
	profile.Identities.User.AuthMode = config.UserAuthModeAuthorizationCode
	profile.Identities.User.Token = &config.Token{AccessToken: "dev-user-token", Expiry: time.Now().Add(time.Hour)}
	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatal(err)
	}
	requests := 0
	app := New(Options{
		Store: store, Stdout: io.Discard, Stderr: io.Discard,
		Logger:          slog.New(slog.NewTextHandler(io.Discard, nil)),
		CredentialStore: &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{}},
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(request *http.Request) (*http.Response, error) {
			requests++
			if request.URL.String() != "https://dev-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract" {
				t.Fatalf("unexpected request during local login lifecycle: %s", request.URL)
			}
			return productionResponse(http.StatusOK, `{"issuer":"https://dev-myaccount.qtech.cn/contract","authorization_endpoint":"https://dev-myaccount.qtech.cn/api/public/oauth/authorize/contract","token_endpoint":"https://dev-myaccount.qtech.cn/api/public/oauth/token/contract","registration_endpoint":"https://dev-myaccount.qtech.cn/api/public/oauth/register/contract"}`), nil
		})},
	})
	if err := app.Run(context.Background(), []string{"config", "add", "--env", "dev", "--name", profile.Name}); err != nil {
		t.Fatal(err)
	}
	saved, err := store.GetProfile(profile.Name)
	if err != nil || saved.Identities.User.Token == nil || saved.Identities.User.Token.AccessToken != "dev-user-token" {
		t.Fatal("reconfiguring the same dev environment must preserve its login")
	}
	if err := app.Run(context.Background(), []string{"auth", "logout", "--profile", profile.Name, "--as", "user"}); err != nil {
		t.Fatal(err)
	}
	saved, err = store.GetProfile(profile.Name)
	if err != nil || saved.Identities.User.Token != nil || saved.Identities.App.Token == nil {
		t.Fatal("legacy OAuth logout must clear only the selected user identity")
	}
	if requests != 1 {
		t.Fatalf("login lifecycle made %d requests, want only metadata discovery", requests)
	}
}

type unreadableDevelopmentCredentials struct {
	credential.Store
	err error
}

func (store unreadableDevelopmentCredentials) Load(string) (credential.DeviceCredential, error) {
	return credential.DeviceCredential{}, store.err
}

func TestDevelopmentBuildPreservesProfileWhenCredentialsCannotBeRead(t *testing.T) {
	store := config.NewStore(t.TempDir())
	profile := developmentProfile(t)
	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatal(err)
	}
	expected := errors.New("credential store unavailable")
	requests := 0
	app := New(Options{
		Store: store, Stdout: io.Discard, Stderr: io.Discard,
		Logger:          slog.New(slog.NewTextHandler(io.Discard, nil)),
		CredentialStore: unreadableDevelopmentCredentials{err: expected},
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
			requests++
			return productionResponse(http.StatusOK, `{}`), nil
		})},
	})
	err := app.Run(context.Background(), []string{"config", "add", "--env", "dev", "--name", profile.Name})
	if !errors.Is(err, expected) || requests != 0 {
		t.Fatalf("unreadable credentials must stop reconfiguration: error=%v requests=%d", err, requests)
	}
	saved, err := store.GetProfile(profile.Name)
	if err != nil || saved.Identities.App.Token == nil || saved.Identities.App.Token.AccessToken != profile.Identities.App.Token.AccessToken {
		t.Fatal("failed reconfiguration must preserve the existing profile")
	}
}
