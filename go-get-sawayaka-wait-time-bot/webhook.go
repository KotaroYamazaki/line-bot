package linebot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/KotaroYamazaki/line-bot/go-get-sawayaka-wait-time-bot/entites"
	"github.com/KotaroYamazaki/line-bot/go-get-sawayaka-wait-time-bot/pkg/firestore"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/api/iterator"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		http.Error(w, "Error init line bot", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	firestore, err := firestore.New(ctx)
	if err != nil {
		http.Error(w, "Error init firestore", http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		http.Error(w, "Error parse request", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	for _, e := range events {
		switch e.Type {
		case linebot.EventTypeMessage:
			switch message := e.Message.(type) {
			case *linebot.TextMessage:
				query := strings.TrimSpace(message.Text)
				itr := firestore.WhereDocumentsItr(ctx, "sawayaka_shops", entites.ShopFiels.Keywords, "array-contains", query)
				var waitTimes []string

				for {
					shop, err := itr.Next()
					if err == iterator.Done || shop == nil {
						break
					}

					s := &entites.Shop{}
					if err := shop.DataTo(s); err != nil {
						log.Print(err)
						continue
					}
					waitTimes = append(waitTimes, s.ConvertText())
				}

				_, err = bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(entites.GenerateMessage(strings.Join(waitTimes, "\n")))).Do()
				if err != nil {
					log.Print(err)
				}
				continue
			}
		}
	}
	fmt.Fprint(w, "ok")
}
