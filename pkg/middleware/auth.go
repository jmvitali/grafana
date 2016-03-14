package middleware

import (
	"net/url"
	"strings"

	"github.com/Unknwon/macaron"

	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

type AuthOptions struct {
	ReqGrafanaAdmin bool
	ReqSignedIn     bool
}

func getRequestUserId(c *Context) int64 {
	userId := c.Session.Get(SESS_KEY_USERID)

	if userId != nil {
		return userId.(int64)
	}

	return 0
}

func getApiKey(c *Context) string {
	header := c.Req.Header.Get("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		key := parts[1]
		return key
	}

	return ""
}

func setRedirect(c *Context) {
	if setting.AuthAutoSignUpEnabled && setting.AuthAutoSignUpName != "" {
		c.Redirect(setting.AppSubUrl + "/login/" + setting.AuthAutoSignUpName )
	} else {
		c.Redirect(setting.AppSubUrl + "/login")
	}
}

func accessForbidden(c *Context) {
	if c.IsApiRequest() {
		c.JsonApiErr(403, "Permission denied", nil)
		return
	}

	c.SetCookie("redirect_to", url.QueryEscape(setting.AppSubUrl+c.Req.RequestURI), 0, setting.AppSubUrl+"/")
	setRedirect(c)
}

func notAuthorized(c *Context) {
	if c.IsApiRequest() {
		c.JsonApiErr(401, "Unauthorized", nil)
		return
	}

	c.SetCookie("redirect_to", url.QueryEscape(setting.AppSubUrl+c.Req.RequestURI), 0, setting.AppSubUrl+"/")
	setRedirect(c)
}

func RoleAuth(roles ...m.RoleType) macaron.Handler {
	return func(c *Context) {
		ok := false
		for _, role := range roles {
			if role == c.OrgRole {
				ok = true
				break
			}
		}
		if !ok {
			accessForbidden(c)
		}
	}
}

func Auth(options *AuthOptions) macaron.Handler {
	return func(c *Context) {

		if !c.IsSignedIn && options.ReqSignedIn && !c.AllowAnonymous {
			// Auto Login with Auth name configured
			if setting.AuthAutoSignUpEnabled && setting.AuthAutoSignUpName != ""{
				c.SetCookie("redirect_to", url.QueryEscape(setting.AppSubUrl+c.Req.RequestURI), 0, setting.AppSubUrl+"/")
	 			c.Redirect(setting.AppSubUrl + "/login/" + setting.AuthAutoSignUpName )
	 		}

			notAuthorized(c)
			return
		}

		if !c.IsGrafanaAdmin && options.ReqGrafanaAdmin {
			accessForbidden(c)
			return
		}
	}
}
