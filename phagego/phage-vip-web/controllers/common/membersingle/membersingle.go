package membersingle

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"net/url"
	"os"
	"phagego/common/utils"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip-web/models/common"
	"strconv"
	"strings"
	"time"
)

type MembersingleIndexController struct {
	sysmanage.BaseController
}

func (this *MembersingleIndexController) Get() {

	//导出
	isExport, _ := this.GetInt("isExport",0)
	if isExport == 1 {
		this.Export()
		return
	}

	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	account := strings.TrimSpace(this.GetString("account"))
	periodName := strings.TrimSpace(this.GetString("PeriodName"))
	//所有期数名称
	var period []Period
	o := orm.NewOrm()
	_, err1 := o.QueryTable(new(Period)).OrderBy("-Rank").All(&period, "PeriodName")
	if err1 != nil {
		beego.Error("获取所有期数名称失败", err1)
	} else {
		this.Data["periodNames"] = period
	}
	//第一次进入的时候使用最新的一期名称
	if period[0].PeriodName != "" && periodName =="" {
		periodName = period[0].PeriodName
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(MemberSingle).Paginate(page, limit, account, periodName)
	pagination.SetPaginator(this.Ctx, limit, total)

	this.Data["condArr"] = map[string]interface{}{"account": account,
		"memberSingleName": periodName}
	this.Data["dataList"] = list
	this.TplName = "common/membersingle/index.html"
}

func (this *MembersingleIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	membersingle := MemberSingle{Id: id}
	o := orm.NewOrm()
	err := o.Read(&membersingle)
	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		this.Redirect("membersingleIndexController.get", 302)
	}
	_, err1 := o.Delete(&membersingle, "Id")
	if err1 != nil {
		beego.Error("删除会员单期投注失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (this *MembersingleIndexController) DelBatch() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	periodname := strings.TrimSpace(this.GetString("PeriodName"))
	if periodname == "" {
		msg = "请选择要删除的期数"
		return
	}
	membersingle := MemberSingle{PeriodName:periodname}
	o := orm.NewOrm()
	num, err1 := o.Delete(&membersingle, "PeriodName")
	if err1 != nil {
		beego.Error("删除会员单期投注失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = fmt.Sprintf("成功删除%d条数据",num)
	}
}

func (this *MembersingleIndexController) CountGift() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)

	periodNmae := this.GetString("PeriodName")
	if periodNmae == "" {
		msg = "请选择要计算的期数"
		return
	}
	o := orm.NewOrm()
	//获取VIP等级
	var levels []Level
	_, err := o.QueryTable(new(Level)).OrderBy("-TotalBet").All(&levels)
	if err != nil {
		beego.Error("获取VIP等级失败", err)
		msg = "VIP等级获取失败，请检查VIP等级配置"
		return
	}
	//获取要计算的期数的数据
	var membersingles []MemberSingle
	_, err1 := o.QueryTable(new(MemberSingle)).Filter("PeriodName", periodNmae).Limit(50000).All(&membersingles)
	if err1 != nil {
		beego.Error("获取数据失败", err1)
		msg = "获取要计算的期数数据失败"
		return
	}
	//获取好运金配置信息
	var luckys []Lucky
	_, err2 := o.QueryTable(new(Lucky)).OrderBy("-MonthBet").All(&luckys)
	if err2 != nil {
		beego.Error("获取好运金配置失败", err2)
		msg = "获取好运金配置失败"
		return
	}

	for _, v := range membersingles {
		for _, j := range levels {
			//获取会员本期和往期投注的所有之和
			var maps [] orm.Params
			_, err := o.Raw("SELECT sum(bet) FROM ph_member_single WHERE account =? and period_seq <= ?", v.Account, v.PeriodSeq).Values(&maps)
			if err != nil {
				beego.Error("会员账号",v.Account,"获取本期和往期投注之和失败", err)
				msg = "获取本期和往期投注之和失败"
				return
			}
			totalbet, _ := strconv.ParseInt(maps[0]["sum(bet)"].(string), 10, 64)
			//获取会员总投注后，与vip等级进行匹配，获得当前vip等级
			if totalbet >= j.TotalBet {
				//判断会员是否晋级
				var level MemberTotal
				err := o.QueryTable(new(MemberTotal)).Filter("Account", v.Account).One(&level, "Level")
				//连续跳级的情况
				if j.VipLevel-level.Level >= 2 {
					//获跳级级奖励
					var maps []orm.Params
					_, err := o.Raw("SELECT sum(level_gift) from ph_level WHERE vip_level >=? and vip_level <=?", level.Level+1, j.VipLevel).Values(&maps)
					if err != nil {
						beego.Error("会员账号",v.Account,"获取本期和往期投注之和失败", err)
						msg = "获取本期和往期投注之和失败"
						return
					}
					//更新跳级奖励
					upgift, _ := strconv.ParseInt(maps[0]["sum(level_gift)"].(string), 10, 64)
					_, _ = o.QueryTable(new(MemberSingle)).Filter("Account", v.Account).Filter("PeriodName", v.PeriodName).Update(orm.Params{"level_gift": upgift})
					//更新好运金
					for _, h := range luckys {
						if v.Bet >= h.MonthBet && h.MaxVipLevel >= j.VipLevel && j.VipLevel >= h.MinVipLevel {
							_, _ = o.QueryTable(new(MemberSingle)).Filter("Account", v.Account).Filter("PeriodName", v.PeriodName).Update(orm.Params{"LuckyGift": h.LuckyGift})
							break
						}
					}
					break
					//晋升一个等级的情况
				} else if j.VipLevel > level.Level || err != nil {
					_, _ = o.QueryTable(new(MemberSingle)).Filter("Account", v.Account).Filter("PeriodName", v.PeriodName).Update(orm.Params{"LevelGift": j.LevelGift})
					//更新好运金
					for _, h := range luckys {
						if v.Bet >= h.MonthBet && h.MaxVipLevel >= j.VipLevel && j.VipLevel >= h.MinVipLevel {
							_, _ = o.QueryTable(new(MemberSingle)).Filter("Account", v.Account).Filter("PeriodName", v.PeriodName).Update(orm.Params{"LuckyGift": h.LuckyGift})
							break
						}
					}
					break
					//未晋级的情况只更新好运金
				} else if j.VipLevel == level.Level {
					//更新好运金
					for _, h := range luckys {
						if v.Bet >= h.MonthBet && h.MaxVipLevel >= j.VipLevel && j.VipLevel >= h.MinVipLevel {
							_, _ = o.QueryTable(new(MemberSingle)).Filter("Account", v.Account).Filter("PeriodName", v.PeriodName).Update(orm.Params{"LuckyGift": h.LuckyGift})
							break
						}
					}
				}
			}
		}
	}

	//在计算后生成会员统计列表
	_, _ = o.Raw("DELETE from ph_member_total WHERE id != 0").Exec()
	var membertotals []MemberTotal
	var accouts []MemberSingle
	_, _ = o.QueryTable(new(MemberSingle)).Distinct().Limit(-1).All(&accouts, "Account")
	for _, v := range accouts {
		var model MemberTotal
		var maps [] orm.Params
		_, err := o.Raw("SELECT sum(bet),sum(level_gift),sum(lucky_gift) FROM ph_member_single where account = ?", v.Account).Values(&maps)
		if err != nil {
			beego.Error("获取本期和往期投注之和失败", err)
			msg = "获取本期和往期投注之和失败"
			return
		}
		totalbet, _ := strconv.ParseInt(maps[0]["sum(bet)"].(string), 10, 64)
		for _, j := range levels {
			if totalbet >= j.TotalBet {
				model.Level = j.VipLevel
				break
			}
		}

		model.Bet = totalbet
		levelgift, _ := strconv.ParseInt(maps[0]["sum(level_gift)"].(string), 10, 64)
		model.TotalLevelGift = levelgift
		luckygift, _ := strconv.ParseInt(maps[0]["sum(lucky_gift)"].(string), 10, 64)
		model.TotalLuckyGift = luckygift
		model.CreateDate = time.Now()
		model.ModifyDate = time.Now()
		model.Version = 0
		model.Creator = this.LoginAdminId
		model.Modifior = this.LoginAdminId
		model.Account = v.Account

		membertotals = append(membertotals, model)
	}
	o.Begin()
	var susNums int64
	//将数组拆分导入，一次1000条
	mlen := len(membertotals)
	if mlen > 0 {
		for i := 0; i <= mlen/1000; i++ {
			end := 0
			if (i+1)*1000 >= mlen {
				end = mlen
			} else {
				end = (i + 1) * 1000
			}
			if i*1000 == end {
				continue
			}
			tmpArr := membertotals[i*1000 : end]
			if nums, err := o.InsertMulti(len(tmpArr), tmpArr); err != nil {
				o.Rollback()
				beego.Error("插入会员总投注失败", err)
				return
			} else {
				susNums += nums
			}
		}
	}
	o.Commit()
	code = 1
	msg = "计算成功"
}

func (this *MembersingleIndexController) Import() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	f, h, err := this.GetFile("file")
	defer f.Close()
	if err != nil {
		beego.Error("导入会员投注失败", err)
		msg = "导入失败，请重试（1）"
		return
	}
	fname := url.QueryEscape(h.Filename)
	suffix := utils.SubString(fname, len(fname), strings.LastIndex(fname, ".")-len(fname))
	if suffix != ".xlsx" {
		msg = "文件必须为xlsx"
		return
	}

	o := orm.NewOrm()
	membersingles := make([]MemberSingle, 0)

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		beego.Error("会员投注导入失败", err)
		msg = "读取excel失败，请重试"
		return
	}
	if xlsx.GetSheetIndex("会员投注") == 0 {
		msg = "不存在《会员投注》sheet页"
		return
	}
	rows := xlsx.GetRows("会员投注")
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			msg = fmt.Sprintf("%s第%d行账号为空<br>", msg, i+1)
			continue
		}
		var membersingle MemberSingle
		periodname := strings.TrimSpace(row[0])
		bool := o.QueryTable(new(Period)).Filter("PeriodName", periodname).Exist()
		if !bool {
			msg = fmt.Sprintf("%s第%d行期数名称不存在<br>", msg, i+1)
			return
		}
		account := strings.TrimSpace(row[1])
		if account == "" {
			msg = fmt.Sprintf("%s第%d行会员账号为空<br>", msg, i+1)
		}
		bet := strings.TrimSpace(row[2])
		if bet == "" {
			msg = fmt.Sprintf("%s第%d行投注金额为空<br>", msg, i+1)
		} else {
			bet1, _ := strconv.ParseInt(bet, 10, 64)
			membersingle.Bet = bet1
		}
		membersingle.PeriodName = periodname
		membersingle.Account = account
        var periodseq Period
		_ = o.QueryTable(new(Period)).Filter("PeriodName", periodname).One(&periodseq, "Rank")
		membersingle.PeriodSeq = periodseq.Rank
		//当会员账号已经存在,更新投注金额
		if account != "" && bool {
			bool1 := o.QueryTable(new(MemberSingle)).Filter("Account", account).Filter("PeriodName", periodname).Exist()
			if bool1 && bet != "" {
				_, err := o.QueryTable(new(MemberSingle)).Filter("Account", account).Filter("PeriodName", periodname).Update(orm.Params{"Bet": bet})
				if err != nil {
					beego.Error("更新已存在会员的投注额失败", err)
					msg = fmt.Sprintf("%s第%d行更新已存在会员的投注额失败<br>", msg, i+1)
				} else {
					continue
				}
			}
		}

		membersingles = append(membersingles, membersingle)
	}
	if msg != "" {
		msg = fmt.Sprintf("请处理以下错误后再导入：<br>%s", msg)
		return
	}
	rlen := len(membersingles)
	if rlen == 0 {
		msg = "没有需要导入的数据"
		return
	}
	var susNums int64
	// 将数组拆分导入，一次1000条
	for i := 0; i <= rlen/1000; i++ {
		end := 0
		if (i+1)*1000 >= rlen {
			end = rlen
		} else {
			end = (i + 1) * 1000
		}
		if i*1000 == end {
			continue
		}
		tmpArr := membersingles[i*1000 : end]
		if nums, err := o.InsertMulti(len(tmpArr), tmpArr); err != nil {
			beego.Error("会员投注记录导入失败", err)
		} else {
			susNums += nums
		}
	}
	code = 1
	msg = fmt.Sprintf("%s成功导入%d条记录", msg, susNums)
	return
}

