FROM golang:1.18-alpine as builder
RUN apk --no-cache add curl git gcc libc-dev

WORKDIR /src
COPY ./ ./

RUN go mod download
RUN go test ./...
RUN GOOS=linux go build -a -o certificates main.go

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /src/certificates .

ENTRYPOINT ["/app/certificates"]
