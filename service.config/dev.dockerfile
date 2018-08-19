FROM golang:latest
RUN go get -u golang.org/x/lint/golint

WORKDIR /go/src/home-automation
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go", "run", "./service.config/main.go"]
