package models

import "time"

type TaskRecord struct {
	ID       string    `json:"id"`
	Content  string    `json:"content"`
	Priority int       `json:"priority"` // Valid range: 0 to 3
	Updated  time.Time `json:"updated"`
	Due      time.Time `json:"due"`
}
