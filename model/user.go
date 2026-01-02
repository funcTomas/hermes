package model

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id         uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	PutDate    int    `gorm:"column:put_date"`
	Phone      string `gorm:"column:phone"`
	UniqId     string `gorm:"column:uniq_id"`
	ChannelId  int    `gorm:"column:channel_id"`
	EnterGroup int64  `gorm:"column:enter_group"`
	Step       int    `gorm:"column:step"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdateAt   int64  `gorm:"column:update_at"`
}

// TableName 表名
func (u *User) TableName() string {
	shardIndex := u.PutDate / 100
	return fmt.Sprintf("user%6d", shardIndex)
}

func (u *User) AddUser(ctx context.Context, db *gorm.DB) error {
	ret := db.WithContext(ctx).Create(u)
	return ret.Error
}

func (u *User) UpdateEnterGroupTimeById(ctx context.Context, db *gorm.DB) error {
	updates := map[string]any{
		"enter_group": u.EnterGroup,
		"update_at":   time.Now().Unix(),
	}
	ret := db.WithContext(ctx).Model(u).Where("id = ?", u.Id).Updates(updates)
	return ret.Error
}
