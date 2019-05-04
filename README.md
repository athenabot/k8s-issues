# Athenabot k8s-issues

This is a simple (WIP) bot, to help with issue management in the Kubernetes org.

# What it does

* Tries to guess SIG labels for issues in kubernetes/kubernetes that have no SIG labels.
* Labels new SIG-network tickets with `triage/unresolved`.
    * Please file an issue if you would like this extended to your SIG!
* Posts periodic reminder comments on `triage/unresolved` that have sat idle with an assignee.

# Contributing

Suggestions and contributions are highly welcome. They can be:

* New functionality (including features beyond classification).
* Technical improvements.
* Sig/keyword changes.
* And more!

# Running

Currently the bot runs in a hacky batch mode, this may change down the line to running as a service.

Install a GitHub credential to secret.txt.

`go install && $GOPATH/bin/k8s-issues`
