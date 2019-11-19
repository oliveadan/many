package front

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	utils2 "phage/utils"
	"phagego/frameweb-v2/controllers/sysmanage"
	"phagego/frameweb-v2/models"
	. "phagego/phage-vip-web/models/common"
	"phagego/phage-vip-web/utils"
	"strings"
)

type FrontIndexController struct {
	sysmanage.BaseController
}

func (this *FrontIndexController) Prepare() {
	this.EnableXSRF = false
}

func (this *FrontIndexController) Get() {
	//前台配置信息
	sc := models.GetSiteConfigMap(utils2.Scname,utils.Scofficial, utils.Scranking, utils.Scregister, utils.Sccust, utils.Scfqa,utils.Scpromotion)
	this.Data["officialSite"] = sc[utils.Scofficial]
	this.Data["officialRegist"] = sc[utils.Scregister]
	this.Data["custServ"] = sc[utils.Sccust]
	this.Data["officialFqa"] = sc[utils.Scfqa]
	this.Data["siteName"] = sc[utils2.Scname]
	this.Data["ranking"] = sc[utils.Scranking]
	this.Data["officialPromot"] = sc[utils.Scpromotion]
	this.TplName = "front/index.html"
}

func (this *FrontIndexController) Query() {
	data := make(map[string]interface{})
	account := strings.TrimSpace(this.GetString("account"))
	o := orm.NewOrm()
	//会员统计信息
	var membertotal MemberTotal
	err := o.QueryTable(new(MemberTotal)).Filter("Account", account).One(&membertotal)
	if err != nil {
		beego.Error("获取会员统计信息失败", err)
	}
	//距离下一级VIP需要的打码量
	var level []Level
	var balance int64
	_, _ = o.QueryTable(new(Level)).OrderBy("-VipLevel").All(&level)
	for _, v := range level {
		if membertotal.Level+1 == v.VipLevel {
			balance = v.TotalBet - membertotal.Bet
		}
	}
	//期数信息
	var membersingle []MemberSingle
	_, _ = o.QueryTable(new(MemberSingle)).Filter("Account", account).OrderBy("-PeriodSeq").All(&membersingle)


	data["balance"] = balance
	data["membertotal"] = membertotal
	data["membersingle"] = membersingle
	this.Ctx.Output.JSON(data, false, false)
}
