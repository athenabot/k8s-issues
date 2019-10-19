FROM golang:1.13 as builder

RUN adduser --disabled-password --gecos '' app

RUN mkdir -p /go/src/github.com/athenabot/k8s-issues
COPY ./athenabot /go/src/github.com/athenabot/k8s-issues/athenabot
COPY main.go go.mod go.sum /go/src/github.com/athenabot/k8s-issues/

WORKDIR /go/src/github.com/athenabot/k8s-issues
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
RUN chmod +x k8s-issues



FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/athenabot/k8s-issues/k8s-issues /k8s-issues

USER app
CMD ["/k8s-issues"]
