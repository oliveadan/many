package check

import (
	"net/url"
	"phage/utils"
	"phagego/phagev2/controllers/sysmanage"
	"strings"

	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"math"
	"os"
	"phagego/phagev2/models"
	. "phagego/phage-check-web/models"
	"strconv"
	"time"
)

type CheckIndexController struct {
	sysmanage.BaseController
}

func (this *CheckIndexController) Prepare() {
	this.EnableXSRF = false
}

func (this *CheckIndexController) Get() {
	beego.Informational("query check")

	//查看后台管理账号是否属于客服管理组
	loginadminid := this.GetSession("loginAdminId")
	var adminrole models.AdminRole
	o := orm.NewOrm()
	_, err1 := o.QueryTable(new(models.AdminRole)).Filter("AdminId", loginadminid).All(&adminrole)
	if err1 != nil {
		beego.Error("查询管理员关联角色关联失败", err1)
	}

	//导出
	isExport, _ := this.GetInt("isExport", 0)
	if isExport == 1 {
		this.Export(adminrole.RoleId, loginadminid.(int64))
		return
	}

	account := strings.TrimSpace(this.GetString("account"))
	idd := strings.TrimSpace(this.GetString("idd"))
	name := strings.TrimSpace(this.GetString("name"))
	agent := strings.TrimSpace(this.GetString("agent"))
	loginInformation := strings.TrimSpace(this.GetString("loginInformation"))
	registerDate := strings.TrimSpace(this.GetString("registerdate"))
	mobile := strings.TrimSpace(this.GetString("mobile"))
	email := strings.TrimSpace(this.GetString("email"))
	qq := strings.TrimSpace(this.GetString("qq"))
	wechat := strings.TrimSpace(this.GetString("wechat"))
	passwordHint := strings.TrimSpace(this.GetString("passwordhint"))
	passwordAnswer := strings.TrimSpace(this.GetString("passwordanswer"))
	withdrawalPassword := strings.TrimSpace(this.GetString("withdrawalpassword"))
	openBank := strings.TrimSpace(this.GetString("openbank"))
	bankAccount := strings.TrimSpace(this.GetString("bankaccount"))
	hierarchy := strings.TrimSpace(this.GetString("hierarchy"))
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")

	list, total := new(Check).Paginate(page, limit, account, idd, name, agent, loginInformation, registerDate, mobile, email, qq, wechat, passwordHint, passwordAnswer, withdrawalPassword, openBank, bankAccount, hierarchy, adminrole.RoleId, loginadminid.(int64))
	pagination.SetPaginator(this.Ctx, limit, total)
	//返回值
	this.Data["condArr"] = map[string]interface{}{"account": account,
		"name":               name,
		"agent":              agent,
		"loginInformation":   loginInformation,
		"registerDate":       registerDate,
		"mobile":             mobile,
		"email":              email,
		"qq":                 qq,
		"wechat":             wechat,
		"passwordHint":       passwordHint,
		"withdrawalPassword": withdrawalPassword,
		"openbank":           openBank,
		"bankaccount":        bankAccount,
		"hierarchy":          hierarchy,
		"idd":                idd}
	//获取字段表
	checkpermission := GetCheckPermission(loginadminid.(int64))
	hierarchys := GetHierarchys()
	//获取所有idd
	//idds := GetIdd()
	//this.Data["idds"] = idds
	this.Data["Hierarchys"] = hierarchys
	this.Data["checkpermission"] = checkpermission
	this.Data["adminroleid"] = adminrole.RoleId
	this.Data["dataList"] = list
	this.TplName = "common/check/index.html"
}

