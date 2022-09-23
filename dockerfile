# syntax=docker/dockerfile:1

FROM golang:1.16-alpine AS builder

WORKDIR /usr/local/go/src/golang-angular

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./
COPY ./handlers/ /usr/local/go/src/golang-angular/handlers
COPY ./dbdriver/ /usr/local/go/src/golang-angular/dbdriver
COPY ./todo/ /usr/local/go/src/golang-angular/todo

RUN ls -laR ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-mod=mod go build -ldflags="-w -s" -o /golangApp

#Step 2 - Build a small image

FROM scratch


COPY --from=builder /golangApp /golangApp

EXPOSE 4300

CMD [ "/golangApp" ]

