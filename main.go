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

var Query struct {
	Repository struct {
		Issues struct {
			Edges struct {
				Node struct {
					Title string
					Url string
				}
			}
		} `graphql:"issues(last: 5)"`
	} `graphql:"repository(owner: \"kubernetes\", name: \"kubernetes\")"`
}



//	repository(owner: "kubernetes", name: "kubernetes") {
//	issues(last: 5, states: OPEN) {
//	edges {
//	node {
//	title
//	url
//	labels(first: 50) {
//	edges {
//	node {
//	name

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: loadSecret()},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	err := client.Query(context.Background(), &Query, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("    Data:", Query.Repository)
}
