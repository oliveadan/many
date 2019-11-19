package customer

import (
	"phagego/phagev2/controllers/sysmanage"
	"github.com/astaxie/beego/orm"
	."phagego/phagev2/models"
	"phagego/phage-check-web/models"
	"github.com/astaxie/beego"

	"strings"
)

type CustomerIndexController struct {
	 sysmanage.BaseController
}

func (this *CustomerIndexController) Get(){
	username := this.GetString("username")
	var adminList []Admin
	o := orm.NewOrm()
	//获取所有客服组的管理员
	ids := make([]int64,0)
	adminroles := GetAdminRoles()
	for _,v := range adminroles {
		ids = append(ids,v.AdminId)
	}

	condArr := make(map[string]string)
	qs := o.QueryTable(new(Admin))
	cond := orm.NewCondition()
	if len(ids)==0 {
		ids = append(ids,0)
	}
	cond = cond.And("Id__in",ids)
	cond = cond.And("IsSystem", 0)
	if username != "" {
		cond = cond.And("username__contains", username)
		condArr["username"] = username
	}
	qs = qs.SetCond(cond)
	qs.All(&adminList)

	// 返回值
	this.Data["dataList"] = adminList
	// 查询条件
	this.Data["condArr"] = condArr
	this.TplName = "sysmanage/customer/index.html"
}


type CustomerEditController struct {
	sysmanage.BaseController
}

func(this *CustomerEditController)  Get(){
	beego.Informational("修改客户权限")
    id,_ := this.GetInt64("id")
    //获取字段权限表
	checkpermission :=  models.GetCheckPermission(id)
	//获取所有分层
	checks := models.GetHierarchys()
	//获取管理员账号所在分层
    hierarchy := models.GetHierarchy(id)
    this.Data["adminid"] = id
    this.Data["hierarchy"] = hierarchy
	this.Data["checks"] = checks
	this.Data["checkpermission"] = checkpermission
	this.TplName = "sysmanage/customer/edit.html"
}

func(this *CustomerEditController) Post() {
     var code int
     var msg  string
     var url  = this.URLFor("CustomerIndexController.Get")
     defer sysmanage.Retjson(this.Ctx,&msg,&code,&url)
     adminid,_ := this.GetInt64("adminid")
	 hierarchy := strings.TrimSpace(this.GetString("hierarchy"))
	 //更新分层信息
	 o := orm.NewOrm()
	 _ , err := o.QueryTable(new(models.CheckRole)).Filter("AdminId",adminid).Update(orm.Params{
	 	"hierarchy": hierarchy})
	if err != nil {
		beego.Informational("更新分层信息失败",err)
		msg = "更新分层信息失败"
		return
	}
	//更新可查看字段

	Account,_ := this.GetInt64("Account")
	idd,_ := this.GetInt64("Idd")
	Name,_ := this.GetInt64("Name")
	Agent,_ := this.GetInt64("Agent")
	LoginInformation,_ := this.GetInt64("LoginInformation")
	RegisterDate,_ := this.GetInt64("RegisterDate")
	Mobile,_ := this.GetInt64("Mobile")
	Email,_ := this.GetInt64("Email")
	Qq,_ := this.GetInt64("Qq")
	Wechat,_ := this.GetInt64("Wechat")
	PasswordHint,_ := this.GetInt64("PasswordHint")
	PasswordAnswer,_ := this.GetInt64("PasswordAnswer")
	WithdrawalPassword,_ := this.GetInt64("WithdrawalPassword")
	OpenBank,_ := this.GetInt64("OpenBank")
	BankAccount,_ := this.GetInt64("BankAccount")
	_ ,err1 := o.QueryTable(new(models.CheckPermission)).Filter("AdminId",adminid).Update(orm.Params{
		"Account":Account,"Idd":idd,"Name":Name,"Agent":Agent,"LoginInformation":LoginInformation,
		"RegisterDate":RegisterDate,"Mobile":Mobile,"Email":Email,"Qq":Qq,"Wechat":Wechat,
		"PasswordHint":PasswordHint,"PasswordAnswer":PasswordAnswer,"WithdrawalPassword":WithdrawalPassword,
		"OpenBank":OpenBank,"BankAccount":BankAccount})
	if err1 != nil {
		beego.Informational("更新字段权限表失败",err)
		msg = "更新字段权限表失败"
		return
	}
	code = 1
	msg = "更新成功"
}

