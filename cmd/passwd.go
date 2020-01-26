package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
)

// CmdPasswd represents the command to change a user's password.
var CmdPasswd = &cli.Command{
	Name:   "passwd",
	Usage:  "Change a user's password",
	Action: passwd,
}

func passwd(c *cli.Context) (err error) {
	settings.LoadConfig()
	engine := models.SetupEngine()
	defer engine.Close()

	fmt.Printf("Username (email): ")
	var username string
	fmt.Scanln(&username)
	u, err := models.GetUser(username)
	if err != nil {
		panic(err)
	}
	var inPass, inPass2 string
	fmt.Printf("New password (WILL echo): ")
	fmt.Scanln(&inPass)
	fmt.Printf("Confirm new password: ")
	fmt.Scanln(&inPass2)
	if inPass != inPass2 {
		fmt.Println("Does not match! Password remains unchanged.")
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(inPass), 10)
	if err != nil {
		panic(err)
	}
	u.Password = string(pass)
	models.UpdateUser(u)

	fmt.Println()
	fmt.Println("Password updated!")
	return nil
}
