package repository

import "time"

type Task struct {
	TaskId    uint `gorm:"primarykey"`
	UserId    uint `gorm:"index"`
	Status    int  `gorm:"default:0"`
	Title     string
	Content   string `gorm:"longtext"`
	StartTime time.Time
	EndTime   time.Time
}
