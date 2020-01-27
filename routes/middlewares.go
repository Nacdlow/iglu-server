package routes

import (
	"github.com/go-macaron/session"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	macaron "gopkg.in/macaron.v1"
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
			} else {
				// Logged in user does not exist
				sess.Set("auth", Unauthenticated)
				ctx.Redirect("/")
				return
			}
		}
	}
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
