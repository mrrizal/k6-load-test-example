package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func generateUser() {
	var users []User
	for i := 0; i < 1000; i++ {
		var user User
		user.Username = fmt.Sprintf("username_%d", i)
		user.Password = fmt.Sprintf("password%d", i)
		user.FirstName = fmt.Sprintf("first_name_%d", i)
		user.LastName = fmt.Sprintf("last_name_%d", i)
		users = append(users, user)
	}

	usersJSON, err := json.MarshalIndent(&users, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	_ = ioutil.WriteFile("users.json", usersJSON, 0644)
}

func main() {
	genUser := flag.Bool("generate-user", false, "a bool")
	flag.Parse()

	if *genUser {
		generateUser()
	}
}
