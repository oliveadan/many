package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type CheckRole struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now_add;type(datetime)"` // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	AdminId    int64     // 管理账号Id
	Hierarchy  string    // 分层
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(CheckRole))
}

//获取分层
func GetHierarchy(id int64) string {
	var checkrole CheckRole
	o := orm.NewOrm()
	o.QueryTable(new(CheckRole)).Filter("AdminId", id).One(&checkrole, "Hierarchy")
	return checkrole.Hierarchy
}

func (model *CheckRole) ReadOrCreate(col1 string, cols ...string) (bool, int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.ReadOrCreate(model, col1, cols...)
}

func (model *CheckRole) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *CheckRole) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	//model.Version =
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *CheckRole) QueryAll() (list []CheckRole) {
	var checkroles []CheckRole
	o := orm.NewOrm()
	o.QueryTable(new(CheckRole)).All(&checkroles)
	return checkroles
}

func (model *CheckRole) Paginate(page int, limit int) (list []CheckRole, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	cond := orm.NewCondition()
	qs := o.QueryTable(new(CheckRole))
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
