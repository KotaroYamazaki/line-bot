package entites

import "fmt"

type Shop struct {
	Name     string `firestore:"name"`
	WaitTime string `firestore:"wait_time"`
	WaitSet  string `firestore:"wait_set"`
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
	return fmt.Sprintf("%s店: %s分 %s組", s.Name, s.WaitTime, s.WaitSet)
}
