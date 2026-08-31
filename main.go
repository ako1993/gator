package main

import (
	"fmt"

	"github.com/ako1993/internal/config"
)

func main() {
	user_config, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR READING FILE:%v", err)
		return
	}
	err = config.Set_user(user_config, "Andreas")
	if err != nil {
		fmt.Printf("ERROR WRITING TO FILE:%v", err)
		return
	}
	user_config, err = config.Read()
	if err != nil {
		fmt.Printf("ERROR READING TO FILE THE SECOND TIME:%v", err)
		return
	}
	fmt.Println(user_config.Current_user_name)
	fmt.Println(user_config.Db_url)
}