func(this *MembersingleIndexController) Export(){
	o := orm.NewOrm()
	var membersingle []MemberSingle
	periodname := this.GetString("PeriodName")
	_, err := o.QueryTable(new(MemberSingle)).Filter("PeriodName",periodname).Limit(-1).All(&membersingle)
	if err != nil {
		beego.Error("导出失败",err)
		return
	}

	xlxs := excelize.NewFile()
	xlxs.SetCellValue("Sheet1","A1","期数名称")
	xlxs.SetCellValue("Sheet1","B1","会员账号")
	xlxs.SetCellValue("Sheet1","C1","投注金额")
	xlxs.SetCellValue("Sheet1","D1","晋级彩金")
	xlxs.SetCellValue("Sheet1","E1","当月好运金")
	for i,value := range membersingle  {
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), value.PeriodName)
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), value.Account)
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), value.Bet)
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), value.LevelGift)
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), value.LuckyGift)
	}
	fileName := fmt.Sprintf("./tmp/excel/membersinglelist_%s.xlsx", time.Now().Format("20060102150405"))
	err1 := xlxs.SaveAs(fileName)
	if err1 != nil {
		beego.Error("Export membersinglelist_ error", err.Error())
	} else {
		defer os.Remove(fileName)
		this.Ctx.Output.Download(fileName)
	}
}

