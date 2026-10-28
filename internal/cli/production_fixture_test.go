package cli_test

import (
	"net/url"

	"cn.qfei/contract-cli/internal/config"
)

func productionProfileFixture(profile config.Profile) config.Profile {
	profile.Environment = "dev"
	profile.OpenPlatformBaseURL = "https://dev-open.qtech.cn"
	profile.Resource = "https://dev-open.qtech.cn"
	profile.AppTokenEndpoint = productionFixtureURL(profile.AppTokenEndpoint, "https://dev-open.qtech.cn")
	profile.ProtectedResourceMetadataURL = productionFixtureURL(profile.ProtectedResourceMetadataURL, "https://dev-open.qtech.cn")
	profile.AuthorizationServerMetadataURL = productionFixtureURL(profile.AuthorizationServerMetadataURL, "https://dev-myaccount.qtech.cn")
	profile.Identities.User.AuthorizationEndpoint = productionFixtureURL(profile.Identities.User.AuthorizationEndpoint, "https://dev-myaccount.qtech.cn")
	profile.Identities.User.DeviceAuthorizationEndpoint = productionFixtureURL(profile.Identities.User.DeviceAuthorizationEndpoint, "https://dev-myaccount.qtech.cn")
	profile.Identities.User.TokenEndpoint = productionFixtureURL(profile.Identities.User.TokenEndpoint, "https://dev-myaccount.qtech.cn")
	profile.Identities.User.RevocationEndpoint = productionFixtureURL(profile.Identities.User.RevocationEndpoint, "https://dev-myaccount.qtech.cn")
	profile.Identities.User.RegistrationEndpoint = productionFixtureURL(profile.Identities.User.RegistrationEndpoint, "https://dev-myaccount.qtech.cn")
	return profile
}

func productionFixtureURL(rawURL, origin string) string {
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
