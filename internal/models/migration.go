package models

import "time"

type Roles struct {
	ID          int64  `gorm:"id,primary_key"`
	Name        string `gorm:"size:255"`
	IsActive    bool
	CreatedDate time.Time
	CreatedBy   int64
	UpdatedDate time.Time
	UpdatedBy   int64
}

type Rooms struct {
	ID           int64  `gorm:"id,primary_key"`
	Title        string `gorm:"size:255"`
	Description  string
	ScheduleTime time.Time
	IsPremium    bool
	CreatedDate  time.Time
	CreatedBy    int64
	UpdatedDate  time.Time
	UpdatedBy    int64
}

type Streams struct {
	ID         int64 `gorm:"id,primary_key"`
	UserID     int64
	RoomID     int64
	StreamPath string
	SessionID  string
	StartDate  time.Time
	EndDate    time.Time
}

type Subcriptions struct {
	ID          int64 `gorm:"id,primary_key"`
	UserID      int64
	RoomID      int64
	CreatedDate time.Time
	CreatedBy   int64
}

type Users struct {
	ID           int64 `gorm:"id,primary_key"`
	Username     string
	Password     string
	Name         string
	LoginType    string
	Email        string
	RoleID       int64
	IsActive     bool
	IsLogin      bool
	RegisterDate time.Time
	UpdatedDate  time.Time
}
