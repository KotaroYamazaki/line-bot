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

					if err := firestore.AddSubCollection(ctx, entites.CollectionRescource, path, entites.CollectionScores, fsScore); err != nil {
						log.Print(err)
						continue
					}

					_, err = bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s%d 記録しました:D", getSign(score), score))).Do()
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
				case "total", "Total":
					itr := firestore.WhereDocumentsItr(ctx, entites.CollectionRescource, entites.CollectionScores, path, entites.ScoreFields.UserUID, "==", e.Source.UserID)
					sum := 0
					cnt := 0
					avg := float64(0)
					for {
						score, err := itr.Next()
						if err == iterator.Done {
							break
						}
						s := &entites.Score{}
						err = score.DataTo(s)
						if err != nil {
							log.Print(err)
							continue
						}
						sum += s.Score
						cnt++
					}
					if cnt != 0 {
						avg = float64(sum) / float64(cnt)
					}

					_, err = bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("=============\n▼総局数\n%d回\n▼累計スコア\n%s%d\n▼平均スコア\n%s%.2f\n==============", cnt, getSign(sum), sum, getSignF(avg), avg))).Do()
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

func getSign(score int) string {
	if score > 0 {
		return "+"
	}
	return ""
}

func getSignF(score float64) string {
	if score > 0 {
		return "+"
	}
	return ""
}
