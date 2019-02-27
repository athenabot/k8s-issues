package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

func loadSecret() string {
	str, err := ioutil.ReadFile("secret.txt")
	if err != nil {
		log.Println(err)
	}
	return string(str)
}

func main() {
	var mostRecent Issue

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: loadSecret()},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	issues, next, err := getIssues(context.Background(), httpClient, nil,5)
	mostRecent = issues[0]
	fmt.Println(issues)
	fmt.Println(next)
	fmt.Println(err)
}
