package setting

type OAuthInfo struct {
	ClientId, ClientSecret string
	Scopes                 []string
	AuthUrl, TokenUrl      string
    AllowInsecureCert      bool
	Enabled                bool
	AllowedDomains         []string
	ApiUrl                 string
	AllowSignup            bool
}

type OAuther struct {
	GitHub, Google, Twitter, Cerberus bool
	OAuthInfos              map[string]*OAuthInfo
}

var OAuthService *OAuther
