package membertotal

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"os"
	"phagego/frameweb-v2/controllers/sysmanage"
	."phagego/phage-vip3-web/models/common"
	"strings"
	"time"
)

type MemberTotalIndexController struct {
	sysmanage.BaseController
}

func (this *MemberTotalIndexController) Get() {
	//导出
	isExport, _ := this.GetInt("isExport", 0)
	if isExport == 1 {
		this.Export()
		return
	}
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	account := strings.TrimSpace(this.GetString("account"))
	limit, _ := beego.AppConfig.Int("pagelimit")
	list,total := new(MemberTotal).Paginate(page,limit,account)
	pagination.SetPaginator(this.Ctx,limit,total)
	this.Data["account"] = account
	this.Data["dataList"] = list
	this.TplName = "common/membertotal/index.html"
}

func (this *MemberTotalIndexController) Delbatch(){
	  var code int
	  var msg  string
	  defer sysmanage.Retjson(this.Ctx,&msg,&code)
	  o := orm.NewOrm()
	  res, err := o.Raw("DELETE from ph_member_total WHERE id != 0").Exec()
	  if err != nil {
		  beego.Error("删除所的会员统计失败",err)
		  msg = "删除失败"
		  return
	  } else {
	  	code = 1
	  	num, _ := res.RowsAffected()
		  msg = fmt.Sprintf("成功删除%d条记录", num)
		  return
	  }
}

func (this *MemberTotalIndexController) Export(){
     o := orm.NewOrm()
     var  membertotal []MemberTotal
     _, err := o.QueryTable(new(MemberTotal)).OrderBy("-Bet").Limit(-1).All(&membertotal)
	if err != nil {
		beego.Error("导出会员统计列表失败",err)
	}
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "会员账号")
	xlsx.SetCellValue("Sheet1", "B1", "VIP等级")
	xlsx.SetCellValue("Sheet1", "C1", "投注额")
	xlsx.SetCellValue("Sheet1", "D1", "晋级总彩金")
	xlsx.SetCellValue("Sheet1", "E1", "俸禄总额")
	for i, value := range membertotal {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), value.Account)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), value.Level)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), value.Bet)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), value.TotalLevelGift)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), value.TotalLuckyGift)

	}
	fileName := fmt.Sprintf("./tmp/excel/memberlist_%s.xlsx", time.Now().Format("20060102150405"))
	err1 := xlsx.SaveAs(fileName)
	if err1 != nil {
		beego.Error("导出会员列表失败", err.Error())
	} else {
		defer os.Remove(fileName)
		this.Ctx.Output.Download(fileName)
	}
}

type MemberTotalEditController struct {
	sysmanage.BaseController
}


func (this *MemberTotalEditController) Get()  {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	mt := MemberTotal{Id: id}
	err := o.Read(&mt)
	if err != nil {
		this.Redirect("membertotaleIndexController.get", 302)
	} else {
		this.Data["data"] = mt
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/membertotal/edit.html"
	}
}

func (this *MemberTotalEditController) Post()  {
	var code int
	var msg string
	url := beego.URLFor("MembertotalIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	mt := MemberTotal{}
	if err := this.ParseForm(&mt); err != nil {
		beego.Error("更新会员总投注信息异常", err)
		msg = "参数异常"
		return
	}
	cols := []string{"Level", "Bet"}
	mt.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	_, err := o.Update(&mt, cols...)
	if err != nil {
		beego.Error("更新会员总投注信息失败", err)
		msg = "更新失败"
		return
	} else {
		code = 1
		msg = "更新成功"
	}
}



