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
	. "phagego/phage-vip3-web/models/common"
	"strconv"
	"strings"
	"time"
)

type MembersingleIndexController struct {
	sysmanage.BaseController
}

func (this *MembersingleIndexController) Get() {

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
	periodName := strings.TrimSpace(this.GetString("PeriodName"))
	levelGift := strings.TrimSpace(this.GetString("levelgift"))
	//所有期数名称
	var period []Period
	o := orm.NewOrm()
	_, err1 := o.QueryTable(new(Period)).OrderBy("-Rank").All(&period, "PeriodName")
	if err1 != nil {
		beego.Error("获取所有期数名称失败", err1)
	} else {
		this.Data["periodNames"] = period
	}
	if len(period) != 0 {
		//第一次进入的时候使用最新的一期名称
		if period[0].PeriodName != "" && periodName == "" {
			periodName = period[0].PeriodName
		}
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(MemberSingle).Paginate(page, limit, account, periodName, levelGift)
	pagination.SetPaginator(this.Ctx, limit, total)

	this.Data["condArr"] = map[string]interface{}{"account": account,
		"memberSingleName": periodName,
		"LevelGift":        levelGift}
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
	membersingle := MemberSingle{PeriodName: periodname}
	o := orm.NewOrm()
	num, err1 := o.Delete(&membersingle, "PeriodName")
	if err1 != nil {
		beego.Error("删除会员单期投注失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = fmt.Sprintf("成功删除%d条数据", num)
	}
}

func (this *MembersingleIndexController) CountGift() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	periodNmae := this.GetString("PeriodName")
	o := orm.NewOrm()
	//不是最新的一期不能计算
	var period Period
	errr := o.QueryTable(new(Period)).OrderBy("-Rank").One(&period)
	if errr != nil {
		beego.Error("获取期数失败", errr)
		msg = "获取期数失败"
		return
	}
	if period.PeriodName != periodNmae {
		msg = "只能计算最新的一期"
		return
	}

	if periodNmae == "" {
		msg = "请选择要计算的期数"
		return
	}
	//获取VIP等级
	var levels []Level
	_, err := o.QueryTable(new(Level)).OrderBy("-TotalBet").All(&levels)
	if err != nil {
		beego.Error("获取VIP等级失败", err)
		msg = "VIP等级获取失败，请检查VIP等级配置"
		return
	}
	//获取最小的vip等级信息
	var level Level
	err2 := o.QueryTable(new(Level)).Filter("VipLevel", 1).One(&level)
	if err2 != nil {
		beego.Error("获取最小VIP等级失败", err)
		msg = "VIP等级获取失败，请检查VIP等级配置"
		return
	}

	//获取要计算的期数的数据
	var membersingles []MemberSingle
	_, err1 := o.QueryTable(new(MemberSingle)).Filter("PeriodName", periodNmae).Filter("EnAble", 0).Limit(-1).All(&membersingles)
	if err1 != nil {
		beego.Error("获取数据失败", err1)
		msg = "获取要计算的期数数据失败"
		return
	}
	if len(membersingles) == 0 {
		msg = "所有数据已经计算过了"
		return
	}
	var membertotals []MemberTotal
	var mt MemberTotal
	var model MemberTotal
	o.Begin()
	for _, v := range membersingles {
		//获取会员以前的信息
		err := o.QueryTable(new(MemberTotal)).Filter("Account", v.Account).One(&mt)
		//新导入的会员，（进行插入）
		if err != nil {
			//如果投注额未达到vip1的要求跳过本次循环
			if v.Bet < level.TotalBet {
				model.Level = 0
				model.Bet = v.Bet
				model.CreateDate = time.Now()
				model.ModifyDate = time.Now()
				model.Version = 0
				model.Creator = this.LoginAdminId
				model.Modifior = this.LoginAdminId
				model.Account = v.Account
				membertotals = append(membertotals, model)
				continue
			}
			for _, j := range levels {
				//获取会员总投注后，与vip等级进行匹配，获得当前vip等级
				if v.Bet >= j.TotalBet {
					//晋一个等级的情况
					if j.VipLevel-0 == 1 {
						model.Bet = v.Bet
						model.CreateDate = time.Now()
						model.ModifyDate = time.Now()
						model.Version = 0
						model.Creator = this.LoginAdminId
						model.Modifior = this.LoginAdminId
						model.Account = v.Account
						model.Level = j.VipLevel
						membertotals = append(membertotals, model)
						//更新当期会员的晋级彩金
						_, err = o.QueryTable(new(MemberSingle)).Filter("Id", v.Id).Update(orm.Params{"LevelGift": j.LevelGift})
						if err != nil {
							beego.Error(v.Account, "更新晋级彩金失败(1)", errr)
							msg = "更新晋级彩金失败(1)"
							o.Rollback()
							return
						}
						break
						//连续跳级的情况
					} else if j.VipLevel-0 >= 2 {
						model.Bet = v.Bet
						model.CreateDate = time.Now()
						model.ModifyDate = time.Now()
						model.Version = 0
						model.Creator = this.LoginAdminId
						model.Modifior = this.LoginAdminId
						model.Account = v.Account
						model.Level = j.VipLevel
						membertotals = append(membertotals, model)
						var sum int64
						for _, v := range levels {
							if v.VipLevel > j.VipLevel {
								continue
							}
							sum += v.LevelGift
						}
						//更新当期会员的晋级彩金
						_, err = o.QueryTable(new(MemberSingle)).Filter("Id", v.Id).Update(orm.Params{"LevelGift": sum})
						if err != nil {
							beego.Error(v.Account, "更新晋级彩金失败(2)", errr)
							msg = "更新晋级彩金失败(2)"
							o.Rollback()
							return
						}
						break
					}
				}
			}
			//老会员（直接更新已有的数据）
		} else {
			//如果投注额未达到vip1的要求跳过本次循环
			if mt.Bet+v.Bet < level.TotalBet {
				_, err = o.QueryTable(new(MemberTotal)).Filter("Account", v.Account).Update(orm.Params{
					"Bet": orm.ColValue(orm.ColAdd, v.Bet)})
				if err != nil {
					beego.Error("未晋级的情况更新失败", err)
					o.Rollback()
					msg = "系统异常0"
					return
				}
				continue
			}
			for _, j := range levels {
				//获取会员总投注后，与vip等级进行匹配，获得当前vip等级
				if mt.Bet+v.Bet >= j.TotalBet {
					//未晋级的情况
					if j.VipLevel == mt.Level {
						_, err = o.QueryTable(new(MemberTotal)).Filter("Id", mt.Id).Update(orm.Params{
							"Bet": orm.ColValue(orm.ColAdd, v.Bet)})
						if err != nil {
							beego.Error("未晋级的情况更新失败", err)
							o.Rollback()
							msg = "系统异常1"
							return
						}
						break
						//晋一个等级的情况
					} else if j.VipLevel-mt.Level == 1 {
						_, err = o.QueryTable(new(MemberTotal)).Filter("Id", mt.Id).Update(orm.Params{
							"Bet":   orm.ColValue(orm.ColAdd, v.Bet),
							"Level": j.VipLevel})
						if err != nil {
							beego.Error("晋一个等级的情况更新失败", err)
							o.Rollback()
							msg = "系统异常2"
							return
						}
						//更新当期会员的晋级彩金
						_, err = o.QueryTable(new(MemberSingle)).Filter("Id", v.Id).Update(orm.Params{"LevelGift": j.LevelGift})
						if err != nil {
							beego.Error(v.Account, "更新晋级彩金失败(3)", errr)
							msg = "更新晋级彩金失败(3)"
							o.Rollback()
							return
						}
						break
						//连续跳级的情况
					} else if j.VipLevel-mt.Level >= 2 {
						_, err = o.QueryTable(new(MemberTotal)).Filter("Id", mt.Id).Update(orm.Params{
							"Bet":   orm.ColValue(orm.ColAdd, v.Bet),
							"Level": j.VipLevel})
						if err != nil {
							beego.Error("连续跳级的情况更新失败", err)
							o.Rollback()
							msg = "系统异常4"
							return
						}
						var sum int64
						for _, v := range levels {
							if v.VipLevel > j.VipLevel || v.VipLevel <= mt.Level {
								continue
							}
							sum += v.LevelGift
						}
						//更新当期会员的晋级彩金
						_, err = o.QueryTable(new(MemberSingle)).Filter("Id", v.Id).Update(orm.Params{"LevelGift": sum})
						if err != nil {
							beego.Error(v.Account, "更新晋级彩金失败(4)", errr)
							msg = "更新晋级彩金失败(4)"
							o.Rollback()
							return
						}
						break
					}
				}
			}
		}
	}

	//在计算后生成会员统计列表

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
	_, _ = o.QueryTable(new(MemberSingle)).Filter("PeriodName", periodNmae).Update(orm.Params{"EnAble": 1})
	code = 1
	msg = "计算成功"
}

//计算好运金
func (this *MembersingleIndexController) CountLuckyGift() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	periodNmae := this.GetString("PeriodName")
	o := orm.NewOrm()
	//不是最新的一期不能计算
	var period Period
	errr := o.QueryTable(new(Period)).OrderBy("-Rank").One(&period)
	if errr != nil {
		beego.Error("获取期数失败", errr)
		msg = "获取期数失败"
		return
	}
	var membersingle MemberSingle
	eerr := o.QueryTable(new(MemberSingle)).Filter("PeriodName", periodNmae).One(&membersingle)
	if eerr != nil {
		beego.Error("获取是否计算失败", eerr)
		msg = "获取期数失败"
		return
	}
	if period.PeriodName != periodNmae {
		msg = "只能计算最新的一期"
		return
	}
	if membersingle.EnAble == 0 {
		msg = "请先计算晋级彩金后再计算好运金"
		return
	}
	//获取要计算的期数的数据
	var membersingles []MemberSingle
	_, err1 := o.QueryTable(new(MemberSingle)).Filter("PeriodName", periodNmae).Limit(-1).All(&membersingles)
	if err1 != nil {
		beego.Error("获取数据失败", err1)
		msg = "获取要计算的期数数据失败"
		return
	}
	var luckys []Lucky
	var lucky Lucky
	_, err10 := o.QueryTable(new(Lucky)).OrderBy("-LuckyGift").All(&luckys)
	err10 = o.QueryTable(new(Lucky)).OrderBy("LuckyGift").One(&lucky)
	if err10 != nil {
		msg = "获取好运金配置失败"
		return
	}
	var mt MemberTotal
	for _, v := range membersingles {
		err := o.QueryTable(new(MemberTotal)).Filter("Account", v.Account).One(&mt)
		if err != nil {
			beego.Error(v.Account, "获取会员总信息失败")
			msg = "获取会员总信息失败"
			return
		}
		if mt.Level < 1 || v.Bet < lucky.MonthBet {
			continue
		}
		o.Begin()
		for _, h := range luckys {
			if v.Bet >= h.MonthBet && h.MaxVipLevel >= mt.Level && mt.Level >= h.MinVipLevel {
				_, err := o.QueryTable(new(MemberSingle)).Filter("Id", v.Id).Update(orm.Params{"LuckyGift": h.LuckyGift})
				if err != nil {
					beego.Error(v.Account, "更新好运金失败", errr)
					msg = "更新好运金失败"
					o.Rollback()
					return
				}
				break
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
			bet1, err := strconv.ParseInt(bet, 10, 64)
			if err != nil {
				msg = fmt.Sprintf("%s第%d行投注金额格式错误<br>", msg, i+1)
			}
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

func (this *MembersingleIndexController) Export() {
	o := orm.NewOrm()
	var membersingle []MemberSingle
	periodname := this.GetString("PeriodName")
	_, err := o.QueryTable(new(MemberSingle)).Filter("PeriodName", periodname).Limit(-1).All(&membersingle)
	if err != nil {
		beego.Error("导出失败", err)
		return
	}

	xlxs := excelize.NewFile()
	xlxs.SetCellValue("Sheet1", "A1", "期数名称")
	xlxs.SetCellValue("Sheet1", "B1", "会员账号")
	xlxs.SetCellValue("Sheet1", "C1", "投注金额")
	xlxs.SetCellValue("Sheet1", "D1", "晋级彩金")
	xlxs.SetCellValue("Sheet1", "E1", "当天好运金")
	for i, value := range membersingle {
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
