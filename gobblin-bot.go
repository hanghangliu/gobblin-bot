package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-joe/cron"
	joehttp "github.com/go-joe/http-server"
	"github.com/go-joe/joe"
	"github.com/go-joe/slack-adapter"
)

type GobblinBot struct {
	*joe.Bot
	Config *Config
}

type DailyDigestEvent struct {
}

type Config struct {
	SlackAppToken     string
	SlackBotUserToken string
	From              string
	To                string
	SendgridToken     string
	Port              string
}

func NewConfig() (*Config, error) {
	config := &Config{
		SlackAppToken:     os.Getenv("SLACK_APP_TOKEN"),
		SlackBotUserToken: os.Getenv("SLACK_BOT_USER_TOKEN"),
		From:              os.Getenv("FROM"),
		To:                os.Getenv("TO"),
		SendgridToken:     os.Getenv("SENDGRID_TOKEN"),
		Port:              os.Getenv("PORT"),
	}
	fmt.Println(config)
	if config.Port == "" {
		config.Port = "80"
	}
	return config, nil
}

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	modules := []joe.Module{
		joehttp.Server(":" + config.Port),
		// Schedule the daily digest cron job at 2:00:00 AM (UTC), which is California time 7 PM
		cron.ScheduleEvent("0 0 2 * * *", DailyDigestEvent{}),
	}
	if config.SlackAppToken != "" && config.SlackBotUserToken != "" {
		modules = append(modules, slack.Adapter(config.SlackBotUserToken))
	}

	b := &GobblinBot{
		Bot:    joe.New("Gobblin-bot", modules...),
		Config: config,
	}

	// Register event handlers
	b.Brain.RegisterHandler(b.HandleDailyDigestEvent)
	b.Brain.RegisterHandler(b.HandleHTTP)
	b.Respond("daily-digest", b.DailyDigest)
	b.Respond("config", b.PrintConfig)
	b.Respond("ping", Pong)
	b.Respond("time", Time)

	b.Say("daily-digest", "Gobblin bot is starting..")
	err = b.Run()
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
}

func (b *GobblinBot) HandleDailyDigestEvent(evt DailyDigestEvent) {
	responseMsg := RunDailyDigest(b.Config)
	b.Say("daily-digest", responseMsg)
}

func (b *GobblinBot) DailyDigest(msg joe.Message) error {
	responseMsg := RunDailyDigest(b.Config)
	msg.Respond(responseMsg)
	return nil
}

func (b *GobblinBot) PrintConfig(msg joe.Message) error {
	fmt.Println("printconfig")
	configMsg := fmt.Sprintf("From: `%s`\n", b.Config.From)
	configMsg += fmt.Sprintf("To: `%s`\n", b.Config.To)
	configMsg += fmt.Sprintf("SlackAppToken: `%s`\n", b.Config.SlackAppToken)
	configMsg += fmt.Sprintf("SlackBotUserToken: `%s`\n", b.Config.SlackBotUserToken)
	configMsg += fmt.Sprintf("SendgridToken: `%s`", b.Config.SendgridToken)
	msg.Respond(configMsg)
	return nil
}

func (b *GobblinBot) HandleHTTP(c context.Context, r joehttp.RequestEvent) {
	if r.URL.Path == "/" {
		fmt.Println("Gobblin bot is running..")
	}
}

func Time(msg joe.Message) error {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Now()
	timeMsg := fmt.Sprintf("Machine local time: `%s`\n", fmt.Sprint(t))
	timeMsg += fmt.Sprintf("Machine local time (in PDT): `%s`", fmt.Sprint(t.In(loc)))
	msg.Respond(timeMsg)
	return nil
}

func Pong(msg joe.Message) error {
	msg.Respond("PONG")
	return nil
}
