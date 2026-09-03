package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ako1993/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	user_config, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR READING FILE:%v", err)
		return
	}
	var state state
	state.config = &user_config
	commands_.register("login", handlerLogin)
	if len(os.Args) < 2 {
		fmt.Println(errors.New("ERROR NOT ENOUGH ARGS PROVIDED"))
		os.Exit(1)
		return
	}
	input := os.Args
	command_name := input[1]
	command_args := input[2:]
	user_command := command{
		Name: command_name,
		Args: command_args,
	}
	if command_name == "login" && len(command_args) < 1 {
		fmt.Println(errors.New("Error A USERNAME IS REQUIRED TO LOGIN"))
		os.Exit(1)
	}
	err = commands_.run(&state, user_command)
	if err != nil {
		fmt.Printf("ERROR RUNNING COMMAND:%v", err)
		return
	}
}
