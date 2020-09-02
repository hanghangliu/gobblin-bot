# gobblin-bot
Send daily digest email for Apache Gobblin slack workplace


# Build and run locally
```
# Install dependencies
$ go get -d ./...

# Running gobblin-bot locally
$ go run gobblin-bot.go digest.go

# Running with all configs
$ SLACK_APP_TOKEN=<app_token> SLACK_BOT_USER_TOKEN=<bot_token> FROM=xxx@domain.com TO=xxx@xx.com SENDGRID_TOKEN=<sendgrid_token> PORT=5005 go run gobblin-bot.go digest.go
```

# Build docker image
```
# Install gox (cross compilation tool)
$ go get github.com/mitchellh/gox

# Add GOROOT/GOPATH to the environmental path.
$ export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
# If not working, try this
$ export GOPATH="$HOME/go"
$ PATH="$GOPATH/bin:$PATH"

# Run makefile script to build docker image
$ make build

# Run locally
$ docker run -a stdin -a stdout -i -t build/gobblin-bot

# Run with all configs
$ docker run -p 5005:80 -e PORT=80 -e SLACK_APP_TOKEN=<slack_token> -e SLACK_BOT_USER_TOKEN=<bot_token> -e FROM=xxx@your-domain.com -e TO=xxx@xx.com -e SENDGRID_TOKEN=<sendgrid_token> build/gobblin-bot

# Publish docker image to docker hub
$ make push
```
