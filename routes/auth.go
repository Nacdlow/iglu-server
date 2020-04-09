package routes

import (
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/models/forms"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/tokens"

	"net/http"

	"github.com/go-macaron/binding"
	"golang.org/x/crypto/bcrypt"
	macaron "gopkg.in/macaron.v1"
)

// LogoutHandler handles logging out.
func LogoutHandler(ctx *macaron.Context, sess session.Store) {
	err := sess.Set("auth", Unauthenticated)
	if err != nil {
		panic(err)
	}
	ctx.Redirect("/")
}

// LoginHandler handles rendering the login page.
func LoginHandler(ctx *macaron.Context, sess session.Store) {
	if sess.Get("auth") == LoggedIn {
		ctx.Redirect("/dashboard")
		return
	}
	ctx.HTML(http.StatusOK, "index")
}

// ForgotHandler handles rendering the forgot password page.
func ForgotHandler(ctx *macaron.Context, sess session.Store) {
	if sess.Get("auth") == LoggedIn {
		ctx.Redirect("/dashboard")
		return
	}
	ctx.Data["EngineerName"] = settings.Config.Get("Engineer.Name")
	ctx.Data["EngineerEmail"] = settings.Config.Get("Engineer.Email")
	ctx.Data["EngineerPhone"] = settings.Config.Get("Engineer.Phone")
	ctx.HTML(http.StatusOK, "forgot")
}

// PostLoginHandler handles the post login page.
func PostLoginHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store,
	form forms.SignInForm, errs binding.Errors, f *session.Flash) {
	if len(errs) > 0 {
		f.Error("Missing required fields!")
		ctx.Redirect("/")
		return
	}
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
	ctx.HTML(http.StatusOK, "register")
}

// AddUserHandler handles adding a new user from registeration.
func AddUserHandler(ctx *macaron.Context, form forms.RegisterForm,
	errs binding.Errors, f *session.Flash) {
	if len(errs) > 0 {
		f.Error("Missing required fields!")
		ctx.Redirect("/")
		return
	}
	ok := tokens.CheckAndConsumeKey(form.InviteCode)
	if !ok {
		f.Error("Invalid invite code. Please ask for an invite code from the home owner.")
		ctx.Redirect("/register")
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
	if err != nil {
		panic(err)
	}
	user := &models.User{
		Username:  form.Email,
		Password:  string(pass),
		FirstName: form.FirstName,
		LastName:  form.LastName,
		FontSize:  "medium",
		FontFace:  "default-roboto",
		Avatar:    "/img/profiles/penguin_pixabay.jpg",
	}
	err = models.AddUser(user)
	if err != nil {
		panic(err)
	}
	ctx.Redirect("/dashboard")
}
