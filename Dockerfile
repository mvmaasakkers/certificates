FROM golang:1.12-alpine as builder
RUN apk --no-cache add curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/mvmaasakkers/certificates
COPY . /go/src/github.com/mvmaasakkers/certificates
RUN dep ensure
RUN go test ./...
RUN GOOS=linux go build -a -o certificates main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /go/src/github.com/mvmaasakkers/certificates/certificates .

ENTRYPOINT ["/app/certificates"]
