package main

import (
	"fmt"
	"strings"
)

var scoreThreshhold float64 = 5

// Returns the list of sigs to classify an issue as.
func getSigLabelsForIssue(issue Issue) []string {
	var sigs []string = nil

	var sizeFactor float64 = 400
	issueSize := float64(len(issue.Title) + len(issue.Body))
	sizeScaling := 0.75 * issueSize / sizeFactor
	if sizeScaling < 1 { // Don't weirdly scale tiny issues
		sizeScaling = 1
	}

	// Quick hack to be less inclined to comment on issues that have already been sorted.
	for _, label := range issue.Labels {
		if strings.HasPrefix(label, "sig/") {
			sizeScaling *= 1.3
			break
		}
	}

	fmt.Println("size scaling", sizeScaling)

	for sigName, scoreData := range getScoresForSigs(issue) {
		fmt.Println("Debug", sigName, scoreData.scoreItems)
		if float64(scoreData.scoreTotal) >= scoreThreshhold*sizeScaling {
			sigs = append(sigs, sigName)
		}
	}

	// Quick hack to catch non-comment sig additions (EG k8s CI bot)
	for _, label := range issue.Labels {
		if strings.HasPrefix(label, "sig/") {
			sigName := strings.Split(label,"sig/")[1]
			for i, pickedLabel := range sigs {
				if sigName == pickedLabel {
					fmt.Println("already labeled", sigName)
					sigs = append(sigs[:i], sigs[i+1:]...)
				}
			}
		}
	}

	return sigs
}
