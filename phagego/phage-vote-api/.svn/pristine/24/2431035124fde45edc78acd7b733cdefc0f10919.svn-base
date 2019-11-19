package votedetail

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"phagego/frameweb-v2/controllers/sysmanage"
	"phagego/phage-vote-api/models"
)

type IndexVoteDetailController struct {
	 sysmanage.BaseController
}

func (v *IndexVoteDetailController) Get()  {
    page, err := v.GetInt("p")
	if err != nil {
		page = 1
	}
    limit, _ :=  beego.AppConfig.Int("pagelimit")
    list, total := new(models.Vote).Paginate(page,limit)
    pagination.SetPaginator(v.Ctx,limit,total)
    v.Data["dataList"] = list
    v.TplName = "common/votedetail/index.html"
}

func (v *IndexVoteDetailController) DelBtch()  {
	var code int
	var msg  string
	defer sysmanage.Retjson(v.Ctx,&msg,&code)
	o := orm.NewOrm()
	num, err := o.QueryTable(new(models.Vote)).Filter("Id__gte",0).Delete()
	if err != nil {
		beego.Error("Delete batch vote error",err)
		msg = "删除失败"
		return
	}
	code = 1
	msg = fmt.Sprintf("成功删除%d条记录", num)
	return
}


