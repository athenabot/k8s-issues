# Athenabot k8s-issues

This is a simple (WIP) bot, for classifying issues in kubernetes/kubernetes.

# How it works

Athenabot reads recent GitHub issue data, and applies basic keyword scoring. If it finds suitable matches, it comments for the k8s ci-bot.

# Contributing

Suggestions and contributions are highly welcome. They can be:

* New functionality (including features beyond classification).
* Technical improvements.
* Sig/keyword changes.
* And more!

# Running

Currently the bot runs in a hacky batch mode, this will change down the line to running as a service.

Install a GitHub credential to secret.txt.

`go install && $GOPATH/bin/k8s-issues`
