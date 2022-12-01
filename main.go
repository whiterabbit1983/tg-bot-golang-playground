package main

import (
	"io"
	"log"
	"net/http"
	"os"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func main() {
	api, err := tg.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))

	if err != nil {
		log.Fatalln(err)
	}

	api.Debug = true

	c := tg.NewSetMyCommands(
		tg.BotCommand{
			Command:     "get",
			Description: "get description",
		},
		tg.BotCommand{
			Command:     "put",
			Description: "put description",
		},
	)

	_, err = api.Request(c)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Document != nil {
			u, err := api.GetFileDirectURL(update.Message.Document.FileID)

			if err != nil {
				log.Println(err)
			}

			b, err := download(u)

			if err != nil {
				log.Println(err)
			}

			log.Println("+++", b)
		}

		if !update.Message.IsCommand() {
			log.Printf("got non-command message %s", update.Message.Text)

			continue
		}

		msg := tg.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "welcome!"
		case "help":
			msg.Text = "you can /get or /put"
		case "get":
			msg.Text = "get executed"
		case "put":
			msg.Text = "put executed"
		default:
			msg.Text = "unknown command"
		}

		if _, err := api.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
