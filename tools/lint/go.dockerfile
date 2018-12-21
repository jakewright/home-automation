FROM golang:1.11
RUN go get -u golang.org/x/lint/golint
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/home-automation

COPY tools/lint/go_fmt.sh /
RUN chmod +x /go_fmt.sh

# Add lock files and install dependencies.
# Hopefully this won't change much and will be cached.
COPY Gopkg.* ./
RUN dep ensure -vendor-only

# Add everything else. This will change a lot.
# COPY . .

CMD ["/go_fmt.sh"]