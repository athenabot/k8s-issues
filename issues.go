package main

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"net/http"
)

type Issue struct {
	Title string
	Url   string
	Labels  []string
}

func getIssues(ctx context.Context, httpClient *http.Client, numIssues int) []Issue {
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
						Labels struct {
							Edges []struct {
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
		"issuesCursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	client := githubv4.NewClient(httpClient)

	issues := make([]Issue, 0)
	for i := 0; i < numIssues; i++ {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			fmt.Println(err)
			break
		}

		for _, issueEdge := range query.Repository.Issues.Edges {
			labels := make([]string, len(issueEdge.Node.Labels.Edges))
			for labelIndex, label := range issueEdge.Node.Labels.Edges {
				labels[labelIndex] = label.Node.Name
			}

			issues = append(issues, Issue{
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
