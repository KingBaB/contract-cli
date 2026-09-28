package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cn.qfei/contract-cli/internal/config"
)

func TestStoreRoundTrip(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	expiry := time.Date(2026, 4, 8, 18, 0, 0, 0, time.UTC)
	profile := config.Profile{
		Name:                "contract",
		Environment:         "dev",
		OpenPlatformBaseURL: "https://dev-open.qtech.cn",
		AppTokenEndpoint:    "https://dev-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal",
		Resource:            "http://example.com/mcp-servers",
		Scopes:              []string{"mcp:tools", "mcp:resources"},
		BusinessType:        "contract",
		ClientName:          "contract-cli",
		DefaultIdentity:     config.IdentityApp,
		Identities: config.Identities{
			User: config.UserIdentity{
				ClientID:              "zsdcli_test",
				TokenEndpoint:         "http://example.com/token",
				AuthorizationEndpoint: "http://example.com/authorize",
				RegistrationEndpoint:  "http://example.com/register",
				RedirectURL:           "http://127.0.0.1:8000/callback",
				Token: &config.Token{
					AccessToken: "token-value",
					TokenType:   "Bearer",
					Scope:       "mcp:tools mcp:resources",
					Expiry:      expiry,
				},
			},
			App: config.AppIdentity{
				AuthMode:     config.AppAuthModeAppCredentials,
				AppID:        "cli_app_123",
				SecretRef:    "contract.app.app_secret",
				ConfiguredAt: expiry,
				Token: &config.Token{
					AccessToken: "app-token",
					TokenType:   "Bearer",
				},
			},
		},
	}

	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}

	got, err := store.GetProfile("contract")
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}

	if got.Identities.User.ClientID != profile.Identities.User.ClientID {
		t.Fatalf("client_id mismatch: got %q want %q", got.Identities.User.ClientID, profile.Identities.User.ClientID)
	}
	if got.Identities.User.Token == nil || !got.Identities.User.Token.Expiry.Equal(expiry) {
		t.Fatalf("user token expiry mismatch: got %+v want %v", got.Identities.User.Token, expiry)
	}
	if got.Identities.App.Token == nil || got.Identities.App.Token.AccessToken != "app-token" {
		t.Fatalf("app token mismatch: got %+v", got.Identities.App.Token)
	}
	if got.AppTokenEndpoint != profile.AppTokenEndpoint {
		t.Fatalf("app token endpoint mismatch: got %q want %q", got.AppTokenEndpoint, profile.AppTokenEndpoint)
	}
	if got.OpenPlatformBaseURL != profile.OpenPlatformBaseURL {
		t.Fatalf("open platform base url mismatch: got %q want %q", got.OpenPlatformBaseURL, profile.OpenPlatformBaseURL)
	}
	if got.DefaultIdentity != config.IdentityApp {
		t.Fatalf("default identity mismatch: got %q want %q", got.DefaultIdentity, config.IdentityApp)
	}
	configContent, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatalf("ReadFile(config) error = %v", err)
	}
	if strings.Contains(string(configContent), "\"server_url\"") {
		t.Fatalf("config should not contain removed server_url field: %s", string(configContent))
	}
	if strings.Contains(string(configContent), "\"bot\"") || strings.Contains(string(configContent), "bot_token_endpoint") {
		t.Fatalf("config should use app field names, got: %s", string(configContent))
	}
	if store.Path() != filepath.Join(storeDir(t, store), "config.json") {
		t.Fatalf("unexpected config path: %s", store.Path())
	}
}

func TestParseIdentityKindAcceptsAppAndLegacyBotAlias(t *testing.T) {
	t.Parallel()

	identity, err := config.ParseIdentityKind("app")
	if err != nil {
		t.Fatalf("ParseIdentityKind(app) error = %v", err)
	}
	if identity != config.IdentityApp {
		t.Fatalf("identity = %q, want %q", identity, config.IdentityApp)
	}

	identity, err = config.ParseIdentityKind("bot")
	if err != nil {
		t.Fatalf("ParseIdentityKind(bot) error = %v", err)
	}
	if identity != config.IdentityApp {
		t.Fatalf("legacy bot identity = %q, want %q", identity, config.IdentityApp)
	}
}

