package main

import (
	"testing"
)

func TestScoreForSig(t *testing.T) {
	sources := []textSource{
		{
			name:       "title",
			multiplier: 5,
			content:    "Error while booting kube-proxy",
		},
		{
			name:       "body",
			multiplier: 1,
			content:    "kube-proxy fails to start, which causes services to fail",
		},
	}
	scoreLineItems := scoreForSig(sigNetwork, sources)

	expectedLineItems := 3
	if len(scoreLineItems) != expectedLineItems {
		t.Errorf("Expected exactly %d score items: %v", expectedLineItems, scoreLineItems)
	}
	expectedScore := 19
	if scoreLineItems[0].points+scoreLineItems[1].points+scoreLineItems[2].points != expectedScore {
		t.Errorf("Expected total score to be %v: %v", expectedScore, scoreLineItems)
	}
}

func TestGetScoresForSigs(t *testing.T) {
	testIssue := Issue{
		Title: "Test issue mentioning services",
		Body:  "This text is nothing of consequence.",
	}
	sigScores := getScoresForSigs(testIssue)

	sigName := "network"
	expectedTotal := 3
	if sigScores[sigName].scoreTotal != expectedTotal {
		t.Errorf("Expected sig %s score total to be %d, got: %v", sigName, expectedTotal, sigScores[sigName])
	}
}
