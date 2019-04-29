FROM golang:1.12-alpine as builder
RUN apk --no-cache add curl git gcc libc-dev

WORKDIR /src
COPY ./ ./

RUN go mod download
RUN go test ./...
RUN GOOS=linux go build -a -o certificates main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /go/certificates .

ENTRYPOINT ["/app/certificates"]
