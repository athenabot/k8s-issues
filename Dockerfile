FROM golang:1.13 as builder

RUN mkdir -p /go/src/github.com/athenabot/k8s-issues
COPY ./athenabot /go/src/github.com/athenabot/k8s-issues/athenabot
COPY main.go go.mod go.sum /go/src/github.com/athenabot/k8s-issues/

WORKDIR /go/src/github.com/athenabot/k8s-issues
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build



FROM alpine

COPY --from=builder /go/src/github.com/athenabot/k8s-issues/k8s-issues /k8s-issues
RUN chmod +x /k8s-issues
CMD ["/k8s-issues"]
