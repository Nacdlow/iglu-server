package forms

// EditAccountForm is the form for editing accounts in settings.
type EditAccountForm struct {
	FirstName  string `form:"fname" binding:"Required;AlphaDash;MinSize(2);MaxSize(50)"`
	LastName   string `form:"lname" binding:"OmitEmpty;AlphaDash;MinSize(1);MaxSize(50)"`
	Email      string `form:"email" binding:"Required;Email"`
	Password   string `form:"password"`
	RePassword string `form:"repassword"`
}