func (this *CheckIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	check := Check{Id: id}
	o := orm.NewOrm()
	_, err := o.Delete(&check, "Id")
	if err != nil {
		beego.Error("Delete member error", err)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (this *CheckIndexController) Delbatch() {
	beego.Informational("批量删除")
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	o := orm.NewOrm()
	hierarchy := this.GetString("hierarchy")
	if hierarchy == "全部" {
		res, err := o.Raw("DELETE from ph_check WHERE 1=1").Exec()
		if err != nil {
			beego.Error("Delete batch memberLottery error", err)
			msg = "删除失败"
		} else {
			code = 1
			num, _ := res.RowsAffected()
			msg = fmt.Sprintf("成功删除%d条记录", num)
		}
	} else {
		i, err := o.QueryTable(new(Check)).Filter("Hierarchy", hierarchy).Delete()
		if err != nil {
			beego.Error("批量删除错误", err)
			msg = "删除失败"
		} else {
			code = 1
			msg = fmt.Sprintf("成功删除%d条记录", i)
		}
	}
}

func (this *CheckIndexController) Import() {
	var code int
	var msg string
	var url1 string
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url1)
	f, h, err := this.GetFile("file")
	defer f.Close()
	if err != nil {
		beego.Error("check upload file get file error", err)
		msg = "上传失败，请重试（1）"
		return
	}
	fname := url.QueryEscape(h.Filename)
	suffix := utils.SubString(fname, len(fname), strings.LastIndex(fname, ".")-len(fname))
	if suffix != ".xlsx" {
		msg = "文件必须为 xlsx"
		return
	}

	o := orm.NewOrm()
	models := make([]Check, 0)
	idsDel := make([]int64, 0)

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		beego.Error("Check import, open excel error", err)
		msg = "读取excel失败，请重试"
		return
	}

	if xlsx.GetSheetIndex("会员信息") == 0 {
		msg = "不存在<<会员信息>>页脚，请重新设置"
	}
	rows := xlsx.GetRows("会员信息")

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 14 {
			msg = fmt.Sprintf("%s第%行，有空白项<br>", msg, i+1)
			continue
		}
		Account := strings.TrimSpace(row[0])
		if Account == "" {
			msg = fmt.Sprintf("%s第%d行会员账号不能为空，请检查<br>", msg, i+1)
			continue
		}
		Name := strings.TrimSpace(row[1])
		Agent := strings.TrimSpace(row[2])
		LoginInformation := strings.TrimSpace(row[3])
		//RegisterDate := strings.TrimSpace(row[4])
		Mobile := strings.TrimSpace(row[5])
		Email := strings.TrimSpace(row[6])
		Qq := strings.TrimSpace(row[7])
		Wechat := strings.TrimSpace(row[8])
		PasswordHint := strings.TrimSpace(row[9])
		PasswordAnswer := strings.TrimSpace(row[10])
		WithdrawalPassword := strings.TrimSpace(row[11])
		OpenBank := strings.TrimSpace(row[12])
		BankAccount := strings.TrimSpace(row[13])
		hierarchy := strings.TrimSpace(row[14])
		idd := strings.TrimSpace(row[15])
		//根据会员账号判断是否已经存在

		var tmpModel Check
		err1 := o.QueryTable(new(Check)).Filter("Account", Account).One(&tmpModel)
		if err != nil && err1 != orm.ErrNoRows {
			msg = fmt.Sprintf("%s第%d行数据中有错误，请重试<br>", msg, i+1)
			continue
		}
		//如果些会员存在，收集其ID后面进行删除，些模式为覆盖导入
		if err == nil {
			idsDel = append(idsDel, tmpModel.Id)
		}
		model := Check{}
		model.Account = Account
		model.Name = Name
		model.Agent = Agent
		model.LoginInformation = LoginInformation
		v1, _ := strconv.ParseFloat(row[4], 64)
		newtime := timeFromExcelTime(v1, false)
		model.RegisterDate = newtime
		model.Mobile = Mobile
		model.Email = Email
		model.Qq = Qq
		model.Wechat = Wechat
		model.PasswordHint = PasswordHint
		model.PasswordAnswer = PasswordAnswer
		model.WithdrawalPassword = WithdrawalPassword
		model.OpenBank = OpenBank
		model.BankAccount = BankAccount
		model.Hierarchy = hierarchy
		model.Idd = idd
		model.Creator = this.LoginAdminId
		model.Modifior = this.LoginAdminId
		model.CreateDate = time.Now()
		model.ModifyDate = time.Now()
		model.Version = 0
		models = append(models, model)
	}
	if msg != "" {
		msg = fmt.Sprintf("请处理以下错误后再导入：<br>%s", msg)
		return
	}
	if len(models) == 0 {
		msg = "导入表格为空，请确认"
		return
	}
	o.Begin()
	if len(idsDel) > 0 {
		idslen := len(idsDel)
		for i := 0; i <= idslen/1000; i++ {
			end := 0
			if (i+1)*1000 >= idslen {
				end = idslen
			} else {
				end = (i + 1) * 1000
			}
			tmpArr := idsDel[i*1000 : end]

			if _, err = o.QueryTable(new(Check)).Filter("Id__in", tmpArr).Delete(); err != nil {
				o.Rollback()
				msg = "导入失败，请重试"
				return
			}
		}
	}
	var susNums int64
	//将数组拆分导入，一次1000条
	mlen := len(models)
	for i := 0; i <= mlen/10; i++ {
		end := 0
		if (i+1)*10 >= mlen {
			end = mlen
		} else {
			end = (i + 1) * 10
		}
		if i*10 == end {
			continue
		}
		tmpArr := models[i*10 : end]
		if nums, err := o.InsertMulti(len(tmpArr), tmpArr); err != nil {
			o.Rollback()
			beego.Error("Check import,insert error", err)
			msg = "导入失败，请重试（2）"
			return
		} else {
			susNums += nums
		}
	}
	o.Commit()
	code = 1
	msg = fmt.Sprintf("成功导入%d条记录", susNums)
	return
}

