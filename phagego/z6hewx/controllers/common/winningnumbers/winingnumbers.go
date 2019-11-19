package winningnumbers

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

type WinningnumbersContorller struct {
	sysmanage.BaseController
}

func (this *WinningnumbersContorller) Get() {
	Winningnumbers := strings.TrimSpace(this.GetString("Winningnumbers"))
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(WinningNumbers).Paginate(page, limit)
	pagination.SetPaginator(this.Ctx, limit, total)
	//返回值
	this.Data["condArr"] = map[string]interface{}{
		"Winningnumbers": Winningnumbers}
	this.Data["dataList"] = list
	this.TplName = "common/winningnumbers/index.html"
}

func (this *WinningnumbersContorller) Delone() {
	var code int
	var msg string
	url := beego.URLFor("WinningnumbersContorller.get","time",time.Now().Unix())
	defer sysmanage.Retjson(this.Ctx, &msg, &code,&url)
	id, _ := this.GetInt64("id")
	Winningnumbers := WinningNumbers{Id: id}
	o := orm.NewOrm()
	_, err1 := o.Delete(&Winningnumbers, "Id")
	if err1 != nil {
		beego.Error("删除会员账号失败", err1)
	} else {
		code = 1
		msg = "删除成功"
	}
}

type WinningnumbersAddContorller struct {
	sysmanage.BaseController
}

func (this *WinningnumbersAddContorller) Get() {
	this.TplName = "common/winningnumbers/add.html"
}

func (this *WinningnumbersAddContorller) Post() {
	var code int
	var msg string
	var url = beego.URLFor("WinningnumbersContorller.Get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	winningnumbers := WinningNumbers{}
	if err := this.ParseForm(&winningnumbers); err != nil {
		msg = "参数异常"
		return
	}
	var s1 string
	if strings.Contains(winningnumbers.Numbers, "，") {
		s1 = strings.Replace(winningnumbers.Numbers, "，", ",", -1)
	} else {
		s1 = winningnumbers.Numbers
	}
	winningnumbers.Numbers = s1
	winningnumbers.Creator = this.LoginAdminId
	winningnumbers.Modifior = this.LoginAdminId
	_, err1 := winningnumbers.Create()
	if err1 != nil {
		msg = "添加失败"
		beego.Error("添加微信账账号失败", err1)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type WinningnumbersEditContorller struct {
	sysmanage.BaseController
}

func (this *WinningnumbersEditContorller) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	Winningnumbers := WinningNumbers{Id: id}

	err := o.Read(&Winningnumbers)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("WinningnumbersContorller.get"), 302)
	} else {
		this.Data["data"] = Winningnumbers
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/winningnumbers/edit.html"
	}
}

func (this *WinningnumbersEditContorller) Post() {
	var code int
	var msg string
	url := beego.URLFor("WinningnumbersContorller.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	winningnumbers := WinningNumbers{}
	if err := this.ParseForm(&winningnumbers); err != nil {
		msg = "参数异常"
		return
	}
	var s1 string
	if strings.Contains(winningnumbers.Numbers, "，") {
		s1 = strings.Replace(winningnumbers.Numbers, "，", ",", -1)
	} else {
		s1 = winningnumbers.Numbers
	}
	winningnumbers.Numbers = s1
	cols := []string{"Period", "Numbers"}
	_, err1 := winningnumbers.Update(cols...)
	if err1 != nil {
		msg = "更新失败"
		beego.Error("更新微信账号失败", err1)
	} else {
		code = 1
		msg = "更新成功"
	}
}
