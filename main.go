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
		Repository struct {
			Issues struct {
				PageInfo struct {
					StartCursor   githubv4.String
					HasPreviousPage bool
				}
				Edges[] struct {
					Node struct {
						Title string
						Url string
						Labels struct {
							Edges[] struct {
								Node struct {
									Name string
								}
							}
						} `graphql:"labels(first: 50)"`
					}
				}
			} `graphql:"issues(last: 1, before: $issuesCursor)"`
		} `graphql:"repository(owner: \"kubernetes\", name: \"kubernetes\")"`
	}

	variables := map[string]interface{}{
		"issuesCursor":  (*githubv4.String)(nil), // Null after argument to get first page.
	}

	client := githubv4.NewClient(httpClient)

	for i := 0; i < 4; i++ {
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("    Data:", query.Repository)
		fmt.Println("    PageInfo:", query.Repository.Issues.PageInfo.StartCursor)

		if !query.Repository.Issues.PageInfo.HasPreviousPage {
			break
		}
		variables["issuesCursor"] = githubv4.NewString(query.Repository.Issues.PageInfo.StartCursor)
		fmt.Println("Cursor now at", variables["issuesCursor"])
	}
}
