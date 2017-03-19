package main

import (
	"log"
	"strings"
	"regexp"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	bot.RemoveWebhook()

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

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
}
