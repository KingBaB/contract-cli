package cli

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"cn.qfei/contract-cli/internal/config"
	"cn.qfei/contract-cli/internal/credential"
	"cn.qfei/contract-cli/internal/openplatform"
)

func TestLoadDeviceAwareProfileRejectsNonProductionProfileAndEndpoints(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*config.Profile)
	}{
		{name: "dev environment", mutate: func(profile *config.Profile) { profile.Environment = "dev" }},
		{name: "dev open platform", mutate: func(profile *config.Profile) { profile.OpenPlatformBaseURL = "https://dev-open.qtech.cn" }},
		{name: "dev resource", mutate: func(profile *config.Profile) { profile.Resource = "https://dev-open.qtech.cn" }},
		{name: "dev app token", mutate: func(profile *config.Profile) {
			profile.AppTokenEndpoint = "https://dev-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal"
		}},
		{name: "dev protected metadata", mutate: func(profile *config.Profile) {
			profile.ProtectedResourceMetadataURL = "https://dev-open.qtech.cn/.well-known/oauth-protected-resource"
		}},
		{name: "dev authorization metadata", mutate: func(profile *config.Profile) {
			profile.AuthorizationServerMetadataURL = "https://dev-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract"
		}},
		{name: "dev authorization endpoint", mutate: func(profile *config.Profile) {
			profile.Identities.User.AuthorizationEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/authorize/contract"
		}},
		{name: "dev device endpoint", mutate: func(profile *config.Profile) {
			profile.Identities.User.DeviceAuthorizationEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/device-authorization/contract"
		}},
		{name: "dev token endpoint", mutate: func(profile *config.Profile) {
			profile.Identities.User.TokenEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract"
		}},
		{name: "dev revoke endpoint", mutate: func(profile *config.Profile) {
			profile.Identities.User.RevocationEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/revoke/contract"
		}},
		{name: "dev registration endpoint", mutate: func(profile *config.Profile) {
			profile.Identities.User.RegistrationEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/register/contract"
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			profile := validProductionProfile()
			test.mutate(&profile)
			if err := store.UpsertProfile(profile, true); err != nil {
				t.Fatal(err)
			}
			app := New(Options{Store: store})

			_, err := app.loadDeviceAwareProfile("contract")
			if err == nil || err.Error() != productionProfileErrorMessage("contract") {
				t.Fatalf("load profile error = %v", err)
			}
		})
	}
}

func TestAuthAndBusinessCommandsRejectDevProfileBeforeNetworkRequest(t *testing.T) {
	commands := []struct {
		name string
		run  func(*App) error
	}{
		{name: "legacy user login", run: func(app *App) error {
			return app.Run(context.Background(), []string{"auth", "login", "--profile", "contract", "--as", "user", "--no-open-browser"})
		}},
		{name: "app login", run: func(app *App) error {
			return app.Run(context.Background(), []string{"auth", "login", "--profile", "contract", "--as", "app", "--app-id", "app-id", "--app-secret", "app-secret"})
		}},
		{name: "device init", run: func(app *App) error {
			return app.Run(context.Background(), []string{"auth", "init", "--profile", "contract", "--output", "json"})
		}},
		{name: "auth status", run: func(app *App) error {
			return app.Run(context.Background(), []string{"auth", "status", "--profile", "contract", "--as", "user"})
		}},
		{name: "auth logout", run: func(app *App) error {
			return app.Run(context.Background(), []string{"auth", "logout", "--profile", "contract", "--as", "user"})
		}},
		{name: "identity switch", run: func(app *App) error {
			return app.Run(context.Background(), []string{"auth", "use", "--profile", "contract", "--as", "app"})
		}},
		{name: "app business request", run: func(app *App) error {
			_, _, err := app.openPlatformClientAndContext("contract", "app", "/open-apis/contract/v1/mcp/contracts/search", openplatform.IdentityPolicyAny)
			return err
		}},
		{name: "device business request", run: func(app *App) error {
			_, _, err := app.openPlatformClientAndContext("contract", "user", "/open-apis/contract/v1/mcp/contracts/search", openplatform.IdentityPolicyUserOnly)
			return err
		}},
	}

	for _, command := range commands {
		t.Run(command.name, func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			profile := validProductionProfile()
			profile.Environment = "dev"
			profile.OpenPlatformBaseURL = "https://dev-open.qtech.cn"
			profile.Resource = "https://dev-open.qtech.cn"
			if err := store.UpsertProfile(profile, true); err != nil {
				t.Fatal(err)
			}
			requestCount := 0
			app := New(Options{
				Store:           store,
				Secrets:         config.NewSecretsStore(t.TempDir()),
				CredentialStore: &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{}},
				HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
					requestCount++
					return nil, errors.New("network request must not run")
				})},
			})

			err := command.run(app)
			if err == nil || err.Error() != productionProfileErrorMessage("contract") {
				t.Fatalf("command error = %v", err)
			}
			if requestCount != 0 {
				t.Fatalf("request count = %d, want 0", requestCount)
			}
		})
	}
}

