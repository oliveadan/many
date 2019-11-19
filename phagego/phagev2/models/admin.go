package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Admin struct {
	Id                int64     `auto`                              // 自增主键
	CreateDate        time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate        time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator           int64                                         // 创建人Id
	Modifior          int64                                         // 更新人Id
	Version           int                                           // 版本
	Enabled           int8                                          // 是否启用
	Locked            int8                                          // 是否锁定
	IsSystem          int8                                          // 是否系统内置
	LockedDate        time.Time `orm:"null"`                        // 锁定时间
	LoginDate         time.Time `orm:"null"`                        // 登录时间
	LoginFailureCount int                                           // 登录失败次数
	LoginIp           string    `orm:"null"`                        // 登录ip
	Salt              string                                        // 盐
	Name              string                                        // 名称
	Password          string                                        // 密码
	Username          string    `orm:"unique"`                      // 用户名
	Email             string    `orm:"null"`                        // 邮箱
	Mobile            string    `orm:"null"`                        // 手机
	LoginVerify       int8                                          // 是否开启登录验证
}

func init() {
	beego.Info("Init model admin")
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(Admin))
}
