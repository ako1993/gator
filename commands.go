package main

import (
	"errors"
	"fmt"

	"github.com/ako1993/internal/config"
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
	err := config.Set_user(*s.config, cmd.Args[0])
	if err != nil {
		fmt.Printf("ERROR SETTING USER NAME:%v", err)
		return err
	}
	fmt.Println("The user has been set!")
	return nil
}
