package forms

// SignInForm is for the sign in page, receiving the email and password of the
// user, also storing if the user wants to be signed in indefinitely
type SignInForm struct {
	Email      string `form:"email" binding:"Required;Email"`
	Password   string `form:"password" binding:"Required;MinSize(8);MaxSize(50)"`
	RememberMe bool   `form:"remember"`
}

// RegisterForm ...
type RegisterForm struct {
	FirstName  string `form:"fname" binding:"Required;AlphaDash;MinSize(2);MaxSize(50)"`
	LastName   string `form:"lname" binding:"OmitEmpty;AlphaDash;MinSize(2);MaxSize(50)"`
	Email      string `form:"email" binding:"Required;Email"`
	Password   string `form:"password" binding:"Required;MinSize(8);MaxSize(50)"`
	RePassword string `form:"repassword" binding:"Required;MinSize(8);MaxSize(50)"`
}

// ForgotPassword ...
type ForgotPassword struct {
	Email string `form:"email" binding:"Required;Email"`
}

// NewPassword ...
type NewPassword struct {
	NewPassword string `form:"newpassword" binding:"Required;MinSize(8);MaxSize(50)"`
	RePassword  string `form:"renewpassword" binding:"Required;MinSize(8);MaxSize(50)"`
}
