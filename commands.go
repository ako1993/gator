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
