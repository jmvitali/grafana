package setting

type OAuthInfo struct {
	ClientId, ClientSecret string
	Scopes                 []string
	AuthUrl, TokenUrl      string
    ReqTokenUrl            string
	Enabled                bool
	AllowedDomains         []string
	ApiUrl                 string
	AllowSignup            bool
    AutoSignUp             bool
}

type OAuther struct {
	GitHub, Google, Twitter, Cerberus bool
	OAuthInfos              map[string]*OAuthInfo
}

var OAuthService *OAuther
