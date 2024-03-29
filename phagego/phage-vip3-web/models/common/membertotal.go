package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type MemberTotal struct {
	Id             int64     `auto`                              // 自增主键
	CreateDate     time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate     time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator        int64     // 创建人Id
	Modifior       int64     // 更新人Id
	Version        int       // 版本
	Account        string    // 会员账号
	Level          int       // vip等级
	Bet            int64     // 总投注额
	TotalLevelGift int64     // 晋级金总额
	TotalLuckyGift int64     // 好运金总额
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(MemberTotal))
}

func (model *MemberTotal) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *MemberTotal) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *MemberTotal) Paginate(page int, limit int, account string) (list []MemberTotal, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(MemberTotal))
	if account != "" {
		qs = qs.Filter("Account", account)
	}
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Level", "Bet")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
