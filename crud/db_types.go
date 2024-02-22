package crud

import "time"

type Category struct {
	CatID  int
	UserID string
	Name   string
	Color  int
}

type Task struct {
	TaskID         int
	UserID         string
	Category       string
	TaskName       string
	Description    string
	StartTime      time.Time
	EndTime        time.Time
	Status         string
	IsRecurring    bool
	IsAllDay       bool
	RecurringType  string
	Difficulty     string
	CronExpression string
}

type TaskPreview struct {
	TaskID      int
	UserID      string
	Category    string
	TaskName    string
	StartTime   time.Time
	EndTime     time.Time
	Status      string
	IsRecurring bool
	IsAllDay    bool
}

type User struct {
	UserID     string // Not known to user, do not expose
	Username   string // Set by user, can be exposed
	Picture    string // A0 stores their profile pics as URLs
	Points     int
	BossId     int
	SocialCode string // Friendly code to uniquely identify (public)
}

type Boss struct {
	BossID int
	Name   string
	Health int
	Image  string
}
