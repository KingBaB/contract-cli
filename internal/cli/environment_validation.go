package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"cn.qfei/contract-cli/internal/config"
	"cn.qfei/contract-cli/internal/credential"
)

func validateProfileEnvironment(profile config.Profile) error {
	preset, err := resolveEnvironment(profile.Environment)
	if err != nil {
		return profileEnvironmentError(profile.Name)
	}
	if !isExactEnvironmentResource(profile.OpenPlatformBaseURL, preset.OpenPlatformBaseURL) || !isExactEnvironmentResource(profile.Resource, preset.Resource) {
		return profileEnvironmentError(profile.Name)
	}
	openPlatformURLs := []string{
		profile.AppTokenEndpoint,
		profile.ProtectedResourceMetadataURL,
	}
	for _, rawURL := range openPlatformURLs {
		if !isEnvironmentOriginURL(rawURL, preset.OpenPlatformBaseURL, false) {
			return profileEnvironmentError(profile.Name)
		}
	}
	accountURLs := []string{
		profile.AuthorizationServerMetadataURL,
		profile.Identities.User.AuthorizationEndpoint,
		profile.Identities.User.DeviceAuthorizationEndpoint,
		profile.Identities.User.TokenEndpoint,
		profile.Identities.User.RevocationEndpoint,
		profile.Identities.User.RegistrationEndpoint,
	}
	for _, rawURL := range accountURLs {
		if !isEnvironmentOriginURL(rawURL, preset.AccountOrigin, false) {
			return profileEnvironmentError(profile.Name)
		}
	}
	return nil
}

func validateDeviceCredentialEnvironment(profile config.Profile, stored credential.DeviceCredential) error {
	profileName := profile.Name
	if stored.DeviceProfile != nil {
		snapshot, err := restoreDeviceProfile(profileName, stored.DeviceProfile)
		if err != nil || validateProfileEnvironment(snapshot) != nil || snapshot.Environment != profile.Environment {
			return profileEnvironmentError(profileName)
		}
	}
	return validatePendingEnvironment(profile, stored.Pending)
}

func validatePendingEnvironment(profile config.Profile, pending *credential.PendingTransaction) error {
	if pending != nil &&
		(!isProfileAccountURL(profile, pending.TokenEndpoint, true) ||
			!isProfileAccountURL(profile, pending.VerificationURIComplete, false)) {
		return profileEnvironmentError(profile.Name)
	}
	return nil
}

func isProfileAccountURL(profile config.Profile, rawURL string, required bool) bool {
	preset, err := resolveEnvironment(profile.Environment)
	return err == nil && isEnvironmentOriginURL(rawURL, preset.AccountOrigin, required)
}

func isExactEnvironmentResource(rawURL, expectedOrigin string) bool {
	parsed, err := parseEnvironmentURL(rawURL)
	if err != nil {
		return false
	}
	return parsed.Scheme+"://"+parsed.Host == expectedOrigin &&
		(parsed.EscapedPath() == "" || parsed.EscapedPath() == "/") && parsed.RawQuery == ""
}

func isEnvironmentOriginURL(rawURL, expectedOrigin string, required bool) bool {
	if strings.TrimSpace(rawURL) == "" {
		return !required
	}
	parsed, err := parseEnvironmentURL(rawURL)
	if err != nil {
		return false
	}
	return parsed.Scheme+"://"+parsed.Host == expectedOrigin
}

func parseEnvironmentURL(rawURL string) (*url.URL, error) {
	trimmed := strings.TrimSpace(rawURL)
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.Fragment != "" {
		return nil, errors.New("invalid environment URL")
	}
	if parsed.Hostname() != strings.TrimSuffix(parsed.Hostname(), ".") {
		return nil, errors.New("trailing-dot hosts are not allowed")
	}
	return parsed, nil
}

func profileEnvironmentError(profileName string) error {
	return fmt.Errorf("profile %q has unsupported or inconsistent environment endpoints; this build only enables test (prod temporarily disabled, dev blocked); configure matching test endpoints and authorize again", profileName)
}

// Match net/http's default redirect limit when the caller has no custom policy.
const defaultEnvironmentRedirectLimit = 10

