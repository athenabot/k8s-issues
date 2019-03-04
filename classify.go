package main

import "fmt"

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
	fmt.Println("size scaling", sizeScaling)

	for sigName, scoreData := range getScoresForSigs(issue) {
		fmt.Println("Debug", sigName, scoreData.scoreItems)
		if float64(scoreData.scoreTotal) >= scoreThreshhold*sizeScaling {
			sigs = append(sigs, sigName)
		}
	}

	return sigs
}
