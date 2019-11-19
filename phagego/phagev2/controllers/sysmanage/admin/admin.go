package admin

import (
	"fmt"
	"html/template"
	"phagego/phagev2/controllers/sysmanage"
	. "phagego/phagev2/models"
	. "phagego/phagev2/utils"
	. "phagego/common/utils"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
	"phagego/phage-check-web/models"
)

func validate(admin *Admin) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(admin.Username, "errmsg").Message("用户名必输")
	valid.AlphaNumeric(admin.Username, "errmsg").Message("用户名必须为字母和数字")
	valid.MaxSize(admin.Username, 30, "errmsg").Message("用户名最长30位")
	valid.Required(admin.Name, "errmsg").Message("名称必输")
	valid.MaxSize(admin.Name, 30, "errmsg").Message("名称最长30位")
	valid.MaxSize(admin.Password, 30, "errmsg").Message("密码最长30位")
	if admin.Email != "" {
		valid.Email(admin.Email, "errmsg").Message("邮箱格式不正确")
	}
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type AdminIndexController struct {
	sysmanage.BaseController
}

func (this *AdminIndexController) Get() {
	username := this.GetString("username")
	var adminList []Admin
	condArr := make(map[string]string)
	o := orm.NewOrm()
	qs := o.QueryTable(new(Admin))
	cond := orm.NewCondition()
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
	this.TplName = "sysmanage/admin/index.html"
}

func (this *AdminIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, err := this.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		beego.Error("Delete Admin error", err)
		return
	}
	o := orm.NewOrm()
	o.Begin()
	// 删除管理员角色关联
	if _, err := o.QueryTable(new(AdminRole)).Filter("AdminId", id).Delete(); err != nil {
		o.Rollback()
		beego.Error("Delete admin error 1", err)
		msg = "删除失败"
		return
	}
	o.Begin()
	//删除客服角色关系
	if _, err := o.QueryTable(new(models.CheckPermission)).Filter("AdminId",id).Delete() ; err != nil {
		o.Rollback()
		beego.Error("Delete checkpermission err 2",err)
		msg = "删除失败"
		return
	}
	o.Begin()
	if _,err := o.QueryTable(new(models.CheckRole)).Filter("AdminId",id).Delete() ;err != nil {
		o.Rollback()
		beego.Error("Delete checkrole err 3",err)
		msg = "删除失败"
		return
	}
	if _, err := o.Delete(&Admin{Id: id}); err != nil {
		o.Rollback()
		beego.Error("Delete admin error 2", err)
		msg = "删除失败"
	} else {
		o.Commit()
		code = 1
		msg = "删除成功"
	}
}

func (this *AdminIndexController) LoginVerify() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, err := this.GetInt64("id")
	t := this.GetString("type")
	verifyCode := this.GetString("code")

	if err != nil {
		msg = "数据错误"
		beego.Error("LoginVerify Admin error", err)
		return
	}
	o := orm.NewOrm()
	model := Admin{}
	model.Id = id
	if err := o.Read(&model); err != nil {
		beego.Error("Read admin error", err)
		msg = "操作失败，请刷新后重试"
		return
	}
	if model.LoginVerify == 1 {
		model.LoginVerify = 0
	} else {
		if model.Email == "" {
			msg = "邮箱未配置，请先配置邮箱"
			return
		}
		if t == "send" {
			if err := SendMailVerifyCode(model.Email); err != nil {
				msg = "验证码发送失败"
				return
			} else {
				code = 1
				msg = "验证码发送成功，请查看收件箱或垃圾箱"
				return
			}
		} else {
			if isVerify := VerifyMailVerifyCode(model.Email, verifyCode); !isVerify {
				msg = "验证失败，请重试"
				return
			}
		}
		// 启用登录验证
		model.LoginVerify = 1
	}

	if _, err := o.Update(&model, "LoginVerify"); err != nil {
		beego.Error("Update admin error", err)
		msg = "操作失败，请刷新后重试"
	} else {
		code = 1
		msg = "操作成功"
	}
}

func (this *AdminIndexController) Locked() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, err := this.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		beego.Error("Locked Admin error", err)
		return
	}
	o := orm.NewOrm()
	model := Admin{}
	model.Id = id
	if err := o.Read(&model); err != nil {
		beego.Error("Read admin error", err)
		msg = "操作失败，请刷新后重试"
		return
	}
	if model.Locked == 1 {
		model.Locked = 0
		model.LoginFailureCount = 0
	} else {
		model.Locked = 1
		model.LockedDate = time.Now()
	}

	if _, err := o.Update(&model, "Locked", "LockedDate"); err != nil {
		beego.Error("Update admin error", err)
		msg = "操作失败，请刷新后重试"
	} else {
		code = 1
		msg = "操作成功"
		if model.Locked == 1 { // 如果是锁定，则一并清楚登录token，强制用户退出
			DelCache(fmt.Sprintf("loginAdminId%d", id))
		}
	}
}

type AdminAddController struct {
	sysmanage.BaseController
}

