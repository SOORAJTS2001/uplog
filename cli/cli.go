package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"uplog/commands"
	"uplog/models"
	"uplog/utils"
	"uplog/setup"
)

var respJSON models.SessionCreateResponse
var tag *string
var pollInterval = 200 * time.Millisecond
var sql_object *sql.DB
var command string
var commandArgs []string



func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	home := setup.UserHomeDir()
	if home == "" {
		fmt.Println("cannot get user home directory")
		os.Exit(1)
	}
	setup.SetupDirectories(home)
	sql_object = setup.InitDB(sql_object,home)

	cmd := os.Args[1]

	switch cmd {

	case "run":
		// commands.RunCmd(sql_object,home, command, commandArgs,respJSON,pollInterval,tag)
		command,commandArgs,pollInterval,tag = utils.RunCmdWithFlags(home, os.Args[2:],pollInterval,tag)
		commands.RunCmd(sql_object,home,command,commandArgs,respJSON,pollInterval,tag)

	case "list":
		commands.ListCmd(sql_object,home)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("usage: uplog delete <session_id>")
			os.Exit(1)
		}
		commands.DeleteCmd(sql_object,home, os.Args[2])

	default:
		usage()
	}
}

func usage() {
	fmt.Println("uplog - simple CLI")
	fmt.Println("Usage:")
	fmt.Println("  uplog run <cmd> [args...]    Run a command and upload logs")
	fmt.Println("  uplog list                   List uploaded sessions")
}
