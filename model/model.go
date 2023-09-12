package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

// User
type User struct {
	ID              int64  `json:"id"`               // 用户id
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数
	Username        string `json:"-" gorm:"unique"`  // 注册用户名，最长32个字符
	PasswordHash    string `json:"-"`                // 密码，最长32个字符   service层完成对应的逻辑操作
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	u.Name = u.Username
	u.Signature = "谢谢你的关注"
	return tx.Model(u).Updates(User{Name: u.Name, Signature: u.Signature}).Error
}

// relation
type Relation struct {
	ID       int64 `json:"-" gorm:"primaryKey"` // 关注记录唯一标识
	AuthorID int64 `json:"-" gorm:"index"`      // 作者ID
	FansID   int64 `json:"-" gorm:"index"`      // 粉丝ID
}

// message
type Message struct {
	ID         int64  `json:"id"`           // 消息id
	FromUserID int64  `json:"from_user_id"` // 消息发送者id
	ToUserID   int64  `json:"to_user_id"`   // 消息接收者id
	Content    string `json:"content"`      // 消息内容
	CreateTime string `json:"create_time"`  // 消息发送时间 yyyy-MM-dd HH:MM:ss
}

func (m *Message) AfterFind(tx *gorm.DB) (err error) {
	t, err := time.Parse("2006-01-02 15:04:05", m.CreateTime)
	if err != nil {
		return err
	}
	m.CreateTime = strconv.FormatInt(t.Unix(), 10)
	return nil
}