func (this *CheckIndexController) Export(adminroleid int64, loginadminid int64) {
	beego.Informational("export check")
	//条件和搜索一致
	account := strings.TrimSpace(this.GetString("account"))
	idd := strings.TrimSpace(this.GetString("Idd"))
	name := strings.TrimSpace(this.GetString("name"))
	agent := strings.TrimSpace(this.GetString("agent"))
	loginInformation := strings.TrimSpace(this.GetString("loginInformation"))
	registerDate := strings.TrimSpace(this.GetString("registerdate"))
	mobile := strings.TrimSpace(this.GetString("mobile"))
	email := strings.TrimSpace(this.GetString("email"))
	qq := strings.TrimSpace(this.GetString("qq"))
	wechat := strings.TrimSpace(this.GetString("wechat"))
	passwordHint := strings.TrimSpace(this.GetString("passwordhint"))
	passwordAnswer := strings.TrimSpace(this.GetString("passwordanswer"))
	withdrawalPassword := strings.TrimSpace(this.GetString("withdrawalpassword"))
	openBank := strings.TrimSpace(this.GetString("openbank"))
	bankAccount := strings.TrimSpace(this.GetString("bankaccount"))
	hierarchy := strings.TrimSpace(this.GetString("hierarchy"))

	page := 1
	limit := 1000
	list, total := new(Check).Paginate(page, limit, account, idd, name, agent, loginInformation, registerDate, mobile, email, qq, wechat, passwordHint, passwordAnswer, withdrawalPassword, openBank, bankAccount, hierarchy, adminroleid, loginadminid)
	totalInt := int(total)
	if totalInt > limit {
		page1 := (float64(totalInt) - float64(limit)) / float64(limit)
		page2 := int(math.Ceil(page1))
		for page = 2; page <= (page2 + 1); page++ {
			list1, _ := new(Check).Paginate(page, limit, account, idd, name, agent, loginInformation, registerDate, mobile, email, qq, wechat, passwordHint, passwordAnswer, withdrawalPassword, openBank, bankAccount, hierarchy, adminroleid, loginadminid)
			for _, v := range list1 {
				list = append(list, v)
			}
		}
	}
	xlsx := excelize.NewFile()
	if adminroleid == 3 {
		checkpermission := GetCheckPermission(loginadminid)
		if checkpermission.Account == 1 {
			xlsx.SetCellValue("Sheet1", "B1", "会员账号")
		}
		if checkpermission.Name == 1 {
			xlsx.SetCellValue("Sheet1", "C1", "会员姓名")
		}
		if checkpermission.Agent == 1 {
			xlsx.SetCellValue("Sheet1", "D1", "代理商")
		}
		if checkpermission.LoginInformation == 1 {
			xlsx.SetCellValue("Sheet1", "E1", "登录信息")
		}
		if checkpermission.RegisterDate == 1 {
			xlsx.SetCellValue("Sheet1", "F1", "注册时间")
		}
		if checkpermission.Mobile == 1 {
			xlsx.SetCellValue("Sheet1", "G1", "手机号码")
		}
		if checkpermission.Email == 1 {
			xlsx.SetCellValue("Sheet1", "H1", "电子邮箱")
		}
		if checkpermission.Qq == 1 {
			xlsx.SetCellValue("Sheet1", "I1", "QQ号码")
		}
		if checkpermission.Wechat == 1 {
			xlsx.SetCellValue("Sheet1", "J1", "微信号码")
		}
		if checkpermission.PasswordHint == 1 {
			xlsx.SetCellValue("Sheet1", "K1", "密码提示问题")
		}
		if checkpermission.PasswordAnswer == 1 {
			xlsx.SetCellValue("Sheet1", "L1", "密码提示答案")
		}
		if checkpermission.WithdrawalPassword == 1 {
			xlsx.SetCellValue("Sheet1", "M1", "取款密码")
		}
		if checkpermission.OpenBank == 1 {
			xlsx.SetCellValue("Sheet1", "N1", "开户银行")
		}
		if checkpermission.BankAccount == 1 {
			xlsx.SetCellValue("Sheet1", "O1", "银行账户")
		}

		for i, value := range list {
			checkpermission := GetCheckPermission(loginadminid)
			if checkpermission.Account == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), value.Account)
			}
			if checkpermission.Name == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), value.Name)
			}
			if checkpermission.Agent == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), value.Agent)
			}
			if checkpermission.LoginInformation == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), value.LoginInformation)
			}
			if checkpermission.RegisterDate == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), value.RegisterDate.Format("2006-01-02 15:04:05"))
			}
			if checkpermission.Mobile == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("G%d", i+2), value.Mobile)
			}
			if checkpermission.Email == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d", i+2), value.Email)
			}
			if checkpermission.Qq == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("I%d", i+2), value.Qq)
			}
			if checkpermission.Wechat == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("J%d", i+2), value.Wechat)
			}
			if checkpermission.PasswordHint == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("K%d", i+2), value.PasswordHint)
			}
			if checkpermission.PasswordAnswer == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("L%d", i+2), value.PasswordAnswer)
			}
			if checkpermission.WithdrawalPassword == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("M%d", i+2), value.WithdrawalPassword)
			}
			if checkpermission.OpenBank == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("N%d", i+2), value.OpenBank)
			}
			if checkpermission.BankAccount == 1 {
				xlsx.SetCellValue("Sheet1", fmt.Sprintf("O%d", i+2), value.BankAccount)
			}
		}

	} else {
		m := map[int64]string{
			0:  "未拨打",
			1:  "有效接听",
			2:  "不是本人",
			3:  "接通挂断",
			4:  "无法接通",
			5:  "无人接听",
			6:  "拒接",
			7:  "关机",
			8:  "停机",
			9:  "空号",
			10: "其他",
		}
		xlsx.SetCellValue("Sheet1", "A1", "分层")
		xlsx.SetCellValue("Sheet1", "B1", "ID")
		xlsx.SetCellValue("Sheet1", "CV1", "会员账号")
		xlsx.SetCellValue("Sheet1", "D1", "会员姓名")
		xlsx.SetCellValue("Sheet1", "E1", "代理商")
		xlsx.SetCellValue("Sheet1", "F1", "登录信息")
		xlsx.SetCellValue("Sheet1", "G1", "注册时间")
		xlsx.SetCellValue("Sheet1", "H1", "手机号码")
		xlsx.SetCellValue("Sheet1", "I1", "电子邮箱")
		xlsx.SetCellValue("Sheet1", "J1", "QQ号码")
		xlsx.SetCellValue("Sheet1", "K1", "微信号码")
		xlsx.SetCellValue("Sheet1", "L1", "密码提示问题")
		xlsx.SetCellValue("Sheet1", "M1", "密码提示答案")
		xlsx.SetCellValue("Sheet1", "N1", "取款密码")
		xlsx.SetCellValue("Sheet1", "O1", "开户银行")
		xlsx.SetCellValue("Sheet1", "P1", "银行账户")
		xlsx.SetCellValue("Sheet1", "Q1", "拨打情况")
		xlsx.SetCellValue("Sheet1", "R1", "备注")
		for i, value := range list {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), value.Hierarchy)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), value.Idd)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), value.Account)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), value.Name)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), value.Agent)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), value.LoginInformation)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("G%d", i+2), value.RegisterDate.Format("2006-01-02 15:04:05"))
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d", i+2), value.Mobile)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("I%d", i+2), value.Email)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("J%d", i+2), value.Qq)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("K%d", i+2), value.Wechat)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("L%d", i+2), value.PasswordHint)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("M%d", i+2), value.PasswordAnswer)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("N%d", i+2), value.WithdrawalPassword)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("O%d", i+2), value.OpenBank)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("P%d", i+2), value.BankAccount)
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("Q%d", i+2), m[value.DialStatus])
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("R%d", i+2), value.Comment)
		}
	}
	fileName := fmt.Sprintf("./tmp/excel/rewardlist_%s.xlsx", time.Now().Format("20060102150405"))
	err := xlsx.SaveAs(fileName)
	if err != nil {
		beego.Error("导出会员信息失败", err.Error())
	} else {
		defer os.Remove(fileName)
		this.Ctx.Output.Download(fileName)
	}
}

type CheckEditController struct {
	sysmanage.BaseController
}

func (this *CheckEditController) Get() {
	Id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	check := Check{Id: Id}
	err := o.Read(&check)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("CheckIndexController.get"), 302)
	} else {
		this.Data["data"] = check
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/check/edit.html"
	}
}

func (this *CheckEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("CheckIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	check := Check{}
	if err := this.ParseForm(&check); err != nil {
		msg = "参数异常"
	}
	cols := []string{"DialStatus", "Comment"}
	_, err := check.Update(cols...)
	if err != nil {
		msg = "更新失败"
		beego.Error("更新查看会员信息更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