func (this *AdminAddController) Get() {
	this.Data["roleList"] = GetRoleList()
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "sysmanage/admin/add.html"
}

func (this *AdminAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	admin := Admin{}
	if err := this.ParseForm(&admin); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&admin); hasError {
		msg = errMsg
		return
	} else if admin.Password == "" {
		msg = "密码不能为空"
		return
	} else if admin.Password != this.GetString("repassword") {
		msg = "两次输入的密码不一致"
		return
	}
	roles := this.GetStrings("roles")
	if len(roles) == 0 {
		msg = "请选择所属权限组"
		return
	}
	salt := GetGuid()
	pa := Md5(Md5(admin.Password), Pubsalt, salt)
	admin.Password = pa
	admin.Salt = salt
	admin.Creator = this.LoginAdminId
	admin.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	o.Begin()
	if created, _, err := o.ReadOrCreate(&admin, "Username"); err != nil {
		o.Rollback()
		msg = "添加失败"
		beego.Error("Insert admin error 1", err)
	} else if created {
		adminRoles := make([]AdminRole, 0)
		for _, v := range roles {
			roleId, _ := strconv.ParseInt(v, 10, 64)
			ar := AdminRole{AdminId: admin.Id, RoleId: roleId}
			adminRoles = append(adminRoles, ar)
		}
		if _, err := o.InsertMulti(len(adminRoles), adminRoles); err != nil {
			o.Rollback()
			msg = "添加失败"
			beego.Error("Insert admin error 3", err)
			return
		}
		o.Commit()

		for _, v := range roles {
			if v == "3" {
				o.Begin()
				checkrole := models.CheckRole{AdminId: admin.Id}
				if _, err := checkrole.Create(); err != nil {
					o.Rollback()
					msg = "添加失败"
					beego.Error("Insert checkrole error 2", err)
					return
				}
				o.Commit()
				o.Begin()
				checkpermission := models.CheckPermission{AdminId: admin.Id}
				if _, err := checkpermission.Create(); err != nil {
					o.Rollback()
					msg = "添加失败"
					beego.Error("Insert checkpermission error4", err)
					return
				}
				o.Commit()
			}
		}
		code = 1
		msg = "添加成功"
	} else {
		msg = "账号已存在"
	}
}

type AdminEditController struct {
	sysmanage.BaseController
}

func (this *AdminEditController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	admin := Admin{Id: id}

	err := o.Read(&admin)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("AdminIndexController.get"), 302)
	} else {
		// 当前管理员所属角色
		var arList orm.ParamsList
		o.QueryTable(new(AdminRole)).Filter("AdminId", id).ValuesFlat(&arList, "RoleId")
		arMap := make(map[int64]bool)
		for _, v := range arList {
			arId, ok := v.(int64)
			if ok {
				arMap[arId] = true
			}
		}

		this.Data["data"] = admin
		this.Data["adminRoleMap"] = arMap
		this.Data["roleList"] = GetRoleList()
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "sysmanage/admin/edit.html"
	}
}

func (this *AdminEditController) Post() {
	var code int
	var msg string
	var reurl = this.URLFor("AdminIndexController.Get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &reurl)
	admin := Admin{}
	if err := this.ParseForm(&admin); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&admin); hasError {
		msg = errMsg
		return
	} else if admin.Password != "" && admin.Password != this.GetString("repassword") {
		msg = "两次输入的密码不一致"
		return
	}
	cols := []string{"Username", "Name", "Enabled", "Email", "ModifyDate"}
	isChangePwd := false
	if admin.Password != "" {
		salt := GetGuid()
		pa := Md5(Md5(admin.Password), Pubsalt, salt)
		admin.Password = pa
		admin.Salt = salt
		cols = append(cols, "Password", "Salt")
		isChangePwd = true
	}
	admin.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	o.Begin()
	if _, err := o.Update(&admin, cols...); err != nil {
		o.Rollback()
		msg = "更新失败"
		beego.Error("Update admin error 1", err)
	} else {
		// 删除旧角色
		if _, err := o.QueryTable(new(AdminRole)).Filter("AdminId", admin.Id).Delete(); err != nil {
			o.Rollback()
			msg = "更新失败"
			beego.Error("Update admin error 2", err)
		}
		// 重新插入角色
		roles := this.GetStrings("roles")
		adminRoles := make([]AdminRole, 0)
		for _, v := range roles {
			roleId, _ := strconv.ParseInt(v, 10, 64)
			ar := AdminRole{AdminId: admin.Id, RoleId: roleId}
			adminRoles = append(adminRoles, ar)
		}

		if _, err := o.InsertMulti(len(adminRoles), adminRoles); err != nil {
			o.Rollback()
			msg = "更新失败"
			beego.Error("Update admin error 3", err)
		}
		o.Commit()
		// 如修改了密码，则重置登录，让用户必须重新登录
		if isChangePwd {
			DelCache(fmt.Sprintf("loginAdminId%d", admin.Id))
		}

		code = 1
		msg = "更新成功"
	}
}
