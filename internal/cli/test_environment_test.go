package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"cn.qfei/contract-cli/internal/config"
	"cn.qfei/contract-cli/internal/credential"
)

func validTestEnvironmentProfile(t *testing.T) config.Profile {
	t.Helper()
	data, err := json.Marshal(validProductionProfile())
	if err != nil {
		t.Fatal(err)
	}
	text := strings.NewReplacer("https://open.qfei.cn", "https://test-open.qtech.cn", "https://myaccount.qfei.cn", "https://test-myaccount.qtech.cn").Replace(string(data))
	var profile config.Profile
	if err := json.Unmarshal([]byte(text), &profile); err != nil {
		t.Fatal(err)
	}
	profile.Environment = "test"
	profile.Identities.User.DeviceClientID = "test-device-client"
	return profile
}

func environmentMetadataResponse(environment string) *http.Response {
	origin := "https://myaccount.qfei.cn"
	if environment == "test" {
		origin = "https://test-myaccount.qtech.cn"
	}
	body := strings.ReplaceAll(`{"issuer":"ACCOUNT/contract","authorization_endpoint":"ACCOUNT/api/public/oauth/authorize/contract","device_authorization_endpoint":"ACCOUNT/api/public/oauth/device-authorization/contract","token_endpoint":"ACCOUNT/api/public/oauth/token/contract","revocation_endpoint":"ACCOUNT/api/public/oauth/revoke/contract","registration_endpoint":"ACCOUNT/api/public/oauth/register/contract"}`, "ACCOUNT", origin)
	return productionResponse(http.StatusOK, body)
}

func TestTestConfigAndExistingUserAndAppCommands(t *testing.T) {
	store := config.NewStore(t.TempDir())
	requests := 0
	app := New(Options{
		Store: store, Secrets: config.NewSecretsStore(t.TempDir()), Stdout: io.Discard, Stderr: io.Discard,
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(request *http.Request) (*http.Response, error) {
			requests++
			switch request.URL.String() {
			case "https://test-myaccount.qtech.cn/.well-known/oauth-authorization-server/contract":
				return environmentMetadataResponse("test"), nil
			case "https://test-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal":
				return productionResponse(http.StatusOK, `{"code":0,"tenant_access_token":"test-app-token","expire":7200}`), nil
			case "https://test-open.qtech.cn/open-apis/contract/v1/mcp/tasks?user_id_type=user_id":
				if got := request.Header.Get("Authorization"); got != "Bearer test-user-token" {
					t.Fatalf("unexpected user authorization header")
				}
				return productionResponse(http.StatusOK, `{"code":0,"success":true,"data":{"items":[]}}`), nil
			default:
				t.Fatalf("unexpected request: %s", request.URL)
				return nil, errors.New("unexpected request")
			}
		})},
	})
	if err := app.Run(context.Background(), []string{"config", "add", "--env", "test", "--name", "contract-test"}); err != nil {
		t.Fatal(err)
	}
	profile, err := store.GetProfile("contract-test")
	if err != nil {
		t.Fatal(err)
	}
	if profile.Environment != "test" || profile.Resource != "https://test-open.qtech.cn" || profile.Identities.User.TokenEndpoint != "https://test-myaccount.qtech.cn/api/public/oauth/token/contract" {
		t.Fatalf("unexpected test configuration")
	}
	if profile.Identities.User.DeviceClientID != "" {
		t.Fatal("test must not inherit the production Device client ID")
	}
	profile.Identities.User.Token = &config.Token{AccessToken: "test-user-token", Expiry: time.Now().Add(time.Hour)}
	if err := store.SaveProfile(profile); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"auth", "status", "--profile", "contract-test", "--as", "user"},
		{"contract", "approval", "task", "list", "--profile", "contract-test", "--as", "user"},
		{"auth", "login", "--profile", "contract-test", "--as", "app", "--app-id", "test-app", "--app-secret", "test-secret"},
	} {
		if err := app.Run(context.Background(), args); err != nil {
			t.Fatalf("%v: %v", args[:3], err)
		}
	}
	if requests != 3 {
		t.Fatalf("requests = %d, want 3", requests)
	}
}

