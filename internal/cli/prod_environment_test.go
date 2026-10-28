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

func TestLoadDeviceAwareProfileRejectsTestProfile(t *testing.T) {
	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(validTestProfile(), true); err != nil {
		t.Fatal(err)
	}
	app := New(Options{Store: store})

	_, err := app.loadDeviceAwareProfile("contract")
	if err == nil || err.Error() != productionProfileErrorMessage("contract") {
		t.Fatalf("load profile error = %v", err)
	}
}

func TestLoadDeviceAwareProfileRejectsUnsupportedProfileAndMismatchedEndpoints(t *testing.T) {
	tests := []struct {
		name    string
		profile config.Profile
	}{
		{name: "prod environment", profile: func() config.Profile {
			profile := validDevProfile()
			profile.Environment = "prod"
			return profile
		}()},
		{name: "dev profile with test open platform", profile: func() config.Profile {
			profile := validDevProfile()
			profile.OpenPlatformBaseURL = "https://test-open.qtech.cn"
			return profile
		}()},
		{name: "dev profile with test auth endpoint", profile: func() config.Profile {
			profile := validDevProfile()
			profile.Identities.User.TokenEndpoint = "https://test-myaccount.qtech.cn/api/public/oauth/token/contract"
			return profile
		}()},
		{name: "test profile with dev resource", profile: func() config.Profile {
			profile := validTestProfile()
			profile.Resource = "https://dev-open.qtech.cn"
			return profile
		}()},
		{name: "test profile with dev auth endpoint", profile: func() config.Profile {
			profile := validTestProfile()
			profile.Identities.User.TokenEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract"
			return profile
		}()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			if err := store.UpsertProfile(test.profile, true); err != nil {
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

func TestOpenPlatformContextRejectsTestProfile(t *testing.T) {
	store := config.NewStore(t.TempDir())
	profile := validTestProfile()
	profile.DefaultIdentity = config.IdentityApp
	profile.Identities.App = config.AppIdentity{AppID: "test-app", Token: &config.Token{AccessToken: "test-app-token", Expiry: time.Now().Add(time.Hour)}}
	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatal(err)
	}
	app := New(Options{Store: store, Secrets: config.NewSecretsStore(t.TempDir())})

	_, _, err := app.openPlatformClientAndContext(
		"contract",
		"app",
		"/open-apis/contract/v1/mcp/contracts/search",
		openplatform.IdentityPolicyAny,
	)
	if err == nil || err.Error() != productionProfileErrorMessage("contract") {
		t.Fatalf("open platform context error = %v", err)
	}
}

func TestDeviceProfileRecoveryRejectsTestSnapshot(t *testing.T) {
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		"contract": {DeviceProfile: snapshotDeviceProfile(validTestProfile())},
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
		t.Fatalf("test snapshot was restored, found=%v err=%v", found, lookupErr)
	}
}

func TestDeviceProfileRecoveryRejectsProdPendingBeforeSavingDevSnapshot(t *testing.T) {
	profile := validDevProfile()
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		"contract": {
			DeviceProfile: snapshotDeviceProfile(profile),
			Pending: &credential.PendingTransaction{
				Status: credential.PendingStatusPending, DeviceCode: "secret-device",
				TokenEndpoint: "https://myaccount.qfei.cn/api/public/oauth/token/contract",
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
		t.Fatalf("snapshot with prod pending state must not be restored, found=%v err=%v", found, lookupErr)
	}
}

func TestDeviceCommandsRejectHistoricalProdStateBeforeNetworkRequest(t *testing.T) {
	tests := []struct {
		name       string
		credential credential.DeviceCredential
		args       []string
	}{
		{
			name: "pending token endpoint", args: []string{"auth", "complete", "--profile", "contract", "--output", "json"},
			credential: credential.DeviceCredential{Pending: &credential.PendingTransaction{
				Status: credential.PendingStatusPending, DeviceCode: "secret-device",
				TokenEndpoint: "https://myaccount.qfei.cn/api/public/oauth/token/contract",
				ClientID:      "zscli_892efdadc11a3f53", ExpiresAt: time.Now().Add(time.Minute),
			}},
		},
		{
			name: "pending verification uri", args: []string{"auth", "init", "--profile", "contract", "--output", "json"},
			credential: credential.DeviceCredential{Pending: &credential.PendingTransaction{
				Status: credential.PendingStatusPending, DeviceCode: "secret-device",
				VerificationURIComplete: "https://myaccount.qfei.cn/device?user_code=secret-user-code",
				TokenEndpoint:           "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract",
				ClientID:                "zscli_892efdadc11a3f53", ExpiresAt: time.Now().Add(time.Minute),
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			if err := store.UpsertProfile(validDevProfile(), true); err != nil {
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

func TestBusinessCommandRejectsHistoricalProdDeviceStateWithDevProfile(t *testing.T) {
	store := config.NewStore(t.TempDir())
	profile := validDevProfile()
	profile.DefaultIdentity = config.IdentityApp
	profile.Identities.App = config.AppIdentity{AppID: "production-app", Token: &config.Token{AccessToken: "production-app-token"}}
	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatal(err)
	}
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		"contract": {Pending: &credential.PendingTransaction{
			Status: credential.PendingStatusPending, DeviceCode: "secret-device",
			TokenEndpoint: "https://myaccount.qfei.cn/api/public/oauth/token/contract",
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

func TestConfigAddDevResetsSameNameProdProfileAndCredentials(t *testing.T) {
	configDir := t.TempDir()
	store := config.NewStore(configDir)
	secrets := config.NewSecretsStore(configDir)
	oldProfile := validProductionProfile()
	oldProfile.Environment = "prod"
	oldProfile.OpenPlatformBaseURL = "https://open.qfei.cn"
	oldProfile.Resource = "https://open.qfei.cn"
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
				DeviceCode: "old-device-code", TokenEndpoint: "https://myaccount.qfei.cn/api/public/oauth/token/contract",
				ClientID: "old-client", ExpiresAt: time.Now().Add(time.Minute),
			},
		},
	}}
	app := New(Options{
		Store: store, Secrets: secrets, CredentialStore: credentials,
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.URL.String() != "https://dev-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract" {
				t.Fatalf("unexpected metadata request: %s", request.URL.String())
			}
			return productionResponse(http.StatusOK, `{"issuer":"common-organization-v2","authorization_endpoint":"https://dev-myaccount.qtech.cn/api/public/oauth/authorize/contract","device_authorization_endpoint":"https://dev-myaccount.qtech.cn/api/public/oauth/device-authorization/contract","token_endpoint":"https://dev-myaccount.qtech.cn/api/public/oauth/token/contract","revocation_endpoint":"https://dev-myaccount.qtech.cn/api/public/oauth/revoke/contract","registration_endpoint":"https://dev-myaccount.qtech.cn/api/public/oauth/register/contract"}`), nil
		})},
	})

	if err := app.Run(context.Background(), []string{"config", "add", "--env", "dev", "--name", "contract"}); err != nil {
		t.Fatalf("config add error = %v", err)
	}
	profile, err := store.GetProfile("contract")
	if err != nil {
		t.Fatal(err)
	}
	if profile.Environment != "dev" || profile.OpenPlatformBaseURL != "https://dev-open.qtech.cn" || profile.Resource != "https://dev-open.qtech.cn" {
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
		"config", "add", "--env", "dev", "--name", "contract",
		"--resource-metadata-url", "https://open.qfei.cn/.well-known/oauth-protected-resource",
	})
	if err == nil || err.Error() != productionProfileErrorMessage("contract") {
		t.Fatalf("config add error = %v", err)
	}
	if requestCount != 0 {
		t.Fatalf("request count = %d, want 0", requestCount)
	}
}

func TestHTTPClientRejectsKnownTestHosts(t *testing.T) {
	requestCount := 0
	app := New(Options{
		Store: config.NewStore(t.TempDir()),
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
			requestCount++
			return productionResponse(http.StatusOK, `{}`), nil
		})},
	})

	for _, rawURL := range []string{"https://test-open.qtech.cn/open-apis/test", "https://test-myaccount.qtech.cn./device"} {
		request, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			t.Fatal(err)
		}
		response, err := app.httpClient.Do(request)
		if response != nil {
			_ = response.Body.Close()
		}
		if err == nil || !strings.Contains(err.Error(), "blocks unsupported environment host") {
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

func validTestProfile() config.Profile {
	profile := validProductionProfile()
	profile.Environment = "test"
	profile.OpenPlatformBaseURL = "https://test-open.qtech.cn"
	profile.AppTokenEndpoint = "https://test-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal"
	profile.ProtectedResourceMetadataURL = "https://test-open.qtech.cn/.well-known/oauth-protected-resource"
	profile.AuthorizationServerMetadataURL = "https://test-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract"
	profile.Resource = "https://test-open.qtech.cn"
	profile.Identities.User.AuthorizationEndpoint = "https://test-myaccount.qtech.cn/api/public/oauth/authorize/contract"
	profile.Identities.User.DeviceAuthorizationEndpoint = "https://test-myaccount.qtech.cn/api/public/oauth/device-authorization/contract"
	profile.Identities.User.TokenEndpoint = "https://test-myaccount.qtech.cn/api/public/oauth/token/contract"
	profile.Identities.User.RevocationEndpoint = "https://test-myaccount.qtech.cn/api/public/oauth/revoke/contract"
	profile.Identities.User.RegistrationEndpoint = "https://test-myaccount.qtech.cn/api/public/oauth/register/contract"
	return profile
}

func validDevProfile() config.Profile {
	profile := validProductionProfile()
	profile.Environment = "dev"
	profile.OpenPlatformBaseURL = "https://dev-open.qtech.cn"
	profile.AppTokenEndpoint = "https://dev-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal"
	profile.ProtectedResourceMetadataURL = "https://dev-open.qtech.cn/.well-known/oauth-protected-resource"
	profile.AuthorizationServerMetadataURL = "https://dev-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract"
	profile.Resource = "https://dev-open.qtech.cn"
	profile.Identities.User.AuthorizationEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/authorize/contract"
	profile.Identities.User.DeviceAuthorizationEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/device-authorization/contract"
	profile.Identities.User.TokenEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/token/contract"
	profile.Identities.User.RevocationEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/revoke/contract"
	profile.Identities.User.RegistrationEndpoint = "https://dev-myaccount.qtech.cn/api/public/oauth/register/contract"
	return profile
}

/*
productionProfileErrorMessage 返回测试断言用的环境端点错误文本。
入参 profileName（string）为配置名；返回 string 错误文本。
*/
func productionProfileErrorMessage(profileName string) string {
	return `profile "` + profileName + `" has unsupported or mismatched environment endpoints; run ` +
		"`contract-cli config add --env dev --name contract` and authorize again"
}

type productionRoundTripFunc func(*http.Request) (*http.Response, error)

func (function productionRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func productionResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}
