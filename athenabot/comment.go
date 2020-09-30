package athenabot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shurcooL/githubv4"
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

func CommentWithSigs(ctx context.Context, httpClient *http.Client, issue *Issue, sigs []string) error {
	if len(sigs) == 0 {
		return nil
	}

	comment := ""
	for _, sigName := range sigs {
		comment += "/sig " + sigName + "\n"
	}
	comment += "\nThese SIGs are my best guesses for this issue. Please comment `/remove-sig <name>` if I am incorrect about one."

	return addComment(ctx, httpClient, issue, comment)
}

func addComment(ctx context.Context, httpClient *http.Client, issue *Issue, comment string) error {
	fmt.Printf("\nComment on issue %v: %v", issue.Url, comment)

	client := githubv4.NewClient(httpClient)
	input := githubv4.AddCommentInput{
		SubjectID: issue.Id,
		Body:      githubv4.String(comment),
	}
	return client.Mutate(ctx, &m, input, nil)
}
