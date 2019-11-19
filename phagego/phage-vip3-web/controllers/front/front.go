package front

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	utils2 "phage/utils"
	"phagego/frameweb-v2/controllers/sysmanage"
	"phagego/frameweb-v2/models"
	. "phagego/phage-vip3-web/models/common"
	"phagego/phage-vip3-web/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FrontIndexController struct {
	sysmanage.BaseController
}

func (this *FrontIndexController) Prepare() {
	this.EnableXSRF = false
}

func (this *FrontIndexController) Get() {
	//前台配置信息
	sc := models.GetSiteConfigMap(utils2.Scname, utils.Scofficial, utils.Scranking, utils.Scregister, utils.Sccust, utils.Scfqa, utils.Scpromotion, utils.Scnotice)
	this.Data["officialSite"] = sc[utils.Scofficial]
	this.Data["officialRegist"] = sc[utils.Scregister]
	this.Data["custServ"] = sc[utils.Sccust]
	this.Data["officialFqa"] = sc[utils.Scfqa]
	this.Data["siteName"] = sc[utils2.Scname]
	this.Data["ranking"] = sc[utils.Scranking]
	this.Data["officialPromot"] = sc[utils.Scpromotion]
	this.Data["notice"] = sc[utils.Scnotice]
	this.TplName = "front/index.html"
}

func (this *FrontIndexController) Query() {
	data := make(map[string]interface{})
	account := strings.TrimSpace(this.GetString("account"))
	o := orm.NewOrm()

	//期数信息
	var membersingle []MemberSingle
	_, _ = o.QueryTable(new(MemberSingle)).Filter("Account", account).OrderBy("-PeriodSeq").All(&membersingle)
	for i, v := range membersingle {
		if v.LuckyEnable == 0 {
			//如果24内后礼物还没有领取将作废
			end := time.Now().Format("2006-01-02 15:04:05")
			start := v.CreateDate.Format("2006-01-02 15:04:05")
			h := utils.GetHourDiffer(start, end)
			if h > 24 {
				_, _ = o.QueryTable(new(MemberSingle)).Filter("Id", v.Id).Update(orm.Params{"LuckyEnable": 3})
				membersingle[i].LuckyEnable = 3
			}
		}
	}
	//获取累计晋级彩金和累计好运金
	var maps []orm.Params
	_, _ = o.Raw("SELECT  sum(level_gift),sum(lucky_gift) FROM ph_member_single where account = ?", account).Values(&maps)
	totallevelgift, _ := strconv.ParseInt(maps[0]["sum(level_gift)"].(string), 10, 64)
	totalluckygift, _ := strconv.ParseInt(maps[0]["sum(lucky_gift)"].(string), 10, 64)
	var levels []Level
	_, _ = o.QueryTable(new(Level)).OrderBy("-TotalBet").All(&levels)

	var memberTotal MemberTotal
	err := o.QueryTable(new(MemberTotal)).Filter("Account", account).One(&memberTotal)
	if err != nil {
		beego.Error("query Membertotal error", err)
	}

	//距离下一级VIP需要的打码量
	var level Level
	/*
		var balance int64
		_, _ = o.QueryTable(new(Level)).OrderBy("-VipLevel").All(&level)
		for _, v := range level {
			if memberTotal.Level+1 == v.VipLevel {
				balance = v.TotalBet - totalbet
			}
		}*/
	err = o.QueryTable(new(Level)).Filter("VipLevel", memberTotal.Level+1).One(&level)
	if err != nil {
		beego.Error("query level err", err)
	}

	data["Account"] = account
	data["balance"] = level.TotalBet - memberTotal.Bet
	data["membertotalbet"] = memberTotal.Bet
	data["level"] = memberTotal.Level
	data["membersingle"] = membersingle
	data["totallevelgift"] = totallevelgift
	data["totalluckygift"] = totalluckygift
	this.Ctx.Output.JSON(data, false, false)
}