func TestDeviceProfileRecoveryRejectsDevSnapshotBeforeSavingProfile(t *testing.T) {
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		"contract": {DeviceProfile: &credential.DeviceProfile{
			Name: "contract", Environment: "dev", OpenPlatformBaseURL: "https://dev-open.qtech.cn",
			Resource: "https://dev-open.qtech.cn", BusinessType: "contract", ClientName: "contract-cli",
			DeviceClientID: "zscli_892efdadc11a3f53", DeviceScope: "contract:full contract-review:full",
			DeviceAuthorizationEndpoint: "https://dev-myaccount.qtech.cn/api/public/oauth/device-authorization/contract",
			TokenEndpoint:               "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract",
			RevocationEndpoint:          "https://dev-myaccount.qtech.cn/api/public/oauth/revoke/contract",
		}},
	}}
	store := config.NewStore(t.TempDir())
	app := New(Options{
		Store: store, CredentialStore: credentials,
		LookupEnv: func(name string) (string, bool) {
			if name == "SKILL_SESSION_WORKSPACE" {
				return t.TempDir(), true
			}
			return "", false
		},
	})

	_, err := app.loadDeviceAwareProfile("contract")
	if err == nil || err.Error() != productionProfileErrorMessage("contract") {
		t.Fatalf("profile recovery error = %v", err)
	}
	if _, found, lookupErr := store.LookupProfile("contract"); lookupErr != nil || found {
		t.Fatalf("dev snapshot must not be restored, found=%v err=%v", found, lookupErr)
	}
}

func TestDeviceProfileRecoveryRejectsDevPendingBeforeSavingProductionSnapshot(t *testing.T) {
	profile := validProductionProfile()
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		"contract": {
			DeviceProfile: snapshotDeviceProfile(profile),
			Pending: &credential.PendingTransaction{
				Status: credential.PendingStatusPending, DeviceCode: "secret-device",
				TokenEndpoint: "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract",
				ClientID:      "zscli_892efdadc11a3f53", ExpiresAt: time.Now().Add(time.Minute),
			},
		},
	}}
	store := config.NewStore(t.TempDir())
	app := New(Options{
		Store: store, CredentialStore: credentials,
		LookupEnv: func(name string) (string, bool) {
			if name == "SKILL_SESSION_WORKSPACE" {
				return t.TempDir(), true
			}
			return "", false
		},
	})

	_, err := app.loadDeviceAwareProfile("contract")
	if err == nil || err.Error() != productionProfileErrorMessage("contract") {
		t.Fatalf("profile recovery error = %v", err)
	}
	if _, found, lookupErr := store.LookupProfile("contract"); lookupErr != nil || found {
		t.Fatalf("snapshot with dev pending state must not be restored, found=%v err=%v", found, lookupErr)
	}
}

func TestDeviceCommandsRejectHistoricalDevStateBeforeNetworkRequest(t *testing.T) {
	tests := []struct {
		name       string
		credential credential.DeviceCredential
		args       []string
	}{
		{
			name: "pending token endpoint", args: []string{"auth", "complete", "--profile", "contract", "--output", "json"},
			credential: credential.DeviceCredential{Pending: &credential.PendingTransaction{
				Status: credential.PendingStatusPending, DeviceCode: "secret-device",
				TokenEndpoint: "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract",
				ClientID:      "zscli_892efdadc11a3f53", ExpiresAt: time.Now().Add(time.Minute),
			}},
		},
		{
			name: "pending verification uri", args: []string{"auth", "init", "--profile", "contract", "--output", "json"},
			credential: credential.DeviceCredential{Pending: &credential.PendingTransaction{
				Status: credential.PendingStatusPending, DeviceCode: "secret-device",
				VerificationURIComplete: "https://dev-myaccount.qtech.cn/device?user_code=secret-user-code",
				TokenEndpoint:           "https://myaccount.qfei.cn/api/public/oauth/token/contract",
				ClientID:                "zscli_892efdadc11a3f53", ExpiresAt: time.Now().Add(time.Minute),
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			if err := store.UpsertProfile(validProductionProfile(), true); err != nil {
				t.Fatal(err)
			}
			credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{"contract": test.credential}}
			requestCount := 0
			app := New(Options{
				Store: store, CredentialStore: credentials,
				HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
					requestCount++
					return nil, errors.New("network request must not run")
				})},
			})

			err := app.Run(context.Background(), test.args)
			if err == nil || err.Error() != productionProfileErrorMessage("contract") {
				t.Fatalf("device command error = %v", err)
			}
			if requestCount != 0 {
				t.Fatalf("request count = %d, want 0", requestCount)
			}
		})
	}
}

