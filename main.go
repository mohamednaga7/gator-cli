package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mohamednaga7/gator/internal/database"

	"github.com/mohamednaga7/gator/internal/config"
)

func main() {
	appConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", appConfig.DBURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	appState := State{
		DB:     dbQueries,
		Config: &appConfig,
	}

	commands := Commands{
		AvailableCommands: map[string]func(s *State, cmd Command) error{
			"login":    LoginHandler,
			"register": RegisterHandler,
			"reset":    ResetHandler,
			"users":    GetUsersHandler,
			"agg":      RSSHandler,
			"addfeed":  AddFeedHandler,
			"feeds":    PrintFeedHandler,
			"follow":   FollowHandler,
		},
	}

	args := os.Args

	if len(args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	err = commands.Run(&appState, Command{Name: args[1], Arguments: args[2:]})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
