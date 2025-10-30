FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY command_audit.sh ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o agent


FROM alpine:latest

RUN apk update && apk add --no-cache jq

RUN jq --version

COPY --from=builder /app/agent .
COPY --from=builder /app/command_audit.sh .

CMD ["./agent"]