func TestBusinessCommandRejectsHistoricalDevDeviceStateWithProductionProfile(t *testing.T) {
	store := config.NewStore(t.TempDir())
	profile := validProductionProfile()
	profile.DefaultIdentity = config.IdentityApp
	profile.Identities.App = config.AppIdentity{AppID: "production-app", Token: &config.Token{AccessToken: "production-app-token"}}
	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatal(err)
	}
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		"contract": {Pending: &credential.PendingTransaction{
			Status: credential.PendingStatusPending, DeviceCode: "secret-device",
			TokenEndpoint: "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract",
			ClientID:      "zscli_892efdadc11a3f53", ExpiresAt: time.Now().Add(time.Minute),
		}},
	}}
	requestCount := 0
	app := New(Options{
		Store: store, CredentialStore: credentials,
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
			requestCount++
			return nil, errors.New("network request must not run")
		})},
	})

	_, _, err := app.openPlatformClientAndContext(
		"contract",
		"app",
		"/open-apis/contract/v1/mcp/contracts/search",
		openplatform.IdentityPolicyAny,
	)
	if err == nil || err.Error() != productionProfileErrorMessage("contract") {
		t.Fatalf("business command error = %v", err)
	}
	if requestCount != 0 {
		t.Fatalf("request count = %d, want 0", requestCount)
	}
}

func TestConfigAddProdResetsSameNameDevProfileAndCredentials(t *testing.T) {
	configDir := t.TempDir()
	store := config.NewStore(configDir)
	secrets := config.NewSecretsStore(configDir)
	oldProfile := validProductionProfile()
	oldProfile.Environment = "dev"
	oldProfile.OpenPlatformBaseURL = "https://dev-open.qtech.cn"
	oldProfile.Resource = "https://dev-open.qtech.cn"
	oldProfile.DefaultIdentity = config.IdentityApp
	oldProfile.Identities.User.Token = &config.Token{AccessToken: "old-user-token"}
	oldProfile.Identities.App = config.AppIdentity{AppID: "old-app", Token: &config.Token{AccessToken: "old-app-token"}}
	if err := store.UpsertProfile(oldProfile, true); err != nil {
		t.Fatal(err)
	}
	if err := secrets.Set(config.AppSecretKey("contract"), "old-app-secret"); err != nil {
		t.Fatal(err)
	}
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		"contract": {
			Token: &config.Token{AccessToken: "old-device-token", RefreshToken: "old-refresh-token"},
			Pending: &credential.PendingTransaction{
				DeviceCode: "old-device-code", TokenEndpoint: "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract",
				ClientID: "old-client", ExpiresAt: time.Now().Add(time.Minute),
			},
		},
	}}
	app := New(Options{
		Store: store, Secrets: secrets, CredentialStore: credentials,
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.URL.String() != "https://myaccount.qfei.cn/.well-known/oauth-authorization-server/contract" {
				t.Fatalf("unexpected metadata request: %s", request.URL.String())
			}
			return productionResponse(http.StatusOK, `{"issuer":"common-organization-v2","authorization_endpoint":"https://myaccount.qfei.cn/api/public/oauth/authorize/contract","device_authorization_endpoint":"https://myaccount.qfei.cn/api/public/oauth/device-authorization/contract","token_endpoint":"https://myaccount.qfei.cn/api/public/oauth/token/contract","revocation_endpoint":"https://myaccount.qfei.cn/api/public/oauth/revoke/contract","registration_endpoint":"https://myaccount.qfei.cn/api/public/oauth/register/contract"}`), nil
		})},
	})

	if err := app.Run(context.Background(), []string{"config", "add", "--env", "prod", "--name", "contract"}); err != nil {
		t.Fatalf("config add error = %v", err)
	}
	profile, err := store.GetProfile("contract")
	if err != nil {
		t.Fatal(err)
	}
	if profile.Environment != "prod" || profile.OpenPlatformBaseURL != "https://open.qfei.cn" || profile.Resource != "https://open.qfei.cn" {
		t.Fatalf("unexpected reset profile: %+v", profile)
	}
	if profile.DefaultIdentity != config.IdentityUser || profile.Identities.User.Token != nil || profile.Identities.App.AppID != "" || profile.Identities.App.Token != nil {
		t.Fatalf("old identities survived reset: %+v", profile.Identities)
	}
	if _, ok, err := secrets.Get(config.AppSecretKey("contract")); err != nil || ok {
		t.Fatalf("old app secret survived reset: ok=%v err=%v", ok, err)
	}
	if _, err := credentials.Load("contract"); !errors.Is(err, credential.ErrCredentialNotFound) {
		t.Fatalf("old device credential survived reset: %v", err)
	}
}

