package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id         int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
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
func (u User) TableName() string {
	shardIndex := u.PutDate / 100
	return fmt.Sprintf("user%6d", shardIndex)
}

type UserRepo interface {
	SaveUser(*gorm.DB, *User) error
	UpdateEnterGroupTime(*gorm.DB, *User) error
}

type myUserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &myUserRepo{db: db}
}

func (up *myUserRepo) SaveUser(db *gorm.DB, u *User) error {
	return db.Table(u.TableName()).Create(u).Error
}

func (up *myUserRepo) UpdateEnterGroupTime(db *gorm.DB, u *User) error {
	updates := map[string]any{
		"enter_group": u.EnterGroup,
		"update_at":   time.Now().Unix(),
	}
	ret := db.Table(u.TableName()).Model(u).Where("id = ?", u.Id).Updates(updates)
	return ret.Error
}
