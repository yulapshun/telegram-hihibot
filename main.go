package main

import (
	"os"
	"log"
	"strings"
	"strconv"
	"regexp"
	"net/http"
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
)

type Configuration struct {
	Debug      bool `json:"debug"`
	UseWebhook bool `json:"useWebhook"`
	Token      string `json:"token"`
	ListenPath string `json:"listenPath"`
	ListenPort int `json:"listenPort"`
}

var bot *tgbotapi.BotAPI
var config Configuration

func main() {

	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)

	_bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}
	bot = _bot

	bot.Debug = config.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	if config.UseWebhook {
		runWebhook()
	} else {
		runPoll()
	}
}

func runPoll() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	bot.RemoveWebhook()
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		run(update)
	}
}

func runWebhook() {
	updates := bot.ListenForWebhook(config.ListenPath + bot.Token)
	go http.ListenAndServe(":" + strconv.Itoa(config.ListenPort), nil)

	for update := range updates {
		run(update)
	}
}

func run(update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	resp := ""

	if strings.Contains(update.Message.Text, "老母") {
		resp = "唔準提老母\U0001F621"
	} else if strings.Contains(update.Message.Text, "屌") {
		resp = "唔好屌\U0001F621"
	} else {
		match, _ := regexp.MatchString("[F|f]uck(ing|ed)?", update.Message.Text)
		if match {
			resp = "No fucking\U0001F621"
		}
	}

	if resp != "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
