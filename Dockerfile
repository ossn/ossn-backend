FROM golang:1-alpine

RUN apk add --no-cache --update git openssh jq curl gcc g++
RUN mkdir -p $GOPATH/src/github.com/ossn/ossn-backend
WORKDIR $GOPATH/src/github.com/ossn/ossn-backend
RUN mkdir ~/.ssh && \
  ssh-keyscan -t rsa github.com > ~/.ssh/known_hosts

# Install dep
RUN curl -fsSL -o /usr/local/bin/dep $(curl -s https://api.github.com/repos/golang/dep/releases/latest | jq -r ".assets[] | select(.name | test(\"dep-linux-amd64\")) |.browser_download_url") && chmod +x /usr/local/bin/dep

# Build app
COPY . .
RUN dep ensure
EXPOSE 8080
RUN go build -ldflags="-s -w" -i -o bin/main main/server.go
CMD ["bin/main"]
