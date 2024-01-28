package main

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	TaskID      int
	UserID      string
	Category    string
	TaskName    string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	IsCompleted bool
	IsRecurring bool
	IsAllDay    bool
}

type TaskPreview struct {
	TaskID      int
	UserID      string
	Category    string
	TaskName    string
	StartTime   time.Time
	EndTime     time.Time
	IsCompleted bool
	IsRecurring bool
	IsAllDay    bool
}
