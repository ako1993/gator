package config

import (
	"testing"
)

func TestSetUser(t *testing.T) {
	user_config := Config{
		Db_url:            "www.test.com",
		Current_user_name: "",
	}
	Set_user(user_config, "Tommy")
}
