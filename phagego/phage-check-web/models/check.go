package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Check struct {
	Id                 int64     `auto`                              // 自增主键
	CreateDate         time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate         time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator            int64     `orm:"null"`                        // 创建人Id
	Modifior           int64     `orm:"null"`                        // 更新人Id
	Version            int       `orm:"null"`                        // 版本
	Account            string    `orm:"null"`                        //会员账号
	Idd                string    `orm:"null"`                        //idd
	Name               string    `orm:"null"`                        //会员姓名
	Agent              string    `orm:"null"`                        //代理商
	LoginInformation   string    `orm:"null"`                        //登录信息
	RegisterDate       time.Time `orm:"null"`                        //注册时间
	Mobile             string    `orm:"null"`                        //手机号码
	Email              string    `orm:"null"`                        //电子邮箱
	Qq                 string    `orm:"null"`                        //QQ号码
	Wechat             string    `orm:"null"`                        //微信号码
	PasswordHint       string    `orm:"null"`                        //密码提示问题
	PasswordAnswer     string    `orm:"null"`                        //密码提示答案
	WithdrawalPassword string    `orm:"null"`                        //取款密码
	OpenBank           string    `orm:"null"`                        //开户银行
	BankAccount        string    `orm:"null"`                        //银行账号
	Hierarchy          string    `orm:"null"`                        //分层
	DialStatus         int64     `orm:"null"`                        //拨打情况
	Comment            string    `orm:"null"`                        //备注
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(Check))
}

//获取idd
func GetIdd() []Check {
	var checks []Check
	o := orm.NewOrm()
	o.QueryTable(new(Check)).Distinct().All(&checks, "Idd")
	return checks
}

//获取分层
func GetHierarchys() []Check {
	var checks []Check
	o := orm.NewOrm()
	o.QueryTable(new(Check)).Distinct().All(&checks, "Hierarchy")
	return checks
}

func (model *Check) ReadOrCreate(col1 string, cols ...string) (bool, int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.ReadOrCreate(model, col1, cols...)
}

func (model *Check) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *Check) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	//model.Version =
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *Check) Paginate(page int, limit int, account string, idd string, name string, agent string, logininformation string, registerdate string, mobile string, email string, qq string, wechat string, passwordhint string, passwordanswer string, withdrawalpassword string, openbank string, bankaccount string, hierarchy string, adminroleid int64, loginadminid int64) (list []Check, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	cond := orm.NewCondition()
	qs := o.QueryTable(new(Check))
	if account != "" {
		cond = cond.And("Account", account)
	}
	if idd != "" {
		cond = cond.And("Idd", idd)
	}
	if name != "" {
		cond = cond.And("Name", name)
	}

	if agent != "" {
		cond = cond.And("Agent", agent)
	}
	if logininformation != "" {
		cond = cond.And("LoginInformation", logininformation)
	}
	if registerdate != "" {
		cond = cond.And("RegisterDate__exact", registerdate)
	}
	if mobile != "" {
		cond = cond.And("Mobile", mobile)
	}
	if email != "" {
		cond = cond.And("Email", email)
	}
	if qq != "" {
		cond = cond.And("Qq", qq)
	}
	if wechat != "" {
		cond = cond.And("Wechat", wechat)
	}
	if passwordhint != "" {
		cond = cond.And("PasswordHint", passwordhint)
	}
	if passwordanswer != "" {
		cond = cond.And("PasswordAnswer", passwordanswer)
	}
	if withdrawalpassword != "" {
		cond = cond.And("WithdrawalPassword", withdrawalpassword)
	}
	if openbank != "" {
		cond = cond.And("OpenBank", openbank)
	}
	if bankaccount != "" {
		cond = cond.And("BankAccount", bankaccount)
	}
	if hierarchy != "" {
		cond = cond.And("Hierarchy", hierarchy)
	}
	if adminroleid == 3 {
		if account == "" {
			cond = cond.And("Account", account)
		}
		lv := GetHierarchy(loginadminid)
		cond = cond.And("Hierarchy", lv)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	/*colss  := make([]string,0)
	//如果属于客服级，过滤返回字段
	if adminroleid == 3{
		qs.All(&list,colss...)
	}else {
		qs.All(&list)
	}*/
	qs.All(&list)
	total, _ = qs.Count()
	return
}
