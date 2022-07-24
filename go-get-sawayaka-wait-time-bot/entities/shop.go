package entites

import (
	"fmt"
	"time"
)

type Shop struct {
	Name      string    `firestore:"name"`
	Keywords  []string  `firestore:"-"`
	Timestamp time.Time `firestore:"timestamp"`
	WaitTime  string    `firestore:"wait_time"`
	WaitSet   string    `firestore:"wait_set"`
}

var ShopFiels = struct {
	Name      string
	Keywords  string
	Timestamp string
	WaitTime  string
	WaitSet   string
}{
	Name:      "name",
	Keywords:  "keywords",
	Timestamp: "timestamp",
	WaitTime:  "wait_time",
	WaitSet:   "wait_set",
}

func (s *Shop) ConvertText() string {
	if s.WaitTime == "-" {
		return fmt.Sprintf("%s店: 営業時間外", s.Name)
	}
	return fmt.Sprintf("%s店: %s分 %s組", s.Name, s.WaitTime, s.WaitSet)
}