func TestConfigAddRejectsNonProductionMetadataOverrideBeforeNetworkRequest(t *testing.T) {
	requestCount := 0
	app := New(Options{
		Store: config.NewStore(t.TempDir()),
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
			requestCount++
			return nil, errors.New("network request must not run")
		})},
	})

	err := app.Run(context.Background(), []string{
		"config", "add", "--env", "prod", "--name", "contract",
		"--resource-metadata-url", "https://dev-open.qtech.cn/.well-known/oauth-protected-resource",
	})
	if err == nil || err.Error() != productionProfileErrorMessage("contract") {
		t.Fatalf("config add error = %v", err)
	}
	if requestCount != 0 {
		t.Fatalf("request count = %d, want 0", requestCount)
	}
}

func TestProductionHTTPClientBlocksKnownDevHosts(t *testing.T) {
	requestCount := 0
	app := New(Options{
		Store: config.NewStore(t.TempDir()),
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
			requestCount++
			return productionResponse(http.StatusOK, `{}`), nil
		})},
	})

	for _, rawURL := range []string{"https://dev-open.qtech.cn/open-apis/test", "https://dev-myaccount.qtech.cn./device"} {
		request, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = app.httpClient.Do(request)
		if err == nil || !strings.Contains(err.Error(), "production build blocks non-production host") {
			t.Fatalf("request %s error = %v", rawURL, err)
		}
	}
	if requestCount != 0 {
		t.Fatalf("request count = %d, want 0", requestCount)
	}
}

func validProductionProfile() config.Profile {
	return config.Profile{
		Name:                           "contract",
		Environment:                    "prod",
		OpenPlatformBaseURL:            "https://open.qfei.cn",
		AppTokenEndpoint:               "https://open.qfei.cn/open-apis/auth/v3/tenant_access_token/internal",
		ProtectedResourceMetadataURL:   "https://open.qfei.cn/.well-known/oauth-protected-resource",
		AuthorizationServerMetadataURL: "https://myaccount.qfei.cn/.well-known/oauth-authorization-server/contract",
		Resource:                       "https://open.qfei.cn",
		Scopes:                         []string{"cli:tools", "cli:resources"},
		BusinessType:                   "contract",
		ClientName:                     "contract-cli",
		DefaultIdentity:                config.IdentityUser,
		Identities: config.Identities{User: config.UserIdentity{
			AuthMode:                    config.UserAuthModeDevice,
			DeviceClientID:              "zscli_892efdadc11a3f53",
			DeviceScope:                 "contract:full contract-review:full",
			AuthorizationEndpoint:       "https://myaccount.qfei.cn/api/public/oauth/authorize/contract",
			DeviceAuthorizationEndpoint: "https://myaccount.qfei.cn/api/public/oauth/device-authorization/contract",
			TokenEndpoint:               "https://myaccount.qfei.cn/api/public/oauth/token/contract",
			RevocationEndpoint:          "https://myaccount.qfei.cn/api/public/oauth/revoke/contract",
			RegistrationEndpoint:        "https://myaccount.qfei.cn/api/public/oauth/register/contract",
		}},
	}
}

func productionProfileErrorMessage(profileName string) string {
	return `profile "` + profileName + `" belongs to a non-production environment and is not allowed in this production build; run ` +
		"`contract-cli config add --env prod --name contract` and authorize again"
}

type productionRoundTripFunc func(*http.Request) (*http.Response, error)

func (function productionRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func productionResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}
