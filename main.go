package main

import (
	"context"
	"fmt"
	"github.com/athenabot/k8s-issues/athenabot"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func loadSecret() string {
	secretPath := os.Getenv("GITHUB_TOKEN")
	str, err := ioutil.ReadFile(secretPath)
	if err != nil {
		log.Println(err)
	}
	return strings.Trim(string(str), "\n")
}

func main() {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: loadSecret()})
	httpClient := oauth2.NewClient(context.Background(), src)
	issues, _, err := athenabot.GetLatestIssues(context.Background(), httpClient, nil, 30)
	if err != nil {
		panic(err)
	}
	for _, issue := range issues {
		labels := athenabot.GetSigLabelsForIssue(issue)
		labels = athenabot.FilterLabels(labels, issue)
		fmt.Println(labels, issue.Url)
		err := athenabot.CommentWithSigs(context.Background(), httpClient, &issue, labels)
		fmt.Println(err)
		athenabot.TriageLabel(context.Background(), httpClient, &issue)
	}

	athenabot.SendTriageReminders(httpClient)

	//err := writeSeenIssues(context.Background(), issues)
	//fmt.Println(err)
}
