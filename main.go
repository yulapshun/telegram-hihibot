package main

import (
	"fmt"
	"os"
	"log"
	"strings"
	"strconv"
	"regexp"
	"net/http"
	"encoding/json"
	"math"
	"math/rand"
	"gopkg.in/telegram-bot-api.v4"
)

type Configuration struct {
	Debug      bool `json:"debug"`
	UseWebhook bool `json:"useWebhook"`
	Token      string `json:"token"`
	ListenPath string `json:"listenPath"`
	ListenPort int `json:"listenPort"`
}

type Rule struct {
	Type     string `json:"type"`
	Patterns []string `json:"patterns"`
	Response string `json:"response"`
}

var bot *tgbotapi.BotAPI
var config Configuration
var ruleSet []Rule

func main() {
	// Read config
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)

	// Read rules
	file, _ = os.Open("ruleset.json")
	decoder = json.NewDecoder(file)
	err = decoder.Decode(&ruleSet)

	// Init bot
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
	updates := bot.ListenForWebhook(config.ListenPath)
	go http.ListenAndServe(":" + strconv.Itoa(config.ListenPort), nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		run(update)
	}
}

func run(update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msgText := update.Message.Text
	match := false
	resp := ""

	for _, rule := range ruleSet {
		switch rule.Type {
			case "match":
			match, resp = compareMatch(rule, msgText)
			case "contain":
			match, resp = compareContain(rule, msgText)
			case "regex":
			match, resp = compareRegex(rule, msgText)
		}
		if match {
			break
		}
	}

	if match {
		if math.Floor(rand.Float64() * 1000) == 0 {
			resp = "Bless you my child 😇"
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func compareMatch(rule Rule, msgText string) (bool, string) {
	for _,pattern := range rule.Patterns {
		if pattern == msgText {
			if strings.Contains(rule.Response, "%s") {
				return true, fmt.Sprintf(rule.Response, pattern)
			} else {
				return true, rule.Response
			}
		}
	}
	return false, ""
}

func compareContain(rule Rule, msgText string) (bool, string) {
	for _,pattern := range rule.Patterns {
		if strings.Contains(msgText, pattern) {
			if strings.Contains(rule.Response, "%s") {
				return true, fmt.Sprintf(rule.Response, pattern)
			} else {
				return true, rule.Response
			}
		}
	}
	return false, ""
}

func compareRegex(rule Rule, msgText string) (bool, string) {
	for _,pattern := range rule.Patterns {
		re := regexp.MustCompile(pattern)
		match := re.FindString(msgText)
		if match != "" {
			if strings.Contains(rule.Response, "%s") {
				return true, fmt.Sprintf(rule.Response, match)
			} else {
				return true, rule.Response
			}
		}
	}
	return false, ""
}
