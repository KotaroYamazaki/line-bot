package entites

import (
	"time"

	"github.com/KotaroYamazaki/line-bot/go-mahjong-score-bot/pkg/firestore"
)

const (
	CollectionScores    firestore.Collection = "scores"
	CollectionRescource firestore.Collection = "resources"
)

type Score struct {
	UserUID     string    `firestore:"user_uid"`
	Score       int       `firestore:"score"`
	CreatedDate time.Time `firestore:"created_date"`
}

var ScoreFields = struct {
	UserUID     string
	Score       string
	CreatedDate string
}{
	UserUID:     "user_uid",
	Score:       "score",
	CreatedDate: "updated_at",
}
