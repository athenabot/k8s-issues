package main

import (
	"context"
	"github.com/shurcooL/githubv4"
	"net/http"
)

type Issue struct {
	Number   int
	Body     string
	Comments []string
	Title    string
	Url      string
	Labels   []string
	Id       string
}

func getIssuesUntilCap(ctx context.Context, httpClient *http.Client, lastIssueSeen Issue) ([]Issue, error) {
	newIssues := make([]Issue, 0)
	var cursor *githubv4.String = nil
	for {
		if len(newIssues) > 30 {
			break
		}
		var issues []Issue
		var err error
		// For some reason, := was causing cursor to be interpreted as unused.
		issues, cursor, err = getIssues(ctx, httpClient, cursor, 10)
		if err != nil {
			return nil, err
		}
		for _, issue := range issues {
			if issue.Title == lastIssueSeen.Title {
				break // TODO use actual IDs
			}
			newIssues = append(newIssues, issue)
		}
	}

	return newIssues, nil
}

func getIssues(ctx context.Context, httpClient *http.Client, cursor *githubv4.String, numIssues int) ([]Issue, *githubv4.String, error) {
	var query struct {
		Repository struct {
			Issues struct {
				PageInfo struct {
					StartCursor     githubv4.String
					HasPreviousPage bool
				}
				Edges []struct {
					Node struct {
						Id       string
						Title    string
						Url      string
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
						Number int
					}
				}
			} `graphql:"issues(last: $numIssues, before: $issuesCursor)"`
		} `graphql:"repository(owner: \"kubernetes\", name: \"kubernetes\")"`
	}

	variables := map[string]interface{}{
		"issuesCursor": cursor, // Null after argument to get first page.
		"numIssues":    (githubv4.Int)(numIssues),
	}

	client := githubv4.NewClient(httpClient)
	issues := make([]Issue, 0)
	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, nil, err
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
			Body:     issueEdge.Node.BodyText,
			Comments: comments,
			Id:       issueEdge.Node.Id,
			Labels:   labels,
			Number:   issueEdge.Node.Number,
			Title:    issueEdge.Node.Title,
			Url:      issueEdge.Node.Url,
		})
	}

	prevPage := githubv4.NewString(query.Repository.Issues.PageInfo.StartCursor)

	return issues, prevPage, nil
}