func TestEnvironmentSwitchClearsAllAuthenticationAndSameEnvironmentPreservesIt(t *testing.T) {
	for _, target := range []string{"prod", "test"} {
		for _, same := range []bool{false, true} {
			t.Run(target+"/same="+map[bool]string{true: "yes", false: "no"}[same], func(t *testing.T) {
				profile := validProductionProfile()
				if (target == "test") == same {
					profile = validTestEnvironmentProfile(t)
				}
				profile.Identities.User.AuthMode = ""
				profile.Identities.User.Token = &config.Token{AccessToken: "user-token"}
				profile.Identities.App = config.AppIdentity{AppID: "app", Token: &config.Token{AccessToken: "app-token"}}
				store := config.NewStore(t.TempDir())
				secrets := config.NewSecretsStore(t.TempDir())
				if err := store.UpsertProfile(profile, true); err != nil {
					t.Fatal(err)
				}
				if err := secrets.Set(config.AppSecretKey(profile.Name), "app-secret"); err != nil {
					t.Fatal(err)
				}
				credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
					profile.Name: {Token: &config.Token{AccessToken: "device-token"}, DeviceProfile: snapshotDeviceProfile(profile)},
				}}
				app := New(Options{Store: store, Secrets: secrets, CredentialStore: credentials, Stdout: io.Discard, Stderr: io.Discard,
					HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(*http.Request) (*http.Response, error) {
						return environmentMetadataResponse(target), nil
					})},
				})
				if err := app.Run(context.Background(), []string{"config", "add", "--env", target, "--name", profile.Name}); err != nil {
					t.Fatal(err)
				}
				saved, err := store.GetProfile(profile.Name)
				if err != nil {
					t.Fatal(err)
				}
				if saved.Environment != target || (saved.Identities.User.Token != nil) != same || (saved.Identities.App.Token != nil) != same {
					t.Fatalf("authentication preservation must match environment")
				}
				if same && saved.Identities.User.DeviceClientID != profile.Identities.User.DeviceClientID {
					t.Fatal("same-environment Device client configuration must be preserved")
				}
				if _, exists, err := secrets.Get(config.AppSecretKey(profile.Name)); err != nil || exists != same {
					t.Fatalf("secret preservation: exists=%v err=%v", exists, err)
				}
				_, err = credentials.Load(profile.Name)
				if same && err != nil || !same && !errors.Is(err, credential.ErrCredentialNotFound) {
					t.Fatalf("credential preservation: %v", err)
				}
			})
		}
	}
}

func TestTestBrowserLoginPreservesOAuthPKCEFlow(t *testing.T) {
	profile := validTestEnvironmentProfile(t)
	profile.Identities.User.AuthMode = ""
	profile.Identities.User.RedirectURL = "http://127.0.0.1:8000/callback"
	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatal(err)
	}
	requests := 0
	app := New(Options{Store: store, Stdout: io.Discard, Stderr: io.Discard,
		HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(request *http.Request) (*http.Response, error) {
			requests++
			if request.URL.Host != "test-myaccount.qtech.cn" {
				t.Fatal("OAuth request reached another environment")
			}
			if request.URL.Path == "/api/public/oauth/register/contract" {
				return productionResponse(http.StatusOK, `{"client_id":"test-registered-client"}`), nil
			}
			if request.URL.Path != "/api/public/oauth/token/contract" {
				t.Fatalf("unexpected path %s", request.URL.Path)
			}
			if err := request.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if request.Form.Get("code_verifier") == "" || request.Form.Get("client_id") != "test-registered-client" || request.Form.Get("resource") != "https://test-open.qtech.cn" {
				t.Fatal("OAuth token exchange must keep PKCE and the test audience")
			}
			return productionResponse(http.StatusOK, `{"access_token":"test-user-token","token_type":"Bearer","expires_in":3600}`), nil
		})},
		OpenBrowser: func(rawURL string) error {
			parsed, err := url.Parse(rawURL)
			if err != nil {
				return err
			}
			if parsed.Host != "test-myaccount.qtech.cn" || parsed.Query().Get("code_challenge") == "" {
				t.Fatal("invalid test authorization URL")
			}
			return nil
		},
	})
	provider := app.userProvider.(userAuthProvider)
	provider.startCallbackServer = func(string) (authorizationCallback, error) {
		return fakeAuthorizationCallback{wait: func(context.Context, string) (string, error) { return "test-code", nil }}, nil
	}
	app.userProvider = provider
	if err := app.Run(context.Background(), []string{"auth", "login", "--profile", profile.Name, "--as", "user"}); err != nil {
		t.Fatal(err)
	}
	saved, err := store.GetProfile(profile.Name)
	if err != nil {
		t.Fatal(err)
	}
	if requests != 2 || saved.Identities.User.Token == nil || saved.Identities.User.Token.AccessToken != "test-user-token" || saved.Identities.User.AuthMode != config.UserAuthModeAuthorizationCode {
		t.Fatal("test browser login did not persist the user authorization")
	}
}

