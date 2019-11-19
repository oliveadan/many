package voteapi

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"phagego/frameweb-v2/controllers/sysmanage"
	"phagego/phage-vote-api/models"
)

type IndexVoteApiIndexController struct {
	beego.Controller
}

func (this *IndexVoteApiIndexController) Prepare() {
	this.EnableXSRF = false
}

func (this *IndexVoteApiIndexController) Get() {
	this.TplName = "front/index.html"
}

func (this *IndexVoteApiIndexController) Vote() {
	this.AllowCross()
	var code int
	var msg string
	data := make(map[string]interface{})
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &data)
	category, _ := this.GetInt("vote")
	if category == 0 {
		msg = "请选择一个选项"
		return
	}
	ip := this.Ctx.Input.IP()
	var v models.Vote
	/*sprintf := fmt.Sprintf("%v", time.Now().Unix())*/
	v.Ip = ip /*+sprintf*/
	v.Category = category
	bool, _, err := v.ReadOrCreate("Ip")
	if err != nil {
		beego.Error("vote error", err)
		msg = "系统异常,请刷新后重试"
		return
	}
	o := orm.NewOrm()
	var vv []models.Vote
	_, e := o.QueryTable(new(models.Vote)).Limit(-1).All(&vv)
	if e != nil {
		beego.Error("query vote error", e)
		msg = "系统异常(1)"
		return
	}
	var votes []int64
	for i := 1; i < 13; i++ {
		count, _ := o.QueryTable(new(models.Vote)).Filter("Category", i).Count()
		var sv models.SetVote
		o.QueryTable(new(models.SetVote)).Filter("Category", i).One(&sv)
		sum := count + sv.Value
		votes = append(votes, sum)
	}
	code = 1
	data["votes"] = votes
	if bool {
		msg = "恭喜您,投票成功"
	} else {
		msg = "您已经投过票了"
	}
}

func (this *IndexVoteApiIndexController) AllVote() {
	this.AllowCross()
	var code int
	var msg string
	data := make(map[string]interface{})
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &data)
	ip := this.Ctx.Input.IP()
	o := orm.NewOrm()
	exist := o.QueryTable(new(models.Vote)).Filter("Ip", ip).Exist()
	if !exist {
		msg = "请投票过后再来查看"
		return
	}
	var v []models.Vote
	_, e := o.QueryTable(new(models.Vote)).Limit(-1).All(&v)
	if e != nil {
		beego.Error("query vote error", e)
		msg = "系统异常(1)"
		return
	}
	var votes []int64
	for i := 1; i < 13; i++ {
		count, _ := o.QueryTable(new(models.Vote)).Filter("Category", i).Count()
		votes = append(votes, count)
	}
	code = 1
	data["votes"] = votes
}

func (this *IndexVoteApiIndexController) AllowCross() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                           //允许访问源
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	this.Ctx.ResponseWriter.Header().Set("content-type", "application/json")
}
