package main

import (
	"context"
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

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: loadSecret()})
	httpClient := oauth2.NewClient(context.Background(), src)
	issues, _, _ := getIssues(context.Background(), httpClient, nil, 10)
	getScoresForSigs(issues[0])
	//err := writeSeenIssues(context.Background(), issues)
	//fmt.Println(err)
}
