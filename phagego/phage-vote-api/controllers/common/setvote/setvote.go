package setvote

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"phage/controllers/sysmanage"
	. "phagego/phage-vote-api/models"
)

type IndexSetVoteController struct {
	sysmanage.BaseController
}

func (v *IndexSetVoteController) Get() {
	o := orm.NewOrm()
	var sv []SetVote
	_, _ = o.QueryTable(new(SetVote)).All(&sv)
	if len(sv) < 11 {
		for i := 1; i < 13; i++ {
			var s SetVote
			s.Category = i
			beego.Info(i)
			s.ReadOrCreate("Category")
		}
	}
	v.Data["dataList"] = sv
	v.TplName = "common/setvote/index.html"
}

type EditSetVoteController struct {
	sysmanage.BaseController
}

func (v *EditSetVoteController) Get() {
	i, _ := v.GetInt("id")
	o := orm.NewOrm()
	var sv SetVote
	o.QueryTable(new(SetVote)).Filter("Id", i).One(&sv)
	v.Data["data"] = sv
	v.TplName = "common/setvote/edit.html"
}

func (v *EditSetVoteController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(v.Ctx, &msg, &code)
	sv := SetVote{}
	if err := v.ParseForm(&sv); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"Value"}
	_, err1 := sv.Update(cols...)
	if err1 != nil {
		msg = "更新失败"
		beego.Error("Update setvote error", err1)
	} else {
		code = 1
		msg = "更新成功"
	}
}
