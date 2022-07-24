package entites

import (
	"fmt"
	"time"
)

func GenerateMessage(text string) string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst).Format("2006/01/02 15:04")
	return fmt.Sprintf("%s \nただいまの待ち時間 \n%s", now, text)
}
