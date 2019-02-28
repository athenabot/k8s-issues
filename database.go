package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"os"
)

func gcloudProject() string {
	return os.Getenv("GCLOUD_PROJECT")
}

func writeSeenIssues(ctx context.Context, seenIssues []Issue) error {
	client, err := firestore.NewClient(ctx, gcloudProject())
	if err != nil {
		return err
	}
	defer client.Close()

	issuesCollection := client.Collection("seenIssues")
	batch := client.Batch()
	for _, issue := range seenIssues {
		issueDoc := issuesCollection.Doc(issue.Id)
		batch = batch.Set(issueDoc, map[string]interface{}{
			"title":     issue.Title,
			"number":    issue.Number,
			"repoOwner": "kubernetes",
			"repo":      "kubernetes",
			"url":       issue.Url,
		})
	}
	_, err = batch.Commit(ctx)

	return err
}