func withEnvironmentNetworkGuard(client *http.Client, logger *slog.Logger) *http.Client {
	guarded := *client
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	guarded.Transport = environmentGuardTransport{next: transport, logger: logger}
	previous := client.CheckRedirect
	guarded.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		for _, previousRequest := range via {
			source := environmentForHost(previousRequest.URL.Hostname())
			target := environmentForHost(request.URL.Hostname())
			logger.Debug("validate environment redirect", "source_environment", source, "target_environment", target)
			if source != "" && target != "" && source != target {
				logger.Error("blocked cross-environment redirect", "source_environment", source, "target_environment", target)
				return errors.New("cross-environment redirect is not allowed")
			}
		}
		if previous != nil {
			return previous(request, via)
		}
		if len(via) >= defaultEnvironmentRedirectLimit {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	return &guarded
}

func environmentForHost(host string) string {
	switch strings.ToLower(strings.TrimSuffix(host, ".")) {
	case "open.qfei.cn", "myaccount.qfei.cn":
		return "prod"
	case "dev-open.qtech.cn", "dev-myaccount.qtech.cn":
		return "dev"
	case "test-open.qtech.cn", "test-myaccount.qtech.cn":
		return "test"
	default:
		return ""
	}
}

func (a *App) loadEnvironmentDeviceCredential(profile config.Profile, store credential.Store) (credential.DeviceCredential, error) {
	stored, err := store.Load(profile.Name)
	if err != nil {
		return credential.DeviceCredential{}, err
	}
	if err := validateDeviceCredentialEnvironment(profile, stored); err != nil {
		a.logger.Error("reject mismatched environment Device credential", "profile", profile.Name, "error", err.Error())
		return credential.DeviceCredential{}, err
	}
	return stored, nil
}

func (a *App) validateDeviceCredentialEnvironmentIfAvailable(profile config.Profile) error {
	store, available, err := a.deviceCredentialStoreIfAvailable()
	if err != nil || !available {
		return err
	}
	stored, err := store.Load(profile.Name)
	if errors.Is(err, credential.ErrCredentialNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := validateDeviceCredentialEnvironment(profile, stored); err != nil {
		a.logger.Error("reject mismatched environment Device credential", "profile", profile.Name, "error", err.Error())
		return err
	}
	return nil
}

func (a *App) deviceCredentialStoreIfAvailable() (credential.Store, bool, error) {
	if a.credentialStore != nil {
		return a.credentialStore, true, nil
	}
	for _, name := range []string{"SKILL_SESSION_WORKSPACE", "CODEBUDDY_SESSION_ID", "SESSION_ID"} {
		if value, ok := a.lookupEnv(name); ok && strings.TrimSpace(value) != "" {
			store, err := a.deviceCredentials()
			return store, err == nil, err
		}
	}
	return nil, false, nil
}

func (a *App) profileRequiresAuthenticationReset(profile config.Profile, environment string) (bool, error) {
	if profile.Environment != environment || validateProfileEnvironment(profile) != nil {
		return true, nil
	}
	store, available, err := a.deviceCredentialStoreIfAvailable()
	if err != nil || !available {
		return false, err
	}
	stored, err := store.Load(profile.Name)
	if errors.Is(err, credential.ErrCredentialNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return validateDeviceCredentialEnvironment(profile, stored) != nil, nil
}

func (a *App) clearProfileAuthenticationState(profileName string) error {
	a.logger.Info("clear mismatched environment profile authentication state", "profile", profileName)
	if err := a.secrets.Delete(config.AppSecretKey(profileName)); err != nil {
		a.logger.Error("clear app secret failed", "profile", profileName, "error", err.Error())
		return err
	}
	store, available, err := a.deviceCredentialStoreIfAvailable()
	if err != nil {
		a.logger.Error("resolve Device credential store for reset failed", "profile", profileName, "error", err.Error())
		return err
	}
	if available {
		if err := store.Delete(profileName); err != nil {
			a.logger.Error("clear Device credential failed", "profile", profileName, "error", err.Error())
			return err
		}
	}
	return nil
}

// Block disabled environments before sending credentials, including redirects.
type environmentGuardTransport struct {
	next   http.RoundTripper
	logger *slog.Logger
}

func (transport environmentGuardTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	host := strings.ToLower(strings.TrimSuffix(request.URL.Hostname(), "."))
	transport.logger.Debug("environment network request", "method", request.Method, "host", host)
	environment := environmentForHost(host)
	if environment == "dev" || (environment == "prod" && !productionEnvironmentEnabled) {
		transport.logger.Error("blocked unsupported environment request", "method", request.Method, "host", host)
		return nil, fmt.Errorf("build blocks unsupported environment host %q", host)
	}
	return transport.next.RoundTrip(request)
}
