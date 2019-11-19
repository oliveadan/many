package wechat

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/z6hewx/models"
	"strings"
	"time"
)

type WechatContorller struct {
	sysmanage.BaseController
}

func (this *WechatContorller) Get() {
	wechtaccount := strings.TrimSpace(this.GetString("wechat"))
	enabled, _ := this.GetInt8("enabled", -1)
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(Wechat).Paginate(page, limit, wechtaccount, enabled)
	pagination.SetPaginator(this.Ctx, limit, total)
	//返回值
	this.Data["condArr"] = map[string]interface{}{
		"wechtaccount": wechtaccount,
		"enabled":      enabled}
	this.Data["dataList"] = list
	this.TplName = "common/wechat/index.html"
}

func (this *WechatContorller) Delone() {
	var code int
	var msg string
	url := beego.URLFor("WechatContorller.get","time",time.Now().Unix())
	defer sysmanage.Retjson(this.Ctx, &msg, &code,&url)
	id, _ := this.GetInt64("id")
	wechat := Wechat{Id: id}
	o := orm.NewOrm()
	_, err1 := o.Delete(&wechat, "Id")
	if err1 != nil {
		beego.Error("删除会员账号失败", err1)
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (this *WechatContorller) Enable() {
	var code int
	var msg string
	//加上时间戳强制清除缓存
	url := beego.URLFor("WechatContorller.get","time",time.Now().Unix())
	defer sysmanage.Retjson(this.Ctx, &msg, &code,&url)
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	wechat := Wechat{Id: id}
	err := o.Read(&wechat)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		msg = "数据不存在，请确认"
		return
	}
	if wechat.Enabled == 0 {
		wechat.Enabled = 1
	} else {
		wechat.Enabled = 0
	}
	wechat.Modifior = this.LoginAdminId
	_, err1 := wechat.Update("Enabled")
	if err1 != nil {
		beego.Error("激活失败", err1)
		msg = "操作失败"
	} else {
		code = 1
		msg = "操作成功"
	}
}

type WechatAddContorller struct {
	sysmanage.BaseController
}

func (this *WechatAddContorller) Get() {
	this.TplName = "common/wechat/add.html"
}

func (this *WechatAddContorller) Post() {
	var code int
	var msg string
	var url = beego.URLFor("WechatContorller.Get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	wechat := Wechat{}
	if err := this.ParseForm(&wechat); err != nil {
		msg = "参数异常"
		return
	}
	wechat.Creator = this.LoginAdminId
	wechat.Modifior = this.LoginAdminId
	_, err1 := wechat.Create()
	if err1 != nil {
		msg = "添加失败"
		beego.Error("添加微信账账号失败", err1)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type WechatEditContorller struct {
	sysmanage.BaseController
}

func (this *WechatEditContorller) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	wechat := Wechat{Id: id}

	err := o.Read(&wechat)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("WechatContorller.get"), 302)
	} else {
		this.Data["data"] = wechat
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/wechat/edit.html"
	}
}

func (this *WechatEditContorller) Post() {
	var code int
	var msg string
	url := beego.URLFor("WechatContorller.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	wechat := Wechat{}
	if err := this.ParseForm(&wechat); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"WxNo", "QrCode"}
	_, err1 := wechat.Update(cols...)
	if err1 != nil {
		msg = "更新失败"
		beego.Error("更新微信账号失败", err1)
	} else {
		code = 1
		msg = "更新成功"
	}
}
