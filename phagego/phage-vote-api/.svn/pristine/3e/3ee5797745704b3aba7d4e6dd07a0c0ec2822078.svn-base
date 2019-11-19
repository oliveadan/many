package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Vote struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64                                         // 创建人Id
	Modifior   int64                                         // 更新人Id
	Version    int                                           // 版本
	Ip         string                                        // IP地址
	Category   int                                           // 类别
}

func init()  {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"),new(Vote))
}

func (model *Vote)ReadOrCreate(col1 string,cols ...string)(bool,int64,error)  {
	 model.CreateDate = time.Now()
	 model.ModifyDate = time.Now()
	 model.Version = 0
	 o := orm.NewOrm()
	 return o.ReadOrCreate(model,col1,cols...)
}

func (model *Vote)Update(cols ...string) (int64,error)  {
	if cols != nil {
		cols = append(cols,"ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model,cols...)
}

func (model *Vote) Paginate (page int, limit int)(list []Vote,total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Vote))
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}



