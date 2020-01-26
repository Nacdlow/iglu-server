package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
)

// CmdAdduser represents the command which allows adding new users to the
// database.
var CmdAdduser = &cli.Command{
	Name:    "adduser",
	Aliases: []string{"useradd"},
	Usage:   "Add a user to the database",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "isadmin",
			Value: true,
			Usage: "Whether the new user is an admin or not.",
		},
	},
	Action: adduser,
}

func adduser(c *cli.Context) (err error) {
	settings.LoadConfig()
	engine := models.SetupEngine()
	defer engine.Close()

	u := new(models.User)
	if c.Bool("isadmin") {
		fmt.Println("Creating a new admin user")
		u.Role = models.AdminRole
	} else {
		fmt.Println("Creating a regular (normal) user")
		u.Role = models.NormalRole
	}
	fmt.Printf("Username (email): ")
	fmt.Scanln(&u.Username)
	fmt.Printf("First name: ")
	fmt.Scanln(&u.FirstName)
	fmt.Printf("Last name: ")
	fmt.Scanln(&u.LastName)
	fmt.Printf("Password (WILL echo): ")
	var inPass, inPass2 string
	fmt.Scanln(&inPass)
	fmt.Printf("Confirm new password: ")
	fmt.Scanln(&inPass2)
	if inPass != inPass2 {
		fmt.Println("Does not match! User not added.")
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(inPass), 10)
	if err != nil {
		panic(err)
	}
	u.Password = string(pass)
	models.AddUser(u)

	fmt.Println()
	fmt.Println("User added!")
	return nil
}