func TestClearIdentityToken(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	profile := config.Profile{
		Name:        "contract",
		Environment: "dev",
		Identities: config.Identities{
			User: config.UserIdentity{
				Token: &config.Token{
					AccessToken: "user-token",
				},
			},
			App: config.AppIdentity{
				Token: &config.Token{
					AccessToken: "app-token",
				},
			},
		},
	}

	if err := store.UpsertProfile(profile, true); err != nil {
		t.Fatalf("UpsertProfile() error = %v", err)
	}
	if err := store.ClearToken("contract", config.IdentityApp); err != nil {
		t.Fatalf("ClearToken() error = %v", err)
	}

	got, err := store.GetProfile("contract")
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if got.Identities.App.Token != nil {
		t.Fatalf("expected app token to be cleared, got %+v", got.Identities.App.Token)
	}
	if got.Identities.User.Token == nil || got.Identities.User.Token.AccessToken != "user-token" {
		t.Fatalf("expected user token to stay intact, got %+v", got.Identities.User.Token)
	}
}

func TestLoadLegacyBotProfileMigratesToAppAndSavesNewShape(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	legacyConfig := `{
  "current_profile": "contract",
  "profiles": {
    "contract": {
      "name": "contract",
      "environment": "dev",
      "open_platform_base_url": "https://dev-open.qtech.cn",
      "bot_token_endpoint": "https://dev-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal",
      "default_identity": "bot",
      "identities": {
        "bot": {
          "auth_mode": "app_credentials",
          "app_id": "legacy-app-id",
          "secret_ref": "contract.bot.app_secret",
          "token": {
            "access_token": "legacy-app-token",
            "token_type": "Bearer"
          }
        }
      }
    }
  }
}`

	if err := os.WriteFile(store.Path(), []byte(legacyConfig), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := store.GetProfile("contract")
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if got.DefaultIdentity != config.IdentityApp {
		t.Fatalf("default identity = %q, want %q", got.DefaultIdentity, config.IdentityApp)
	}
	if got.AppTokenEndpoint != "https://dev-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal" {
		t.Fatalf("app token endpoint = %q", got.AppTokenEndpoint)
	}
	if got.Identities.App.AppID != "legacy-app-id" {
		t.Fatalf("app id = %q", got.Identities.App.AppID)
	}
	if got.Identities.App.Token == nil || got.Identities.App.Token.AccessToken != "legacy-app-token" {
		t.Fatalf("app token = %+v", got.Identities.App.Token)
	}

	if err := store.SaveProfile(got); err != nil {
		t.Fatalf("SaveProfile() error = %v", err)
	}
	content, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatalf("ReadFile(config) error = %v", err)
	}
	text := string(content)
	if strings.Contains(text, `"bot"`) || strings.Contains(text, "bot_token_endpoint") {
		t.Fatalf("saved config should not keep legacy bot fields: %s", text)
	}
	if !strings.Contains(text, `"app"`) || !strings.Contains(text, "app_token_endpoint") {
		t.Fatalf("saved config should use app fields: %s", text)
	}
}

