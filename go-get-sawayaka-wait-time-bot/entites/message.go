package entites

import (
	"fmt"
	"time"
)

func GenerateMessage(text string) string {
	now := time.Now().Local().Format("2006/01/02 15:04")
	return fmt.Sprintf("%s \nただいまの待ち時間 \n%s", now, text)
}
