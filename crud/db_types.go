package crud

import "time"

type Category struct {
	CatID  int
	UserID string `json:"-"`
	Name   string
	Color  int
}

type Team struct {
	TeamID  int64
	Name    string
	Members []map[string]interface{}
}

type Task struct {
	TaskID         int
	UserID         string `json:"-"`
	Category       string
	TaskName       string
	Description    string
	StartTime      time.Time
	EndTime        time.Time
	Status         string
	IsRecurring    bool
	IsAllDay       bool
	Difficulty     string
	CronExpression string
	TeamID         int
}

type RecurTypeTask struct {
	TaskID       int
	UserID       string `json:"-"`
	Category     string
	TaskName     string
	StartTime    time.Time
	EndTime      time.Time
	Status       string
	IsRecurring  bool
	IsAllDay     bool
	Difficulty   string
	RecurrenceId int
}

type User struct {
	UserID     string `json:"-"` // Not known to user, do not expose
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
