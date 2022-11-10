# syntax=docker/dockerfile:1

#BUILD GO BACKEND
FROM golang:1.18-alpine AS go-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /usr/local/go/src/performance-analyzer
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
COPY ./handlers/ /usr/local/go/src/performance-analyzer/handlers
COPY ./modules/ /usr/local/go/src/performance-analyzer/modules
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOFLAGS=-mod=mod go build -ldflags="-w -s" -o /PerfAnalyzerApi

#BUILD WEBAPP
FROM node:latest as node-builder

WORKDIR /app
COPY ./webapp/package.json ./
COPY ./webapp/package-lock.json ./
RUN npm install --force
COPY ./webapp .
RUN npm install -g @angular/cli
RUN ng build --configuration=production-local --base-href=/performance-analyzer/ --output-path=/webapp/dist

#BUILD A SMALL FOOTPRINT IMAGE
FROM alpine:latest

COPY --from=go-builder /PerfAnalyzerApi /PerfAnalyzerApi
COPY --from=node-builder /webapp/dist/ /webapp/dist/performance-analyzer

EXPOSE 4300

ENTRYPOINT [ "PerfAnalyzerApi" ]