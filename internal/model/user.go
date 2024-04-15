package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username      string     `gorm:"type:varchar(64);unique;not null;comment:用户登录名"`
	UniqueNumber  string     `gorm:"type:varchar(30);unique;comment:编码"`
	Password      string     `gorm:"type:varchar(128);not null;comment:用户登录密码"`
	Status        int        `gorm:"type:tinyint;default:1;comment:状态"`
	NickName      string     `gorm:"type:varchar(64);default:'';comment:用户昵称"`
	Avatar        string     `gorm:"type:varchar(256);default:'';comment:用户头像"`
	Gender        int        `gorm:"type:tinyint;default:0;comment:性别"`
	Phone         string     `gorm:"type:varchar(11);default:'';comment:用户手机号"`
	Email         string     `gorm:"type:varchar(64);unique;not null;comment:用户邮箱"`
	Telegram      *string    `gorm:"type:varchar(64);unique;default:null;comment:telegram账号"`
	Discord       *string    `gorm:"type:varchar(64);unique;default:null;comment:discord账号"`
	Banner        string     `gorm:"type:varchar(128);default:'';comment:背景图url"`
	Intro         string     `gorm:"type:varchar(512);default:'';comment:简介"`
	Name          string     `gorm:"type:varchar(64);default:'';comment:姓名"`
	BirthDay      *time.Time `gorm:"type:datetime;default:null;comment:生日"`
	Address       string     `gorm:"type:varchar(128);default:'';comment:所在地"`
	WithdrawState int        `gorm:"type:tinyint;default:0;comment:'提取状态'"`
	ParentId      uint       `gorm:"default:0;type:bigint;comment:'上级id'"`
	InviteNum     int        `gorm:"type:int;default:0;comment:邀请人数"`
	HasPaid       bool       `gorm:"tinyint;default:0;comment:'是否购买过盲盒'"`
	IsInternal    bool       `gorm:"tinyint;default:0;comment:'是否是内部账号'"`
	IsVip         bool       `gorm:"tinyint;default:0;comment:'是否是VIP用户'"`
	IsActive      bool       `gorm:"tinyint;default:0;comment:'用户是否活跃，48小时计算'"`
	VipLevel      int        `gorm:"type:smallint(2);default:0;comment:'VIP等级'"`
	FragCnt       int        `gorm:"type:integer;default:0;comment:'装备碎片的数量'"`
	ComAmt        float64    `gorm:"default:0;type:decimal(40,20);comment:'给上级贡献总金额'"`
	CalcVipTime   time.Time  `gorm:"type:datetime;default:'1970-01-01 00:00:00';comment:上级计算VIP等级的时间"`
	VipTime       time.Time  `gorm:"type:datetime;default:'1970-01-01 00:00:00';comment:上次vip等级更新时间"`
	ActiveTime    time.Time  `gorm:"type:datetime;default:'1970-01-01 00:00:00';comment:最近活跃时间"`
}

const (
	UserStatusForbidden  = -2 // 禁用
	UserStatusLocked     = -1 // 锁定
	UserStatusRegistered = 1  // 已注册
	UserStatusVerified   = 2  // 已验证
)

const (
	UserGenderUnknown = iota
	UserGenderMale
	UserGenderFemale

	UserWithdrawPermit = 1
)

func (u *User) WithdrawPermit() bool {
	return u.WithdrawState == UserWithdrawPermit
}

func (u *User) IsLocked() bool {
	return u.Status == UserStatusLocked
}

func (u *User) IsForbidden() bool {
	return u.Status == UserStatusForbidden
}

func (u *User) TelegramString() string {
	if u.Telegram == nil {
		return ""
	}
	return *u.Telegram
}

func (u *User) DiscordString() string {
	if u.Discord == nil {
		return ""
	}
	return *u.Discord
}
