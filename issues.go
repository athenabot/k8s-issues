package main

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"net/http"
	"strings"
	"time"
)

type Issue struct {
	Assignees        []string
	LastAssignedTime time.Time
	Number           int
	Body             string
	Comments         []IssueComment
	Title            string
	Url              string
	Labels           []string
	Id               string
}

type IssueComment struct {
	Body        string
	CreatedTime time.Time
	User        string
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
		for _, line := range strings.Split(comment.Body, "\n") {
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
		"issueState":   []githubv4.IssueState{githubv4.IssueStateOpen},
	}

	client := githubv4.NewClient(httpClient)
	issues := make([]Issue, 0)
	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, nil, err
	}

	for _, issueEdge := range query.Repository.Issues.Edges {

		// TODO clean up this use of comments. Use label history instead.
		comments := make([]IssueComment, len(issueEdge.Node.Comments.Nodes))
		for index, ghComment := range issueEdge.Node.Comments.Nodes {
			comments[index] = IssueComment{
				Body: ghComment.Body,
			}
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

func getUnresolvedIssues(ctx context.Context, httpClient *http.Client) ([]Issue, error) {
	issues := make([]Issue, 0)

	firstLoop := true
	blankCursor := githubv4.NewString("") // Use for comparison. Cannot use len().
	var cursor *githubv4.String = nil // nil -> "get first page"
	for firstLoop || (cursor != nil && cursor != blankCursor) {
		firstLoop = false

		// Hard to tell if the cursor is "done" or not. Use the hasCursor indicator instead.
		issueBatch, hasCursor, newCursor, err := getUnresolvedIssuesBatch(ctx, httpClient, cursor, 25)
		cursor = newCursor // Goland complaints when assigning cursor directly.
		if err != nil {
			return nil, err // TODO retry, EG for rate limits.
		}

		for _, issue := range issueBatch {
			issues = append(issues, issue)
		}
		fmt.Println("Got issues: ", len(issueBatch))

		if !hasCursor {
			break
		}
	}

	return issues, nil
}

func getUnresolvedIssuesBatch(ctx context.Context, httpClient *http.Client, cursor *githubv4.String, numIssues int) ([]Issue, bool, *githubv4.String, error) {
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
						Comments struct {
							Nodes []struct {
								Author struct {
									Login string
								}
								Body      string
								CreatedAt string
							}
						} `graphql:"comments(last: 50)"`
						Labels struct {
							Edges []struct {
								Node struct {
									Name string
								}
							}
						} `graphql:"labels(first: 30)"`
						Assignees struct {
							Edges []struct {
								Node struct {
									Login string
								}
							}
						} `graphql:"assignees(first: 10)"`
						TimelineItems struct {
							Nodes []struct {
								AssignedEvent struct {
									CreatedAt string
								} `graphql:"... on AssignedEvent"`
								LabeledEvent struct {
									CreatedAt string
								} `graphql:"... on AssignedEvent"`
							}
						} `graphql:"timeline(last: 50)"`
						Number int
						Url string
					}
				}
			} `graphql:"issues(last: $numIssues, before: $issuesCursor, states:OPEN, labels:[\"triage/unresolved\"])"`
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
		return nil, false, nil, err
	}

	for _, issueEdge := range query.Repository.Issues.Edges {

		comments := make([]IssueComment, len(issueEdge.Node.Comments.Nodes))
		for index, ghComment := range issueEdge.Node.Comments.Nodes {
			commentTime, err := ghStringToTime(ghComment.CreatedAt)
			if err != nil {
				panic(err)
			}
			comments[index] = IssueComment{
				Body:        ghComment.Body,
				CreatedTime: commentTime,
				User:        ghComment.Author.Login,
			}
		}

		labels := make([]string, len(issueEdge.Node.Labels.Edges))
		for index, label := range issueEdge.Node.Labels.Edges {
			labels[index] = label.Node.Name
		}

		assignees := make([]string, len(issueEdge.Node.Assignees.Edges))
		for index, user := range issueEdge.Node.Assignees.Edges {
			if user.Node.Login != "" {
				assignees[index] = "@" + user.Node.Login
			}
		}

		// TODO this is messy, clean up
		// TODO double check order assumptions
		lastAssignedTimeStr := ""
		for _, timelineItem := range issueEdge.Node.TimelineItems.Nodes {
			if timelineItem.AssignedEvent.CreatedAt != "" {
				lastAssignedTimeStr = timelineItem.AssignedEvent.CreatedAt
				break // Get most recent only
			}
		}
		lastAssignedTime, err := ghStringToTime(lastAssignedTimeStr)
		if err != nil && lastAssignedTimeStr != "" {
			fmt.Println(lastAssignedTimeStr)
			panic(err)
		}

		issues = append(issues, Issue{
			Assignees:        assignees,
			LastAssignedTime: lastAssignedTime,
			Comments:         comments,
			Id:               issueEdge.Node.Id,
			Labels:           labels,
			Number:           issueEdge.Node.Number,
			Title:            issueEdge.Node.Title,
			Url:              issueEdge.Node.Url,
		})

	}

	prevPage := githubv4.NewString(query.Repository.Issues.PageInfo.StartCursor)

	return issues, query.Repository.Issues.PageInfo.HasPreviousPage, prevPage, nil
}

// TODO use label history
// Removes SIG labels from the list if they had already been added in the past.
func filterLabels(labels []string, issue Issue) []string {
	sigsCommented := make(map[string]bool)
	for _, comment := range issue.Comments {
		for _, line := range strings.Split(comment.Body, "\n") {
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

func ghStringToTime(timestr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestr)
}
