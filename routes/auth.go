package routes

import (
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

// LoginHandler handles rendering the login page.
func LoginHandler(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

// PostLoginHandler handles the post login page.
func PostLoginHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store) {
	if sess.Get("auth") == LoggedIn {
		ctx.Redirect("/dashboard")
		return
	}
	// TODO authenticate

	ctx.Redirect("/dashboard")
}

// RegisterHandler handles the registration page.
func RegisterHandler(ctx *macaron.Context) {
	ctx.HTML(200, "register")
}

// PostRegisterHandler handles the post registration page.
func PostRegisterHandler(ctx *macaron.Context) {
	ctx.Redirect("/login")
}
