package siteconfig

import (
	"html/template"
	"phagego/phagev2/controllers/sysmanage"
	. "phagego/phagev2/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

func validate(siteConfig *SiteConfig) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(siteConfig.Code, "errmsg").Message("代码必选")
	valid.Required(siteConfig.Value, "errmsg").Message("值必填")
	valid.MaxSize(siteConfig.Value, 255, "errmsg").Message("值最长255位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type SiteConfigIndexController struct {
	sysmanage.BaseController
}

func (this *SiteConfigIndexController) Get() {
	var siteConfigList []SiteConfig
	o := orm.NewOrm()
	qs := o.QueryTable(new(SiteConfig))
	qs.All(&siteConfigList)
	// 返回值
	this.Data["dataList"] = siteConfigList
	this.TplName = "sysmanage/siteconfig/index.html"
}

func (this *SiteConfigIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	siteConfig := SiteConfig{Id: id}
	o := orm.NewOrm()
	err := o.Read(&siteConfig)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	} else if siteConfig.IsSystem == 1 {
		msg = "系统内置，不能删除"
		return
	}
	_, err1 := o.Delete(&SiteConfig{Id: id})
	if err1 != nil {
		beego.Error("Delete siteconfig error", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type SiteConfigAddController struct {
	sysmanage.BaseController
}

func (this *SiteConfigAddController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "sysmanage/siteconfig/add.html"
}

func (this *SiteConfigAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	siteConfig := SiteConfig{}
	if err := this.ParseForm(&siteConfig); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&siteConfig); hasError {
		msg = errMsg
		return
	}
	siteConfig.Creator = this.LoginAdminId
	siteConfig.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&siteConfig, "Code");err != nil {
		msg = "添加失败,请重试"
		beego.Error("Insert siteconfig error", err)
	} else if !created {
		msg = "添加失败，配置已存在"
	} else {
		code = 1
		msg = "添加成功"
	}
}

type SiteConfigEditController struct {
	sysmanage.BaseController
}

func (this *SiteConfigEditController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	siteConfig := SiteConfig{Id: id}

	err := o.Read(&siteConfig)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("SiteConfigIndexController.get"), 302)
	} else {
		this.Data["data"] = siteConfig
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "sysmanage/siteconfig/edit.html"
	}
}

func (this *SiteConfigEditController) Post() {
	var code int
	var msg string
	var reurl = this.URLFor("SiteConfigIndexController.Get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &reurl)
	siteConfig := SiteConfig{}
	if err := this.ParseForm(&siteConfig); err != nil {
		msg = "参数异常"
		return
	}
	siteConfig.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&siteConfig, "Value", "ModifyDate"); err != nil {
		msg = "更新失败"
		beego.Error("Update siteconfig error", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
