package main

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"net/http"
	"strings"
)

var m struct {
	AddComment struct {
		CommentEdge struct {
			Node struct {
				Body githubv4.String
			}
		}
		Subject struct {
			ID githubv4.ID
		}
	} `graphql:"addComment(input: $input)"`
}

func commentWithSigs(ctx context.Context, httpClient *http.Client, issueId string, sigs []string) error {
	if len(sigs) == 0 {
		return nil
	}

	comment := ""
	for _, sigName := range sigs {
		comment += "/sig " + sigName + "\n"
	}
	comment += "\nThese SIGs are my best guesses for this issue. Please comment `/remove-sig <name>` if I am incorrect about one."

	return addComment(ctx, httpClient, issueId, comment)
}

func commentTriageReminder(ctx context.Context, httpClient *http.Client, issue *Issue, assignees []string) {
	comment := strings.Join(assignees, " ") + "\n"
	comment += "If this issue has been triaged, please comment `/remove-triage unresolved`."
	comment += "\n\nMeta:\n/athenabot mark-triage-reminder"
	fmt.Println(comment)
	//addComment(ctx, httpClient, issue.Id, comment)
}

func addComment(ctx context.Context, httpClient *http.Client, issueId string, comment string) error {
	signature := "\n\n🤖 I am a bot run by @vllry. 👩‍🔬"
	client := githubv4.NewClient(httpClient)

	input := githubv4.AddCommentInput{
		SubjectID: issueId,
		Body:      githubv4.String(comment + signature),
	}
	return client.Mutate(ctx, &m, input, nil)
}
