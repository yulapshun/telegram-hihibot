package main

import (
	"log"
	"strings"
	"regexp"
	"net/http"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://5a52b38d.ngrok.io/"+bot.Token, "cert.pem"))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8080", "cert.pem", "key.pem", nil)

	for update := range updates {
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
