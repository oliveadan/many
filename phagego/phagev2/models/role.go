package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id          int64     `auto`                              // 自增主键
	CreateDate  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator     int64                                         // 创建人Id
	Modifior    int64                                         // 更新人Id
	Version     int                                           // 版本
	Enabled     int8                                          // 是否启用
	Description string    `orm:"null"`                        // 描述
	IsSystem    int8                                          // 是否内置(内置不可选择)
	Name        string                                        // 名称
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(Role))
}

func GetRoleList() (roles []Role) {
	var roleList []Role
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
	cond := orm.NewCondition()
	qs = qs.SetCond(cond.And("Enabled", 1).And("IsSystem", 0))
	qs.All(&roleList)
	return roleList
}
