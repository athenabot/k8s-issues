package main

import (
	"context"
	"github.com/shurcooL/githubv4"
	"net/http"
	"strings"
)

type Issue struct {
	Assignees []string
	Number    int
	Body      string
	Comments  []string
	Title     string
	Url       string
	Labels    []string
	Id        string
}

func (issue *Issue) hasLabel(searchFor string) bool {
	for _, label := range issue.Labels {
		if label == searchFor {
			return true
		}
	}
	return false
}

func (issue *Issue) hasCommentWithCommand(command string, key string) bool {
	for _, comment := range issue.Comments {
		for _, line := range strings.Split(comment, "\n") {
			if strings.HasPrefix(line, command) && strings.Contains(line, key) {
				return true
			}
		}
	}
	return false
}

func getLatestIssues(ctx context.Context, httpClient *http.Client, cursor *githubv4.String, numIssues int) ([]Issue, *githubv4.String, error) {
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
						Assignees struct {
							Edges []struct {
								Node struct {
									Name string
								}
							}
						} `graphql:"assignees(first: 5)"`
						Number int
					}
				}
			} `graphql:"issues(last: $numIssues, before: $issuesCursor, states: $issueState)"`
		} `graphql:"repository(owner: \"kubernetes\", name: \"kubernetes\")"`
	}

	variables := map[string]interface{}{
		"issuesCursor": cursor, // Null after argument to get first page.
		"numIssues":    (githubv4.Int)(numIssues),
		"issueState": []githubv4.IssueState{githubv4.IssueStateOpen},
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

		assignees := make([]string, len(issueEdge.Node.Assignees.Edges))
		for index, user := range issueEdge.Node.Assignees.Edges {
			assignees[index] = user.Node.Name
		}

		issues = append(issues, Issue{
			Assignees: assignees,
			Body:      issueEdge.Node.BodyText,
			Comments:  comments,
			Id:        issueEdge.Node.Id,
			Labels:    labels,
			Number:    issueEdge.Node.Number,
			Title:     issueEdge.Node.Title,
			Url:       issueEdge.Node.Url,
		})
	}

	prevPage := githubv4.NewString(query.Repository.Issues.PageInfo.StartCursor)

	return issues, prevPage, nil
}

func getUnresolvedIssues(ctx context.Context, httpClient *http.Client, cursor *githubv4.String, numIssues int) ([]Issue, *githubv4.String, error) {
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
						Assignees struct {
							Edges []struct {
								Node struct {
									Name string
								}
							}
						} `graphql:"assignees(first: 5)"`
						Number int
					}
				}
			} `graphql:"issues(last: $numIssues, before: $issuesCursor, states:OPEN, labels:[\"sig/network\", \"triage/unresolved\"])"`
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

		assignees := make([]string, len(issueEdge.Node.Assignees.Edges))
		for index, user := range issueEdge.Node.Assignees.Edges {
			assignees[index] = user.Node.Name
		}

		issues = append(issues, Issue{
			Assignees: assignees,
			Body:      issueEdge.Node.BodyText,
			Comments:  comments,
			Id:        issueEdge.Node.Id,
			Labels:    labels,
			Number:    issueEdge.Node.Number,
			Title:     issueEdge.Node.Title,
			Url:       issueEdge.Node.Url,
		})
	}

	prevPage := githubv4.NewString(query.Repository.Issues.PageInfo.StartCursor)

	return issues, prevPage, nil
}

// Removes SIG labels from the list if they had already been added in the past.
func filterLabels(labels []string, issue Issue) []string {
	sigsCommented := make(map[string]bool)
	for _, comment := range issue.Comments {
		for _, line := range strings.Split(comment, "\n") {
			if strings.HasPrefix(line, "/sig") {
				splitSigComment := strings.Split(line, " ")
				if len(splitSigComment) > 1 {
					sigsCommented[splitSigComment[1]] = true
				}
			}
		}
	}

	var uniqueLabels []string = nil
	for _, label := range labels {
		if _, found := sigsCommented[label]; !found {
			uniqueLabels = append(uniqueLabels, label)
		}
	}
	return uniqueLabels
}
