# syntax=docker/dockerfile:1
#Step 1 - Build golang application
FROM golang:1.18-alpine AS builder

WORKDIR /usr/local/go/src/performance-analyzer
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
COPY ./handlers/ /usr/local/go/src/performance-analyzer/handlers
COPY ./modules/ /usr/local/go/src/performance-analyzer/modules
RUN ls -laR ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-mod=mod go build -ldflags="-w -s" -o /PerfAnalyzerApi

#Step 2 - Build a small image
FROM nginx:alpine

COPY --from=builder /PerfAnalyzerApi /PerfAnalyzerApi
EXPOSE 4300
CMD [ "/PerfAnalyzerApi" ]

