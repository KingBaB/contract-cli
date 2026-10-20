package cli_test

import (
	"net/url"

	"cn.qfei/contract-cli/internal/config"
)

func productionProfileFixture(profile config.Profile) config.Profile {
	profile.Environment = "prod"
	profile.OpenPlatformBaseURL = "https://open.qfei.cn"
	profile.Resource = "https://open.qfei.cn"
	profile.AppTokenEndpoint = productionFixtureURL(profile.AppTokenEndpoint, "https://open.qfei.cn")
	profile.ProtectedResourceMetadataURL = productionFixtureURL(profile.ProtectedResourceMetadataURL, "https://open.qfei.cn")
	profile.AuthorizationServerMetadataURL = productionFixtureURL(profile.AuthorizationServerMetadataURL, "https://myaccount.qfei.cn")
	profile.Identities.User.AuthorizationEndpoint = productionFixtureURL(profile.Identities.User.AuthorizationEndpoint, "https://myaccount.qfei.cn")
	profile.Identities.User.DeviceAuthorizationEndpoint = productionFixtureURL(profile.Identities.User.DeviceAuthorizationEndpoint, "https://myaccount.qfei.cn")
	profile.Identities.User.TokenEndpoint = productionFixtureURL(profile.Identities.User.TokenEndpoint, "https://myaccount.qfei.cn")
	profile.Identities.User.RevocationEndpoint = productionFixtureURL(profile.Identities.User.RevocationEndpoint, "https://myaccount.qfei.cn")
	profile.Identities.User.RegistrationEndpoint = productionFixtureURL(profile.Identities.User.RegistrationEndpoint, "https://myaccount.qfei.cn")
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
