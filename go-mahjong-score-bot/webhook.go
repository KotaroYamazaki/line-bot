package linebot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/KotaroYamazaki/line-bot/go-mahjong-score-bot/entites"
	"github.com/KotaroYamazaki/line-bot/go-mahjong-score-bot/pkg/firestore"
	"github.com/line/line-bot-sdk-go/linebot"
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
				text := strings.TrimSpace(message.Text)
				var path string
				switch e.Source.Type {
				case linebot.EventSourceTypeUser:
					path = e.Source.UserID
				case linebot.EventSourceTypeGroup:
					path = e.Source.GroupID
				default:
					continue
				}
				if strings.HasPrefix(text, "+") || strings.HasPrefix(text, "-") {
					sign := text[0]
					text = text[1:]
					score, err := strconv.Atoi(text)
					if err != nil {
						log.Print(err)
						continue
					}
					if sign == '-' {
						score = -score
					}
					fsScore := &entites.Score{
						UserUID:     e.Source.UserID,
						Score:       score,
						CreatedDate: time.Now(),
					}

					if err := firestore.Set(ctx, entites.CollectionScores, path, fsScore); err != nil {
						log.Print(err)
						continue
					}

					_, err = bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%d 記録しました", score))).Do()
					if err != nil {
						log.Print(err)
					}
					continue
				}

				switch text {
				case "undo", "Undo":
					_, err := bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage("undo")).Do()
					if err != nil {
						log.Print(err)
					}
					continue
				case "total":
					score, err := firestore.Get(ctx, entites.CollectionScores, path, &entites.Score{})
					if err != nil {
						log.Print(err)
						continue
					}
					s, ok := score.(*entites.Score)
					if !ok {
						log.Print("Error type assertion")
						continue
					}
					_, err = bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s-(%d)-%s", s.UserUID, s.Score, s.CreatedDate.Format("2005/01/02")))).Do()
					if err != nil {
						log.Print(err)
					}
					continue
				}
			}
		}
	}
	fmt.Fprint(w, "ok")
}
