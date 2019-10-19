package athenabot

import (
	"context"
	"net/http"
)

// An issue needs triage if:
//	* Is labeled sig/network (they've opted in to this)
//	* Does not already have label triage/unresolved
//	* Does not have "/triage resolved" comment
// TODO: use issue events, not comments. Some people edit in the UIl.
func issueNeedsTriage(issue *Issue) bool {
	// Only sig-network has opted in.
	if !issue.hasLabel("sig/network") {
		return false
	}

	// Don't double-comment.
	if issue.hasLabel("triage/unresolved") {
		return false
	}

	// Don't relabel resolved issues.
	if issue.hasCommentWithCommand("/remove-triage", "unresolved") {
		return false
	}

	return true
}

func TriageLabel(ctx context.Context, httpClient *http.Client, issue *Issue) {
	if issueNeedsTriage(issue) {
		comment := "/triage unresolved"
		comment += "\n\nComment `/remove-triage unresolved` when the issue is assessed and confirmed."
		addComment(ctx, httpClient, issue, comment)
	}
}
