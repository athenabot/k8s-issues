package main

import (
	"testing"
	"time"
)

func aInB(A []Issue, B []Issue) bool {
	bMap := make(map[string]Issue)
	for _, b := range B {
		bMap[b.Id] = b
	}

	for _, a := range A {
		_, found := bMap[a.Id]
		if !found {
			return false
		}
	}

	return true
}

func errorIfSetsNotEqual(t *testing.T, actual []Issue, expected []Issue) {
	if !aInB(actual, expected) || !aInB(expected, actual) {
		t.Errorf("got: %v\nexpected: %v", actual, expected)
	}
}

func TestFilterIssuesAssignedLongerThan(t *testing.T) {
	issues := []Issue{
		{
			Id:               "no1",
			Assignees:        []string{"someone"},
			LastAssignedTime: time.Now(),
		},
		{
			Id:               "yes1",
			Assignees:        []string{"someone"},
			LastAssignedTime: time.Now().Add(-2 * time.Hour),
		},
		{
			Id:               "yes2",
			Assignees:        []string{"someone"},
			LastAssignedTime: time.Now().Add(-time.Hour - time.Second),
		},
	}

	expect := []Issue{
		{
			Id: "yes1",
		},
		{
			Id: "yes2",
		},
	}

	filtered := filterIssuesAssignedLongerThan(issues, time.Hour)
	errorIfSetsNotEqual(t, filtered, expect)
}

func TestFilterToIssuesNeedingReminders(t *testing.T) {
	issues := []Issue{
		{
			Id:               "yes1",
			Assignees:        []string{"someone"},
			LastAssignedTime: time.Now().Add(-time.Hour * 24 * 8),
		},
		{
			Id:               "no1",
			Assignees:        []string{"someone"},
			LastAssignedTime: time.Now().Add(-time.Hour * 24 * 6),
		},
		{
			Id:        "no2",
			Assignees: []string{"someone"},
			Comments: []IssueComment{
				{
					Body:        "If this issue has been triaged,",
					CreatedTime: time.Now().Add(-time.Hour * 24 * 25),
					User:        "athenabot",
				},
			},
			LastAssignedTime: time.Now().Add(-time.Hour * 24 * 20),
		},
		{
			Id:        "yes2",
			Assignees: []string{"someone"},
			Comments: []IssueComment{
				{
					Body:        "If this issue has been triaged,",
					CreatedTime: time.Now().Add(-time.Hour * 24 * 31),
					User:        "athenabot",
				},
			},
			LastAssignedTime: time.Now().Add(-time.Hour * 24 * 50),
		},
	}

	expect := []Issue{
		{
			Id: "yes1",
		},
		{
			Id: "yes2",
		},
	}

	filtered := filterToIssuesNeedingReminders(issues)
	errorIfSetsNotEqual(t, filtered, expect)
}
