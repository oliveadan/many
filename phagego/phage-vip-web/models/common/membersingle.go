package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type MemberSingle struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	Account    string    // 会员账号
	Bet        int64     // 投注金额
	LevelGift  int64     // 晋级彩金
	LuckyGift  int64     // 当月好运金
	PeriodName string    // 期数名称
	PeriodSeq  int       // 期数排序字段
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(MemberSingle))
}

func (model *MemberSingle) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}
func (model *MemberSingle) Upate(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model)
}
func (model *MemberSingle) Paginate(page int, limit int,account string,periodname string) (list []MemberSingle, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(MemberSingle))
	cond := orm.NewCondition()
	if account != ""{
		cond = cond.And("Account",account)
	}
	if periodname != ""{
		cond = cond.And("PeriodName",periodname)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-PeriodName")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
