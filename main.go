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

	commands := Commands{}

	err = registerHandlers(&commands)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

func registerHandlers(commands *Commands) error {
	if err := commands.Register("login", LoginHandler); err != nil {
		return err
	}
	if err := commands.Register("register", RegisterHandler); err != nil {
		return err
	}
	if err := commands.Register("reset", ResetHandler); err != nil {
		return err
	}
	if err := commands.Register("users", GetUsersHandler); err != nil {
		return err
	}
	if err := commands.Register("agg", RSSHandler); err != nil {
		return err
	}
	if err := commands.Register("addfeed", middlewareLoggedIn(AddFeedHandler)); err != nil {
		return err
	}
	if err := commands.Register("feeds", PrintFeedHandler); err != nil {
		return err
	}
	if err := commands.Register("follow", middlewareLoggedIn(FollowHandler)); err != nil {
		return err
	}
	if err := commands.Register("following", middlewareLoggedIn(FeedFollowsHandler)); err != nil {
		return err
	}
	if err := commands.Register("unfollow", middlewareLoggedIn(UnfollowHandler)); err != nil {
		return err
	}
	return nil
}
