package sysmanage

import (
	. "phagego/phagev2/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	. "phagego/phagev2/models"
	"github.com/astaxie/beego/orm"
	"time"
)

type NestPreparer interface {
	NestPrepare()
}

type BaseController struct {
	beego.Controller
	LoginAdminId int64
}

func (this *BaseController) Prepare() {
	adminId, ok := this.GetSession("loginAdminId").(int64)
	if !ok {
		this.Data["json"] = map[string]interface{}{"msg": "请登录"}
		this.ServeJSON()
		this.StopRun()
	}
	this.LoginAdminId = adminId
	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

func (this *BaseController) Index() {
	lid, _ := this.GetSession("loginAdminId").(int64)
	// 获取左侧菜单
	o := orm.NewOrm()
	sql := "select * from ph_permission a where a.enabled = 1 and display = 1 and exists(select b.id from ph_role_permission b, ph_admin_role c where b.role_id = c.role_id and b.permission_id = a.id and c.admin_id = ?) order by a.pid, a.sort, a.id"
	var permissions []Permission
	_, err := o.Raw(sql, lid).QueryRows(&permissions)
	if err != nil {
		beego.Error("Query admin permission error", err)
		this.Abort("内部错误，请重试")
	} else {
		var mainMenuList []Permission
		secdMenuMap := make(map[int64][]Permission)
		for _, pe := range permissions {
			// 构建菜单
			if pe.Pid == 0 {
				mainMenuList = append(mainMenuList, pe)
			} else {
				if val, ok := secdMenuMap[pe.Pid]; ok {
					val = append(val, pe)
					secdMenuMap[pe.Pid] = val
				} else {
					var menuList []Permission
					menuList = append(menuList, pe)
					secdMenuMap[pe.Pid] = menuList
				}
			}
		}
		this.Data["loginAdminName"] = this.GetSession("loginAdminName")
		this.Data["mainMenuList"] = mainMenuList
		this.Data["secdMenuMap"] = secdMenuMap
		// 站点信息
		this.Data["siteName"] = GetSiteConfigValue(Scname)
		this.Data["year"] = time.Now().Year()
	}
	this.TplName = "sysmanage/base.html"
}

func Retjson(ctx *context.Context, msg *string, code *int, data ...interface{}) {
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	if len(data) > 0 {
		d := data[0]
		switch d.(type) {
		case string:
			ret["url"] = d
			break
		case *string:
			ret["url"] = d
			break
		}
		ret["data"] = d
	}
	ctx.Output.JSON(ret, false, false)
}