type MembersingleAddController struct {
	sysmanage.BaseController
}

func (this *MembersingleAddController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "common/membersingle/add.html"
}

func (this *MembersingleAddController) Post() {
	var code int
	var msg string
	url := beego.URLFor("membersingleIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	membersingle := MemberSingle{}
	if err := this.ParseForm(&membersingle); err != nil {
		beego.Error("会员单期投注参数异常", err)
		msg = "参数异常"
		return
	}
	o := orm.NewOrm()
	_, err1 := o.Insert(&membersingle)
	if err1 != nil {
		beego.Error("添加会员单期投注失败", err1)
		msg = "添加失败"
		return
	} else {
		code = 1
		msg = "添加成功"
	}
}

type MembersingleEditController struct {
	sysmanage.BaseController
}

func (this *MembersingleEditController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	membersingle := MemberSingle{Id: id}
	err := o.Read(&membersingle)
	if err != nil {
		this.Redirect("membersingleIndexController.get", 302)
	} else {
		this.Data["data"] = membersingle
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/membersingle/edit.html"
	}
}

func (this *MembersingleEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("MembersingleIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	membersingle := MemberSingle{}
	if err := this.ParseForm(&membersingle); err != nil {
		beego.Error("修改单期投注参数异常", err)
		msg = "参数异常"
		return
	}
	cols := []string{"Account", "Bet", "LevelGift", "LuckyGift"}
	membersingle.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	_, err := o.Update(&membersingle, cols...)
	if err != nil {
		beego.Error("更新会员单期投注失败", err)
		msg = "更新失败"
		return
	} else {
		code = 1
		msg = "更新成功"
	}
}
