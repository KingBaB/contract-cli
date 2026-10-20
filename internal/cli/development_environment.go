package cli

const (
	developmentEnvironment        = "dev"
	developmentOpenPlatformOrigin = "https://dev-open.qtech.cn"
	developmentAccountOrigin      = "https://dev-myaccount.qtech.cn"
)

// developmentEnvironmentPreset 仅为显式启用 dev 标签的本地构建提供开发环境预设。
func developmentEnvironmentPreset(name string) (environmentPreset, bool) {
	if !developmentBuildEnabled || name != developmentEnvironment {
		return environmentPreset{}, false
	}
	return environmentPreset{
		OpenPlatformBaseURL:            developmentOpenPlatformOrigin,
		AppTokenEndpoint:               developmentOpenPlatformOrigin + "/open-apis/auth/v3/tenant_access_token/internal",
		AuthorizationServerMetadataURL: developmentAccountOrigin + "/.well-known/oauth-authorization-server/contract",
		Resource:                       developmentOpenPlatformOrigin,
		RedirectURL:                    "http://127.0.0.1:8000/callback",
		Scopes:                         []string{"cli:tools", "cli:resources"},
		BusinessType:                   "contract",
		ClientName:                     "contract-cli",
	}, true
}

func supportedEnvironmentNames() string {
	if developmentBuildEnabled {
		return "prod, dev"
	}
	return "prod"
}

func configurationEnvironmentFlag() string {
	if developmentBuildEnabled {
		return "--env <prod|dev>"
	}
	return "--env <prod>"
}

// developmentOriginAllowed 只放行固定 HTTPS origin；端口、用户信息和尾点域名仍拒绝。
func developmentOriginAllowed(rawURL string) bool {
	return developmentBuildEnabled &&
		(isProductionOriginURL(rawURL, developmentOpenPlatformOrigin, true) ||
			isProductionOriginURL(rawURL, developmentAccountOrigin, true))
}
