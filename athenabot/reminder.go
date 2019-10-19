package athenabot

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func SendTriageReminders(httpClient *http.Client) {
	issues, err := getUnresolvedIssues(context.Background(), httpClient)
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
			if strings.Contains(comment.Body, "issue has been triaged") {
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
			assignedDuration := now.Sub(issue.LastAssignedTime)
			if assignedDuration >= duration {
				filteredIssues = append(filteredIssues, issue)
			}
		}
	}

	return filteredIssues
}

func commentTriageReminder(ctx context.Context, httpClient *http.Client, issue *Issue, assignees []string) {
	comment := strings.Join(assignees, " ") + "\n"
	comment += "If this issue has been triaged, please comment `/remove-triage unresolved`." +
		"\n\nIf you aren't able to handle this issue, consider unassigning yourself and/or adding the `help-wanted` label."
	_ = addComment(ctx, httpClient, issue, comment)
}
