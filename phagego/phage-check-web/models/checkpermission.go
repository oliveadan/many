package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CheckPermission struct {
	Id                 int64     `auto`                              // 自增主键
	CreateDate         time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate         time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator            int64     // 创建人Id
	Modifior           int64     // 更新人Id
	Version            int       // 版本
	Account            int64     //会员账号
	Idd                int64     //Idd
	Name               int64     //会员姓名
	Agent              int64     //代理商
	LoginInformation   int64     //登录信息
	RegisterDate       int64     //注册时间
	Mobile             int64     //手机号码
	Email              int64     //电子邮箱
	Qq                 int64     //QQ号码
	Wechat             int64     //微信号码
	PasswordHint       int64     //密码提示问题
	PasswordAnswer     int64     //密码提示答案
	WithdrawalPassword int64     //取款密码
	OpenBank           int64     //开户银行
	BankAccount        int64     //银行账号
	AdminId            int64     //管理员Id
	Hierarchy          int64     // 分层
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(CheckPermission))
}

func GetCheckPermission(id int64) CheckPermission {
	var checkpermission CheckPermission
	o := orm.NewOrm()
	o.QueryTable(new(CheckPermission)).Filter("AdminId", id).One(&checkpermission)
	return checkpermission
}
func (model *CheckPermission) ReadOrCreate(col1 string, cols ...string) (bool, int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.ReadOrCreate(model, col1, cols...)
}

func (model *CheckPermission) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *CheckPermission) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	//model.Version =
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *CheckPermission) Paginate(page int, limit int, param1 string, status int) (list []CheckPermission, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(CheckPermission))
	cond := orm.NewCondition()
	if status != -1 {
		cond = cond.And("Status", status)
	}
	if param1 != "" {
		cond = cond.And("Account__contains", param1)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
