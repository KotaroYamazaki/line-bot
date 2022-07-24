package entites

import "fmt"

type Shop struct {
	name     string   `firestore:"name"`
	keywords []string `firestore:"-"`
	waitTime string   `firestore:"wait_time"`
	waitSet  string   `firestore:"wait_set"`
}

var ShopFiels = struct {
	Name     string
	Keywords string
	WaitTime string
	WaitSet  string
}{
	Name:     "name",
	Keywords: "keywords",
	WaitTime: "wait_time",
	WaitSet:  "wait_set",
}

func (s *Shop) ConvertText() string {
	if s.waitTime == "-" {
		return fmt.Sprintf("%s店: 営業時間外", s.name)
	}
	return fmt.Sprintf("%s店: %s分 %s組", s.name, s.waitTime, s.waitSet)
}