func TestLoadProfilePrefersAppFieldsOverLegacyBotFields(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	rawConfig := `{
  "current_profile": "contract",
  "profiles": {
    "contract": {
      "name": "contract",
      "environment": "dev",
      "open_platform_base_url": "https://dev-open.qtech.cn",
      "app_token_endpoint": "https://dev-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal",
      "bot_token_endpoint": "https://legacy-open.example/open-apis/auth/v3/tenant_access_token/internal",
      "default_identity": "app",
      "identities": {
        "app": {
          "auth_mode": "app_credentials",
          "app_id": "new-app-id",
          "secret_ref": "contract.app.app_secret",
          "token": {
            "access_token": "new-app-token",
            "token_type": "Bearer"
          }
        },
        "bot": {
          "auth_mode": "app_credentials",
          "app_id": "legacy-app-id",
          "secret_ref": "contract.bot.app_secret",
          "token": {
            "access_token": "legacy-app-token",
            "token_type": "Bearer"
          }
        }
      }
    }
  }
}`

	if err := os.WriteFile(store.Path(), []byte(rawConfig), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := store.GetProfile("contract")
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if got.AppTokenEndpoint != "https://dev-open.qtech.cn/open-apis/auth/v3/tenant_access_token/internal" {
		t.Fatalf("app token endpoint = %q", got.AppTokenEndpoint)
	}
	if got.Identities.App.AppID != "new-app-id" {
		t.Fatalf("app id = %q", got.Identities.App.AppID)
	}
	if got.Identities.App.Token == nil || got.Identities.App.Token.AccessToken != "new-app-token" {
		t.Fatalf("app token = %+v", got.Identities.App.Token)
	}
}

func TestLoadLegacyProfileMigration(t *testing.T) {
	t.Parallel()

	store := config.NewStore(t.TempDir())
	legacyConfig := `{
  "current_profile": "contract",
  "profiles": {
    "contract": {
      "name": "contract",
      "environment": "dev",
      "protected_resource_metadata_url": "http://example.com/.well-known/oauth-protected-resource",
      "authorization_server_metadata_url": "http://example.com/.well-known/oauth-authorization-server/contract",
      "resource": "http://example.com/mcp-servers",
      "authorization_endpoint": "http://example.com/oauth/authorize",
      "token_endpoint": "http://example.com/oauth/token",
      "registration_endpoint": "http://example.com/oauth/register",
      "redirect_url": "http://127.0.0.1:8000/callback",
      "scopes": ["mcp:tools","mcp:resources"],
      "business_type": "contract",
      "client_name": "democli",
      "client_id": "legacy-client-id",
      "token": {
        "access_token": "legacy-token",
        "token_type": "Bearer"
      }
    }
  }
}`

	if err := os.WriteFile(store.Path(), []byte(legacyConfig), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := store.GetProfile("contract")
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}

	if got.DefaultIdentity != config.IdentityUser {
		t.Fatalf("default identity mismatch: got %q want %q", got.DefaultIdentity, config.IdentityUser)
	}
	if got.Identities.User.ClientID != "legacy-client-id" {
		t.Fatalf("legacy user client_id not migrated: got %q", got.Identities.User.ClientID)
	}
	if got.Identities.User.Token == nil || got.Identities.User.Token.AccessToken != "legacy-token" {
		t.Fatalf("legacy user token not migrated: got %+v", got.Identities.User.Token)
	}
	if got.Identities.User.RedirectURL != "http://127.0.0.1:8000/callback" {
		t.Fatalf("legacy redirect_url not migrated: got %q", got.Identities.User.RedirectURL)
	}
	if got.ClientName != "contract-cli" {
		t.Fatalf("legacy client_name should be normalized: got %q", got.ClientName)
	}
}

func TestDefaultDirUsesContractCLIConventions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CONTRACT_CLI_CONFIG_DIR", "")
	t.Setenv("DEMOCLI_CONFIG_DIR", "")

	dir, err := config.DefaultDir()
	if err != nil {
		t.Fatalf("DefaultDir() error = %v", err)
	}

	if !strings.HasSuffix(dir, filepath.Join(home, ".contract-cli")) {
		t.Fatalf("DefaultDir() = %q, want suffix %q", dir, filepath.Join(home, ".contract-cli"))
	}
}

func TestDefaultDirPrefersNewEnvAndFallsBackToLegacyEnv(t *testing.T) {
	t.Setenv("CONTRACT_CLI_CONFIG_DIR", "/tmp/contract-cli-config")
	t.Setenv("DEMOCLI_CONFIG_DIR", "/tmp/democli-config")

	dir, err := config.DefaultDir()
	if err != nil {
		t.Fatalf("DefaultDir() error = %v", err)
	}
	if dir != "/tmp/contract-cli-config" {
		t.Fatalf("DefaultDir() = %q, want /tmp/contract-cli-config", dir)
	}

	t.Setenv("CONTRACT_CLI_CONFIG_DIR", "")

	dir, err = config.DefaultDir()
	if err != nil {
		t.Fatalf("DefaultDir() legacy fallback error = %v", err)
	}
	if dir != "/tmp/democli-config" {
		t.Fatalf("DefaultDir() with legacy env = %q, want /tmp/democli-config", dir)
	}
}

func TestDefaultDirFallsBackToLegacyHomeDirectoryWhenPresent(t *testing.T) {
	home := t.TempDir()
	legacyDir := filepath.Join(home, ".democli")
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	t.Setenv("HOME", home)
	t.Setenv("CONTRACT_CLI_CONFIG_DIR", "")
	t.Setenv("DEMOCLI_CONFIG_DIR", "")

	dir, err := config.DefaultDir()
	if err != nil {
		t.Fatalf("DefaultDir() error = %v", err)
	}

	if dir != legacyDir {
		t.Fatalf("DefaultDir() = %q, want legacy dir %q", dir, legacyDir)
	}
}

func TestSecretsStoreRoundTrip(t *testing.T) {
	t.Parallel()

	secrets := config.NewSecretsStore(t.TempDir())
	const key = "contract.bot.app_secret"
	const value = "super-secret"

	if err := secrets.Set(key, value); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	got, ok, err := secrets.Get(key)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !ok || got != value {
		t.Fatalf("Get() = (%q, %v), want (%q, true)", got, ok, value)
	}

	if err := secrets.Delete(key); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	got, ok, err = secrets.Get(key)
	if err != nil {
		t.Fatalf("Get() after delete error = %v", err)
	}
	if ok || got != "" {
		t.Fatalf("expected secret to be deleted, got (%q, %v)", got, ok)
	}
}

func storeDir(t *testing.T, store *config.Store) string {
	t.Helper()
	return filepath.Dir(store.Path())
}