func (this *FrontIndexController) GetGift() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	account := strings.TrimSpace(this.GetString("account"))
	id := this.GetString("id")
	typee, _ := this.GetInt("type")
	gift, _ := this.GetInt64("gift", 0)
	o := orm.NewOrm()
	//判断会员是否已经领取了彩金
	var membersingle MemberSingle
	_ = o.QueryTable(new(MemberSingle)).Filter("Id", id).One(&membersingle)
	if typee == 0 {
		if membersingle.LevelEnable == 1 {
			msg = "您已经成功领取晋级彩金"
			return
		}
	} else {
		if membersingle.LuckyEnable == 1 {
			msg = "您已经成功领取当天好运金"
			return
		}
	}
	var rewardlog RewardLog
	rewardlog.Account = account
	if typee == 0 {
		rewardlog.GiftName = "晋级奖励"
		rewardlog.GiftContent = fmt.Sprintf("%d", gift)
	} else {
		rewardlog.GiftName = "好运金"
		rewardlog.GiftContent = fmt.Sprintf("%d", gift)
	}
	var lock sync.RWMutex
	lock.Lock()
	_, err := rewardlog.Create()
	lock.Unlock()
	if err != nil {
		beego.Error("生成中奖记录失败", err)
		msg = "系统异常，请刷新后重试"
		return
	} else {
		if typee == 0 {
			_, err := o.QueryTable(new(MemberSingle)).Filter("Id", id).Update(orm.Params{"LevelEnable": 1})
			if err != nil {
				beego.Error("更新晋级礼物失败", err)
			}

		} else {
			_, err := o.QueryTable(new(MemberSingle)).Filter("Id", id).Update(orm.Params{"LuckyEnable": 1})
			if err != nil {
				beego.Error("更新好运礼物失败", err)
			}
		}
	}
	code = 1
	msg = "恭喜您领取成功"
}

//总信息修复
func (this *FrontIndexController) RepairMemberTotal() {
	var data = make(map[string]interface{})
	en := this.GetString("Enable")
	var erraccount string
	o := orm.NewOrm()
	if en == "1" {
		var mt []MemberTotal
		_, _ = o.QueryTable(new(MemberTotal)).Limit(-1).All(&mt)
		for _, v := range mt {
			//获取累计晋级彩金和累计好运金
			var maps []orm.Params
			_, _ = o.Raw("SELECT sum(bet),sum(level_gift),sum(lucky_gift) FROM ph_member_single where account = ?", v.Account).Values(&maps)
			totalbet, _ := strconv.ParseInt(maps[0]["sum(bet)"].(string), 10, 64)
			var levels []Level
			_, _ = o.QueryTable(new(Level)).OrderBy("-TotalBet").All(&levels)
			var llevel int
			//获取当前的VIP等级
			for _, l := range levels {
				if totalbet >= l.TotalBet {
					llevel = l.VipLevel
					break
				}
			}
			//当总投注不相同的时候，更新
			if totalbet != v.Bet {
				o.Begin()
				_, e := o.QueryTable(new(MemberTotal)).Filter("Id", v.Id).Update(orm.Params{"Level": llevel, "Bet": totalbet})
				if e != nil {
					erraccount += v.Account + "__"
					o.Rollback()
					continue
				}
				o.Commit()
			}
		}
	}
	if en == "2" {
		o.QueryTable(new(MemberTotal)).Filter("Id__gte", 0).Delete()
		//获取累计晋级彩金和累计好运金
		var mts []MemberTotal
		var maps []orm.Params
		_, _ = o.Raw("SELECT account,sum(bet),sum(level_gift),sum(lucky_gift) FROM ph_member_single GROUP BY account").Values(&maps)
		for _, v := range maps {
			var mt MemberTotal
			totalbet, _ := strconv.ParseInt(v["sum(bet)"].(string), 10, 64)
			level, _ := strconv.ParseInt(v["sum(level_gift)"].(string), 10, 64)
			gift, _ := strconv.ParseInt(v["sum(lucky_gift)"].(string), 10, 64)
			var levels []Level
			_, _ = o.QueryTable(new(Level)).OrderBy("-TotalBet").All(&levels)
			var llevel int
			//获取当前的VIP等级
			for _, l := range levels {
				if totalbet >= l.TotalBet {
					llevel = l.VipLevel
					break
				}
			}
			mt.CreateDate = time.Now()
			mt.ModifyDate = time.Now()
			mt.TotalLevelGift = level
			mt.TotalLuckyGift = gift
			mt.Account = v["account"].(string)
			mt.Bet = totalbet
			mt.Level = llevel
			mts = append(mts, mt)
		}
		var susNums int64
		// 将数组拆分导入，一次1000条
		mlen := len(mts)
		for i := 0; i <= mlen/1000; i++ {
			end := 0
			if (i+1)*1000 >= mlen {
				end = mlen
			} else {
				end = (i + 1) * 1000
			}
			tmpArr := mts[i*1000 : end]

			if nums, err := o.InsertMulti(len(tmpArr), tmpArr); err != nil {
				beego.Error("member import, insert error", err)
			} else {
				susNums += nums
			}
		}
		data["number"] = susNums
	}
	data["erraccount"] = erraccount
	data["code"] = 1
	data["msg"] = "计算成功"
	this.Data["json"] = &data
	this.ServeJSON()
}
