package main

import (
	"context"
	"github.com/shurcooL/githubv4"
	"net/http"
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
	comment += "\nThese SIGs are my best guesses for this issue. Please comment `/remove-sig <name>` if I am incorrect about one." +
		"\n🤖 I am an (alpha) bot run by @vllry. 👩‍🔬"
	return addComment(ctx, httpClient, issueId, comment)
}

func addComment(ctx context.Context, httpClient *http.Client, issueId string, comment string) error {
	client := githubv4.NewClient(httpClient)

	input := githubv4.AddCommentInput{
		SubjectID: issueId,
		Body:      githubv4.String(comment),
	}
	return client.Mutate(ctx, &m, input, nil)
}