func TestTestMCPApprovalAndDownloadCommands(t *testing.T) {
	profile := validTestEnvironmentProfile(t)
	profile.Identities.User.AuthMode = config.UserAuthModeAuthorizationCode
	profile.Identities.User.Token = &config.Token{AccessToken: "test-user-token"}
	for _, command := range []struct {
		args         []string
		method, path string
	}{
		{[]string{"approval", "get", "process-1"}, "GET", "/process_instances/process-1"},
		{[]string{"approval", "comment", "list", "process-1"}, "GET", "/process_instances/process-1/comments"},
		{[]string{"approval", "comment", "create", "process-1", "--data", `{"content":"测试评论"}`}, "POST", "/process_instances/process-1/comments"},
		{[]string{"approval", "task", "list"}, "POST", "/tasks"},
		{[]string{"approval", "task", "approve", "task-1"}, "POST", "/tasks/task-1/approval"},
		{[]string{"approval", "task", "reject", "task-1", "--comment", "测试驳回"}, "POST", "/tasks/task-1/approval"},
		{[]string{"download-file", "file-1", "--contract", "contract-1", "--raw"}, "GET", "/contracts/contract-1/files/file-1/download"},
	} {
		t.Run(strings.Join(command.args[:2], " ")+command.path+strings.Join(command.args[2:], " "), func(t *testing.T) {
			store := config.NewStore(t.TempDir())
			if err := store.UpsertProfile(profile, true); err != nil {
				t.Fatal(err)
			}
			stdout := &bytes.Buffer{}
			requests := 0
			app := New(Options{Store: store, Stdout: stdout, Stderr: io.Discard,
				HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(request *http.Request) (*http.Response, error) {
					requests++
					if request.URL.Host == "files.example.test" {
						if request.Header.Get("Authorization") != "" {
							t.Fatal("signed download leaked user token")
						}
						response := productionResponse(http.StatusOK, "file")
						response.ContentLength = 4
						return response, nil
					}
					if request.URL.Host != "test-open.qtech.cn" || request.URL.Path != "/open-apis/contract/v1/mcp"+command.path || request.Method != command.method || request.Header.Get("Authorization") != "Bearer test-user-token" {
						t.Fatal("MCP request has the wrong environment, route, method or identity")
					}
					if command.args[0] == "download-file" {
						return productionResponse(http.StatusOK, `{"code":0,"data":{"file_size":4,"download_url":"https://files.example.test/contract?signature=secret"}}`), nil
					}
					return productionResponse(http.StatusOK, `{"code":0,"data":{}}`), nil
				})},
			})
			args := append([]string{"contract"}, command.args...)
			args = append(args, "--profile", profile.Name, "--as", "user")
			if err := app.Run(context.Background(), args); err != nil {
				t.Fatal(err)
			}
			if command.args[0] == "download-file" {
				if requests != 2 || stdout.String() != "file" {
					t.Fatal("test signed download did not complete")
				}
			} else if requests != 1 {
				t.Fatalf("requests = %d", requests)
			}
		})
	}
}

func TestTestProfileRejectsProductionDeviceSnapshot(t *testing.T) {
	profile := validTestEnvironmentProfile(t)
	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatal(err)
	}
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		profile.Name: {DeviceProfile: snapshotDeviceProfile(validProductionProfile())},
	}}
	app := New(Options{Store: store, CredentialStore: credentials, Stdout: io.Discard, Stderr: io.Discard})
	if _, err := app.loadDeviceAwareProfile(profile.Name); err == nil {
		t.Fatal("cross-environment snapshot must be rejected")
	}
}

