package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"

	"gitlab.com/group-nacdlow/nacdlow-server/models"
	"gitlab.com/group-nacdlow/nacdlow-server/modules/settings"
)

// CmdDumpfakes represents a command which adds fake data to the database.
var CmdDumpfakes = &cli.Command{
	Name:    "dumpfakedata",
	Aliases: []string{"dumpfakes"},
	Usage:   "A tool which allows you to dump fake data to the database.",
	Action:  dumpfakes,
}

func dumpfakes(c *cli.Context) (err error) {
	settings.LoadConfig()
	engine := models.SetupEngine()
	defer fmt.Println("Saving...")
	defer engine.Close()
	defer fmt.Println("Done")

	fmt.Println("========== WARNING ==========")
	fmt.Println("This tool allows you to dump randomly generated fake data to your database.")
	fmt.Println("Are you sure you want to continue?")
	fmt.Printf("Type YES (in uppercase): ")
	var verify string
	fmt.Scanln(&verify)
	if verify != "YES" {
		fmt.Println("Cancelling...")
		return nil
	}

	fmt.Println()
	fmt.Println("Type \"quit\" to exit")
	fmt.Println()
	for {
		var table, entriesStr string

		fmt.Println("Available tables: Device, Room, RoomStat, Statistic, User")
		fmt.Printf("Which table would you like to fake? ")
		fmt.Scanln(&table)
		if strings.ToLower(table) == "quit" {
			os.Exit(0)
		}
		fmt.Printf("How much entries do you want to generate? ")
		fmt.Scanln(&entriesStr)
		entries, err := strconv.Atoi(entriesStr)
		if err != nil {
			fmt.Println("Invalid number!")
			continue
		}

		table = strings.ToLower(table)
		switch table {
		case "device":
			for i := 0; i < entries; i++ {
				models.AddDevice(models.GetFakeDevice())
			}
		case "room":
			for i := 0; i < entries; i++ {
				models.AddRoom(models.GetFakeRoom())
			}
		case "roomstat":
			for i := 0; i < entries; i++ {
				models.AddRoomStat(models.GetFakeRoomStat())
			}
		case "statistic":
			for i := 0; i < entries; i++ {
				models.AddStat(models.GetFakeStat())
			}
		case "user":
			for i := 0; i < entries; i++ {
				models.AddUser(models.GetFakeUser())
			}
		default:
			fmt.Println("Invalid table name!")
		}
		fmt.Println("Added!")
		fmt.Println()
	}
}
