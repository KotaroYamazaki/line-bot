package entities

import (
	"fmt"
)

func GenerateMessage(text string) string {
	if text == "" {
		return fmt.Sprintln("該当する店舗がありませんでした")
	}
	return fmt.Sprintf("%s", text)
}