func TestTestProfileRecoveryUsesMatchingEncryptedSnapshot(t *testing.T) {
	profile := validTestEnvironmentProfile(t)
	store := config.NewStore(t.TempDir())
	credentials := &deviceMemoryCredentialStore{values: map[string]credential.DeviceCredential{
		profile.Name: {DeviceProfile: snapshotDeviceProfile(profile)},
	}}
	workspace := t.TempDir()
	app := New(Options{Store: store, CredentialStore: credentials, Stdout: io.Discard, Stderr: io.Discard,
		LookupEnv: func(name string) (string, bool) { return workspace, name == "SKILL_SESSION_WORKSPACE" },
	})
	restored, err := app.loadDeviceAwareProfile(profile.Name)
	if err != nil {
		t.Fatal(err)
	}
	if restored.Environment != "test" || restored.Resource != profile.Resource {
		t.Fatal("restored wrong environment")
	}
}

func TestEnvironmentHTTPClientAllowsTestEnvironmentAndBlocksCrossEnvironmentRedirects(t *testing.T) {
	for _, start := range []string{"https://open.qfei.cn", "https://test-open.qtech.cn"} {
		for _, end := range []string{"https://open.qfei.cn", "https://test-open.qtech.cn"} {
			t.Run(start+"->"+end, func(t *testing.T) {
				requests := 0
				app := New(Options{Store: config.NewStore(t.TempDir()), Stderr: io.Discard,
					HTTPClient: &http.Client{Transport: productionRoundTripFunc(func(request *http.Request) (*http.Response, error) {
						requests++
						if request.URL.Path == "/start" {
							response := productionResponse(http.StatusTemporaryRedirect, "")
							response.Header.Set("Location", end+"/finish")
							return response, nil
						}
						return productionResponse(http.StatusOK, `{}`), nil
					})},
				})
				response, err := app.httpClient.Get(start + "/start")
				if response != nil {
					response.Body.Close()
				}
				if start == end {
					if err != nil || requests != 2 {
						t.Fatalf("same-environment request failed: %v, requests=%d", err, requests)
					}
				} else if err == nil || requests != 1 {
					t.Fatalf("cross-environment redirect must stop before second request: %v, requests=%d", err, requests)
				}
			})
		}
	}
}

func TestEnvironmentRedirectChecksEntireChain(t *testing.T) {
	app := New(Options{Store: config.NewStore(t.TempDir()), Stderr: io.Discard})
	first, _ := http.NewRequest(http.MethodGet, "https://files.example.test/start", nil)
	second, _ := http.NewRequest(http.MethodGet, "https://test-open.qtech.cn/continue", nil)
	target, _ := http.NewRequest(http.MethodGet, "https://open.qfei.cn/end", nil)
	if err := app.httpClient.CheckRedirect(target, []*http.Request{first, second}); err == nil {
		t.Fatal("redirect via an external host must not bypass environment isolation")
	}
}

func TestDeviceProfileRecoveryDoesNotSuggestSwitchingToProd(t *testing.T) {
	err := deviceProfileRecoveryError("contract-test", "snapshot missing")
	if strings.Contains(err.Error(), "--env prod") {
		t.Fatal("recovery must not suggest switching a test profile to prod")
	}
}

func TestTestAppRecoveryKeepsEnvironment(t *testing.T) {
	profile := validTestEnvironmentProfile(t)
	profile.AppTokenEndpoint = ""
	store := config.NewStore(t.TempDir())
	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatal(err)
	}
	app := New(Options{Store: store, Secrets: config.NewSecretsStore(t.TempDir()), Stderr: io.Discard})
	err := app.Run(context.Background(), []string{"auth", "login", "--profile", profile.Name, "--as", "app", "--app-id", "test-app", "--app-secret", "test-secret"})
	if err == nil || !strings.Contains(err.Error(), "--env test") {
		t.Fatalf("wrong recovery hint: %v", err)
	}
}
