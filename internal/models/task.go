package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	Priority    int       `json:"priority"` // the lower the number, the higher the priority
	UpdatedTime time.Time `json:"updatedTime"`
	Due         time.Time `json:"due"`
	Completed   bool      `json:"completed"`
}

func (t Task) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%+v", t)
	}
	return string(b)
}
