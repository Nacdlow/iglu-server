package routes

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/Nacdlow/plugin-sdk"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/plugin"
)

// LoginStatus represents the status of a user's session.
type LoginStatus int64

// LoginStatus enums.
const (
	Unauthenticated = iota
	LoggedIn
	Verification // For OTP
)

// ContextInit initialises the Macaron context to load the authenticated user
// from the database, and set other page fields such as the app name.
func ContextInit() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store) {
		ctx.Data["AppName"] = "iglü"
		if sess.Get("auth") == LoggedIn {
			if user, err := models.GetUser(sess.Get("username").(string)); err == nil {
				ctx.Data["User"] = user
				ctx.Data["LoggedIn"] = 1
				if user.Role == models.AdminRole {
					ctx.Data["IsAdmin"] = 1
				}
			} else {
				// Logged in user does not exist
				err := sess.Set("auth", Unauthenticated)
				if err != nil {
					panic(err)
				}
				ctx.Redirect("/")
				return
			}
		}
		var extraCSS, extraJS strings.Builder

		// Load WebExtensions from plugins
		for _, pl := range plugin.GetLoadedPlugins() {
			if pl.State != plugin.Running {
				continue
			}
			exts := pl.Plugin.GetWebExtensions()
			if exts != nil {
				for _, ext := range exts {
					if matchPathRegex(ext.PathMatchRegex, ctx.Req.URL.Path) {
						switch ext.Type {
						case sdk.CSS:
							extraCSS.WriteString(fmt.Sprintf("/* Injected by %s */\n", pl.Manifest.Name))
							extraCSS.WriteString(ext.Source)
							extraCSS.WriteString("\n\n")
						case sdk.JavaScript:
							extraJS.WriteString(fmt.Sprintf("/* Injected by %s */\n", pl.Manifest.Name))
							extraJS.WriteString(ext.Source)
							extraJS.WriteString("\n\n")
						}
					}
				}
			}
		}

		ctx.Data["ExtraCSS"] = template.CSS(extraCSS.String())
		ctx.Data["ExtraJS"] = template.JS(extraJS.String())
	}
}

func matchPathRegex(regex, path string) bool {
	if regex == "*" || regex == "" {
		return true
	}
	r := regexp.MustCompile(regex)
	return r.Match([]byte(path))

}

// RequireAdmin is a per-route middleware which requires the user to be an
// admin.
func RequireAdmin(ctx *macaron.Context, sess session.Store) {
	if u, err := models.GetUser(sess.Get("username").(string)); err == nil {
		if u.Role != models.AdminRole {
			ctx.Redirect("/")
			return
		}
	} else {
		ctx.Redirect("/")
		return
	}
}

// RequireLogin is a per-route middleware which requires the user to be logged
// in.
func RequireLogin(ctx *macaron.Context, sess session.Store) {
	if sess.Get("auth") != LoggedIn {
		ctx.Redirect("/")
		return
	}
}
