package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

func loadSecret() string {
	str, err := ioutil.ReadFile("secret.txt")
	if err != nil {
		log.Println(err)
	}
	return string(str)
}

func main() {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: loadSecret()})
	httpClient := oauth2.NewClient(context.Background(), src)
	issues, _, _ := getIssues(context.Background(), httpClient, nil, 30)
	for _, issue := range issues {
		labels := getSigLabelsForIssue(issue)
		labels = filterLabels(labels, issue)
		fmt.Println(labels, issue.Url)
		err := commentWithSigs(context.Background(), httpClient, issue.Id, labels)
		fmt.Println(err)
		triageLabel(context.Background(), httpClient, &issue)
	}
	//err := writeSeenIssues(context.Background(), issues)
	//fmt.Println(err)
}
