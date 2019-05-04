package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func sendReminders(httpClient *http.Client) {
	issues, _, err := getUnresolvedIssues(context.Background(), httpClient, nil, 30)
	if err != nil {
		panic(err)
	}

	for _, issue := range filterToIssuesNeedingReminders(issues) {
		commentTriageReminder(context.Background(), httpClient, &issue, issue.Assignees)
	}
}

func filterToIssuesNeedingReminders(issues []Issue) []Issue {
	needReminder := make([]Issue, 0)

	for _, issue := range filterIssuesAssignedLongerThan(issues, 7*24*time.Hour) {
		timeSinceReminder := timeSinceTriageReminder(issue)

		if timeSinceReminder == nil || *timeSinceReminder > 30*24*time.Hour {
			needReminder = append(needReminder, issue)
		}
	}

	return needReminder
}

func timeSinceTriageReminder(issue Issue) *time.Duration {
	// Comments are in chronological order.
	for i := len(issue.Comments) - 1; i >= 0; i-- {
		comment := issue.Comments[i]
		if comment.User == "athenabot" {
			if strings.Contains(comment.Body, "mark-triage") {
				duration := time.Now().Sub(comment.CreatedTime)
				return &duration
			}
		}
	}
	return nil
}

func filterIssuesAssignedLongerThan(issues []Issue, duration time.Duration) []Issue {
	filteredIssues := make([]Issue, 0)

	now := time.Now()
	for _, issue := range issues {
		if len(issue.Assignees) != 0 {
			fmt.Println(issue.Title, strings.Join(issue.Assignees, ", "))

			assignedDuration := now.Sub(issue.LastAssignedTime)
			fmt.Println("issue has been assigned for: ", assignedDuration)
			if assignedDuration >= duration {
				filteredIssues = append(filteredIssues, issue)
			}
		}
	}

	return filteredIssues
}

func commentTriageReminder(ctx context.Context, httpClient *http.Client, issue *Issue, assignees []string) {
	comment := strings.Join(assignees, " ") + "\n"
	comment += "If this issue has been triaged, please comment `/remove-triage unresolved`."
	comment += "\n\nMeta:\n/athenabot mark-triage-reminder"
	fmt.Println(comment)
	//addComment(ctx, httpClient, issue.Id, comment)
}
