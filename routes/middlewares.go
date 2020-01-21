package routes

import (
	"github.com/go-macaron/session"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	macaron "gopkg.in/macaron.v1"
)

type LoginStatus int64

const (
	Unauthenticated = iota
	LoggedIn
)

func ContextInit() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store) {
		ctx.Data["AppName"] = "AppName"
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

func RequireLogin(ctx *macaron.Context, sess session.Store) {
	if sess.Get("auth") != LoggedIn {
		ctx.Redirect("/")
		return
	}
}
