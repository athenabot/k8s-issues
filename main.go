package main

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
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
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: loadSecret()},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	var query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
	}

	client := githubv4.NewClient(httpClient)
	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		// Handle error.
	}
	fmt.Println("    Login:", query.Viewer.Login)
	fmt.Println("CreatedAt:", query.Viewer.CreatedAt)
}
