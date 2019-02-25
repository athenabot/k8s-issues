package main

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"math"
	"net/http"
)

type Issue struct {
	Body string
	Comments []string
	Title string
	Url   string
	Labels  []string
}

func getIssues(ctx context.Context, httpClient *http.Client, minIssuesToReturn int) []Issue {
	var query struct {
		Repository struct {
			Issues struct {
				PageInfo struct {
					StartCursor     githubv4.String
					HasPreviousPage bool
				}
				Edges []struct {
					Node struct {
						Title  string
						Url    string
						BodyText string
						Comments struct {
							Nodes []struct {
								Body string
							}
						} `graphql:"comments(first: 50)"`
						Labels struct {
							Edges []struct {
								Node struct {
									Name string
								}
							}
						} `graphql:"labels(first: 50)"`
					}
				}
			} `graphql:"issues(last: 10, before: $issuesCursor)"`
		} `graphql:"repository(owner: \"kubernetes\", name: \"kubernetes\")"`
	}

	variables := map[string]interface{}{
		"issuesCursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	client := githubv4.NewClient(httpClient)

	numLoops := int(math.Ceil(float64(minIssuesToReturn) / 10.0))
	issues := make([]Issue, 0)
	for i := 0; i < numLoops; i++ {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			fmt.Println(err)
			break
		}

		for _, issueEdge := range query.Repository.Issues.Edges {
			comments := make([]string, len(issueEdge.Node.Comments.Nodes))
			for index, comment := range issueEdge.Node.Comments.Nodes {
				comments[index] = comment.Body
			}

			labels := make([]string, len(issueEdge.Node.Labels.Edges))
			for index, label := range issueEdge.Node.Labels.Edges {
				labels[index] = label.Node.Name
			}

			issues = append(issues, Issue{
				Body: issueEdge.Node.BodyText,
				Comments: comments,
				Labels: labels,
				Title: issueEdge.Node.Title,
				Url: issueEdge.Node.Url,
			})
		}

		if !query.Repository.Issues.PageInfo.HasPreviousPage {
			break
		}
		variables["issuesCursor"] = githubv4.NewString(query.Repository.Issues.PageInfo.StartCursor)
	}

	return issues
}
