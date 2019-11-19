package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type AdminRole struct {
	Id      int64
	AdminId int64
	RoleId  int64
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(AdminRole))
}

func (model *AdminRole) TableUnique() [][]string {
	return [][]string{
		[]string{"AdminId", "RoleId"},
	}
}

//获取所有客服组
func GetAdminRoles()[]AdminRole{
	 var adminrole []AdminRole
	 o := orm.NewOrm()
	 o.QueryTable(new(AdminRole)).Filter("RoleId",3).All(&adminrole)
	 return  adminrole
}
