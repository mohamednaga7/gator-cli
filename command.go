package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mohamednaga7/gator/internal/config"
	"github.com/mohamednaga7/gator/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	AvailableCommands map[string]func(s *State, cmd Command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func LoginHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 || strings.TrimSpace(cmd.Arguments[0]) == "" {
		return errors.New("username argument required")
	}

	_, err := s.DB.GetUserByName(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	err = s.Config.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("the user has been set")

	return nil
}

func RegisterHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 || strings.TrimSpace(cmd.Arguments[0]) == "" {
		return errors.New("username argument required")
	}

	name := strings.TrimSpace(cmd.Arguments[0])

	userFound := true

	_, err := s.DB.GetUserByName(context.Background(), name)
	if err != nil {
		if "sql: no rows in result set" == err.Error() {
			userFound = false
		} else {
			return err
		}
	}

	if userFound {
		return errors.New("user already exists")
	}

	user, err := s.DB.CreateUser(context.Background(), name)
	if err != nil {
		return err
	}

	err = s.Config.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("user %s with id %s has been created at %s\n", user.Name, user.ID.String(), time.Now().Format(time.RFC3339))

	return nil
}

func ResetHandler(s *State, _ Command) error {
	err := s.DB.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("all users have been deleted")

	return nil
}

func GetUsersHandler(s *State, _ Command) error {
	users, err := s.DB.GetAllUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		textToPrint := "* " + user.Name
		if user.Name == s.Config.CurrentUserName {
			textToPrint += " (current)"
		}
		fmt.Println(textToPrint)
	}

	return nil
}

func AddFeedHandler(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}

	name := cmd.Arguments[0]
	url := cmd.Arguments[1]

	createdFeed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: &user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}

	newUserFeedParams := database.AddFeedFollowParams{
		FeedID: &createdFeed.ID,
		UserID: &user.ID,
	}

	_, err = s.DB.AddFeedFollow(context.Background(), newUserFeedParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed created successfully:\n")

	return nil
}

func RSSHandler(_ *State, _ Command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func PrintFeedHandler(s *State, _ Command) error {
	feeds, err := s.DB.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("* %s - %s - %s\n", feed.Name, feed.Url, feed.UserName)
	}

	return nil
}

func FollowHandler(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("usage: follow <url>")
	}

	url := cmd.Arguments[0]

	feedItem, err := s.DB.GetFeedByUrl(context.Background(), url)
	if err != nil {
		if "sql: no rows in result set" == err.Error() {
			return errors.New("feed not found")
		}
		return err
	}

	newUserFeedParams := database.AddFeedFollowParams{
		FeedID: &feedItem.ID,
		UserID: &user.ID,
	}

	newUserFeed, err := s.DB.AddFeedFollow(context.Background(), newUserFeedParams)
	if err != nil {
		return err
	}

	fmt.Printf("Followed %s - %s\n", newUserFeed.UserName, newUserFeed.FeedName)

	return nil
}

func FeedFollowsHandler(s *State, _ Command, user database.User) error {
	feeds, err := s.DB.GetFeedFollowsByUserId(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.FeedName)
	}

	return nil
}

func UnfollowHandler(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("usage: unfollow <url>")
	}

	url := cmd.Arguments[0]

	feedItem, err := s.DB.GetFeedByUrl(context.Background(), url)
	if err != nil {
		if "sql: no rows in result set" == err.Error() {
			return errors.New("feed not found")
		}
		return err
	}

	deleteInput := database.DeleteByUserIdAndFeedIdParams{
		UserID: &user.ID,
		FeedID: &feedItem.ID,
	}

	err = s.DB.DeleteByUserIdAndFeedId(context.Background(), deleteInput)
	if err != nil {
		return err
	}

	fmt.Printf("Unfollowed %s\n", feedItem.Name)

	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	availableCmd, ok := c.AvailableCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return availableCmd(s, cmd)
}

func (c *Commands) Register(name string, f func(s *State, cmd Command) error) error {
	if c.AvailableCommands == nil {
		c.AvailableCommands = make(map[string]func(s *State, cmd Command) error)
	}
	c.AvailableCommands[name] = f
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	if res.StatusCode < 200 && res.StatusCode >= 400 {
		return nil, errors.New("error fetching Rss Feed")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed

	err = xml.Unmarshal(bodyBytes, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUserByName(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
