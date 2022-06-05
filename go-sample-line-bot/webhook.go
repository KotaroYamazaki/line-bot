package linebot

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	client, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	events, err := client.ParseRequest(r)
	if err != nil {
		http.Error(w, "Error parse request", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	for _, e := range events {
		switch e.Type {
		case linebot.EventTypeMessage:
			message := linebot.NewTextMessage("Test")
			_, err := client.ReplyMessage(e.ReplyToken, message).Do()
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
	fmt.Fprint(w, "ok")
}
