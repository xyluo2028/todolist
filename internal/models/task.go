package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`    // title of the task, upto 100 characters
	Content     string    `json:"content"`  // detailed description about the task, upto 1k characters
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
