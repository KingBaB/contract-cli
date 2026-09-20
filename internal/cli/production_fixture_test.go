package cli_test

import (
	"net/url"

	"cn.qfei/contract-cli/internal/config"
)

func testEnvironmentProfileFixture(profile config.Profile) config.Profile {
	profile.Environment = "test"
	profile.OpenPlatformBaseURL = "https://test-open.qtech.cn"
	profile.Resource = "https://test-open.qtech.cn"
	profile.AppTokenEndpoint = testEnvironmentFixtureURL(profile.AppTokenEndpoint, "https://test-open.qtech.cn")
	profile.ProtectedResourceMetadataURL = testEnvironmentFixtureURL(profile.ProtectedResourceMetadataURL, "https://test-open.qtech.cn")
	profile.AuthorizationServerMetadataURL = testEnvironmentFixtureURL(profile.AuthorizationServerMetadataURL, "https://test-myaccount.qtech.cn")
	profile.Identities.User.AuthorizationEndpoint = testEnvironmentFixtureURL(profile.Identities.User.AuthorizationEndpoint, "https://test-myaccount.qtech.cn")
	profile.Identities.User.DeviceAuthorizationEndpoint = testEnvironmentFixtureURL(profile.Identities.User.DeviceAuthorizationEndpoint, "https://test-myaccount.qtech.cn")
	profile.Identities.User.TokenEndpoint = testEnvironmentFixtureURL(profile.Identities.User.TokenEndpoint, "https://test-myaccount.qtech.cn")
	profile.Identities.User.RevocationEndpoint = testEnvironmentFixtureURL(profile.Identities.User.RevocationEndpoint, "https://test-myaccount.qtech.cn")
	profile.Identities.User.RegistrationEndpoint = testEnvironmentFixtureURL(profile.Identities.User.RegistrationEndpoint, "https://test-myaccount.qtech.cn")
	return profile
}

func testEnvironmentFixtureURL(rawURL, origin string) string {
	if rawURL == "" {
		return ""
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	productionOrigin, err := url.Parse(origin)
	if err != nil {
		panic(err)
	}
	parsed.Scheme = productionOrigin.Scheme
	parsed.Host = productionOrigin.Host
	return parsed.String()
}
