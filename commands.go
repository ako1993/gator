package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ako1993/internal/config"
	"github.com/ako1993/internal/database"
	"github.com/google/uuid"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	handler map[string]func(*state, command) error
}

var commands_ = commands{
	make(map[string]func(*state, command) error),
}

func (c *commands) run(s *state, cmd command) error {
	err := c.handler[cmd.Name](s, cmd)
	if err != nil {
		fmt.Printf("ERROR COULD NOT RUN COMMAND:%v\n", err)
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handler[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Command Arguments Empty\n")
	}
	if cmd.Args[0] == "unknown" {
		os.Exit(1)
	}
	err := config.Set_user(*s.config, cmd.Args[0])
	if err != nil {
		fmt.Printf("ERROR SETTING USER NAME:%v", err)
		return err
	}
	fmt.Println("The user has been set!")
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.dbQueries.GetUserByName(context.Background(), s.config.Current_user_name)
		if err != nil {
			fmt.Printf("ERROR GETTING CURRENT USER:%v", err)
			return err
		}
		handler(s, cmd, user)
		return nil
	}
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Error! No Args provided")
	}
	name := cmd.Args[0]
	checkforuser, err := s.dbQueries.GetUserByName(context.Background(), name)
	if checkforuser.Name != "" {
		os.Exit(1)
	}
	if name == "unknown" {
		fmt.Println(name)
		os.Exit(1)
	}
	CreateUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	new_user, err := s.dbQueries.CreateUser(context.Background(), CreateUserParams)
	if err != nil {
		fmt.Printf("ERROR CREATING USER:%v", err)
		return err
	}
	config.Set_user(*s.config, new_user.Name)
	fmt.Println("The user was created")
	fmt.Printf("User name:%v\n", new_user.Name)
	fmt.Printf("Created At:%v\n", new_user.CreatedAt)
	fmt.Printf("Updated At:%v\n", new_user.UpdatedAt)
	fmt.Printf("User ID:%v\n", new_user.ID)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.dbQueries.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Printf("ERROR REMOVING ALL USERS:%v\n", err)
		os.Exit(1)
	}
	fmt.Println("Users cleared from Database!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.dbQueries.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("ERROR GETTING ALL USERS FROM DATABASE:%v", err)
		return err
	}
	current_user := s.config.Current_user_name
	for _, user := range users {
		if user.Name == current_user {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Printf("ERROR FETCHING RSS FEED IN HANDLER:%v\n", err)
		return err
	}
	for _, i := range feed.Channel.Item {
		fmt.Println(i)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		os.Exit(1)
	}
	feed_name := cmd.Args[0]
	feed_url := cmd.Args[1]
	new_id_ := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true, // must set to true to indicate it's not NULL
	}
	feed_params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feed_name,
		Url:       feed_url,
		UserID:    new_id_,
	}
	new_feed, err := s.dbQueries.CreateFeed(context.Background(), feed_params)
	if err != nil {
		fmt.Printf("ERROR ADDING NEW FEED TO DATABASE:%v\n", err)
		return err
	}
	fmt.Println("New feed created:")
	fmt.Printf("Feed ID:%v\n", new_feed.ID)
	fmt.Printf("Created At:%v\n", new_feed.CreatedAt)
	fmt.Printf("Updated At:%v\n", new_feed.UpdatedAt)
	fmt.Printf("Name:%v\n", new_feed.Name)
	fmt.Printf("URL:%v\n", new_feed.Url)
	fmt.Printf("User ID:%v\n", new_feed.UserID)
	new_feed_id_ := uuid.NullUUID{
		UUID:  new_feed.ID,
		Valid: true, // must set to true to indicate it's not NULL
	}
	feed_follows_params := database.CreateFeedFollowsRecordParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    new_id_,
		FeedID:    new_feed_id_,
	}
	new_follows_record, err := s.dbQueries.CreateFeedFollowsRecord(context.Background(), feed_follows_params)
	if err != nil {
		fmt.Printf("ERROR CREATING FEED FOLLOW AFTER ADDING FEED:%v", err)
		return err
	}
	fmt.Printf("New feed follows record created for:%v\n", new_follows_record.UserName)
	return nil
}

func HandlerFeeds(s *state, cmd command) error {
	feeds, err := s.dbQueries.GetAllFeeds(context.Background())
	if err != nil {
		fmt.Printf("ERROR FETCHING FEEDS:%v\n", err)
		return err
	}
	for _, feed := range feeds {
		user_name, err := s.dbQueries.GetUserNameByID(context.Background(), feed.UserID.UUID)
		if err != nil {
			fmt.Printf("ERROR GETTING USERNAME BY ID:%v\n", err)
			return err
		}
		fmt.Printf("Feed Name:%v\n", feed.Name)
		fmt.Printf("Feed URL:%v\n", feed.Url)
		fmt.Printf("Feed Name:%v\n", feed.Name)
		fmt.Printf("User Name:%v\n", user_name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return errors.New("ERROR NO ARGS PROVIDED")
	}
	url := cmd.Args[0]
	selected_feed, err := s.dbQueries.GetFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Printf("ERROR FETCHING FEED BY URL:%v", err)
		return err
	}
	new_id_ := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true, // must set to true to indicate it's not NULL
	}
	new_feed_id_ := uuid.NullUUID{
		UUID:  selected_feed.ID,
		Valid: true, // must set to true to indicate it's not NULL
	}
	create_feed_follow_params := database.CreateFeedFollowsRecordParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    new_id_,
		FeedID:    new_feed_id_,
	}
	new_feeds_follow_record, err := s.dbQueries.CreateFeedFollowsRecord(context.Background(), create_feed_follow_params)
	if err != nil {
		fmt.Printf("ERROR CREATING FEED FOLLOWS RECORD:%v", err)
		return err
	}
	fmt.Println("New feed follows record created!")
	fmt.Printf("Name of feed:%v", new_feeds_follow_record.FeedName)
	fmt.Printf("Current user:%v", new_feeds_follow_record.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	new_id_ := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true, // must set to true to indicate it's not NULL
	}
	feeds_being_followed, err := s.dbQueries.GetFeedFollowsForUser(context.Background(), new_id_)
	if err != nil {
		fmt.Printf("ERROR GETTING FEED FOLLOWS FOR USER :%v", err)
		return err
	}
	for _, feed := range feeds_being_followed {
		fmt.Println(feed.Name)
	}
	return nil
}
