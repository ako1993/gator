package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

const configFileName = "/.gatorconfig.json"

func Read() (Config, error) {
	var user_config Config
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("ERROR GETTING FILE PATH:%v", err)
		return user_config, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("ERROR READING FILE:%v", err)
		return user_config, err
	}
	err = json.Unmarshal(content, &user_config)
	if err != nil {
		fmt.Printf("ERROR DECODING JSON:%v", err)
		return user_config, err
	}
	return user_config, nil
}

func Set_user(user_config Config, user_name string) error {
	user_config.Current_user_name = user_name
	data, err := json.Marshal(user_config)
	if err != nil {
		fmt.Printf("ERROR MARSHALLING DATA:%v", err)
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("ERROR GETTING FILE PATH:%v", err)
		return err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Printf("ERROR WRITING TO FILE:%v", err)
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("ERROR:%v", err)
		return "", err
	}
	path := home_dir + configFileName
	return path, nil
}
