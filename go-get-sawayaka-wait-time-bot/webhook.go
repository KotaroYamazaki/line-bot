package linebot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	fs "cloud.google.com/go/firestore"
	"github.com/KotaroYamazaki/line-bot/go-get-sawayaka-wait-time-bot/entities"
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
				query, err := entities.NewQuery(message.Text)
				if err != nil {
					log.Print(err)
					continue
				}
				var itr *fs.DocumentIterator
				if query.IsALL() {
					itr = firestore.GetCollectionDocs(ctx, "sawayaka_shops")
				} else {
					itr = firestore.WhereDocumentsItr(ctx, "sawayaka_shops", entities.ShopFiels.Keywords, "array-contains", query)
				}
				var waitTimes []string

				for {
					shop, err := itr.Next()
					if err == iterator.Done || shop == nil {
						break
					}

					s := &entities.Shop{}
					if err := shop.DataTo(s); err != nil {
						log.Print(err)
						continue
					}
					if len(waitTimes) == 0 {
						jst := time.FixedZone("Asia/Tokyo", 9*60*60)
						waitTimes = append(waitTimes, fmt.Sprintf("%s 時点", s.Timestamp.In(jst).Format("2006/01/02 15:04")))
					}
					waitTimes = append(waitTimes, s.ConvertToMessage())
				}

				_, err = bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(entities.GenerateMessage(strings.Join(waitTimes, "\n")))).Do()
				if err != nil {
					log.Print(err)
				}
				continue
			}
		}
	}
	fmt.Fprint(w, "ok")
}
