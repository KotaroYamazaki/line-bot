package entities

import (
	"errors"
	"strings"
)

type Query string

var all = []string{"all", "すべて", "全て", "全部", "あ"}

func NewQuery(s string) (Query, error) {
	q := Query(strings.TrimSpace(s))
	if err := q.validate(); err != nil {
		return "", err
	}
	return q, nil
}

func (q Query) validate() error {
	if q.String() == "" {
		return errors.New("not valid empty")
	}
	return nil
}

func (q Query) String() string {
	return string(q)
}

func (q Query) IsALL() bool {
	for _, a := range all {
		if q.String() == a {
			return true
		}
	}
	return false
}
