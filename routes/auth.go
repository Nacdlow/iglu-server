package routes

import (
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

// LogoutHandler handles logging out.
func LogoutHandler(ctx *macaron.Context, sess session.Store) {
	sess.Set("auth", Unauthenticated)
	ctx.Redirect("/")
}

// LoginHandler handles rendering the login page.
func LoginHandler(ctx *macaron.Context, sess session.Store) {
	if sess.Get("auth") == LoggedIn {
		ctx.Redirect("/dashboard")
		return
	}
	ctx.HTML(200, "index")
}

// PostLoginHandler handles the post login page.
func PostLoginHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store,
	form forms.SignInForm, f *session.Flash) {
	if sess.Get("auth") == LoggedIn {
		ctx.Redirect("/dashboard")
		return
	}
	var u *models.User
	var err error
	if u, err = models.GetUser(form.Email); err != nil {
		f.Error("Invalid username or password")
		ctx.Redirect("/")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(form.Password))
	if err != nil {
		f.Error("Invalid username or password")
		ctx.Redirect("/")
		return
	}

	err = sess.Set("auth", LoggedIn)
	if err != nil {
		panic(err)
	}
	err = sess.Set("username", u.Username)
	if err != nil {
		panic(err)
	}

	ctx.Redirect("/dashboard")
}

// RegisterHandler handles the registration page.
func RegisterHandler(ctx *macaron.Context, sess session.Store) {
	if sess.Get("auth") == LoggedIn {
		ctx.Redirect("/dashboard")
		return
	}
	ctx.HTML(200, "register")
}

// PostRegisterHandler handles the post registration page.
/*func PostRegisterHandler(ctx *macaron.Context, sess session.Store) {
	if sess.Get("auth") == LoggedIn {
		ctx.Redirect("/dashboard")
		return
	}
	ctx.Redirect("/login")
} */
