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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		os.Exit(1)
	}
	feed_name := cmd.Args[0]
	feed_url := cmd.Args[1]
	current_user, err := s.dbQueries.GetUserByName(context.Background(), s.config.Current_user_name)
	if err != nil {
		fmt.Printf("ERROR GETTING CURRENT USER:%v", err)
		return err
	}
	new_id_ := uuid.NullUUID{
		UUID:  current_user.ID,
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
