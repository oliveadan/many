package login

import (
	"fmt"
	"html/template"
	. "phagego/phagev2/models"
	. "phagego/common/utils"
	"time"
	. "phagego/phagev2/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	beego.Warn("LoginController Get from ip:", this.Ctx.Input.IP())
	var pubkey, prikey string
	if this.GetSession("loginpubkey") != nil {
		pubkey = this.GetSession("loginpubkey").(string)
	} else {
		pubkey, prikey = RsaGenerateKey(1024)
		this.SetSession("loginprikey", prikey)
		this.SetSession("loginpubkey", pubkey)
	}
	if beego.BConfig.RunMode == "dev" {
		this.Data["username"] = "admin"
		this.Data["pass"] = "111111"
		this.Data["captchaValue"] = "1"
	} else {
		this.Data["username"] = ""
		this.Data["pass"] = ""
		this.Data["captchaValue"] = ""
	}
	this.Data["year"] = time.Now().Year()
	this.Data["pubkey"] = pubkey
	this.Data["siteName"] = GetSiteConfigValue(Scname)
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "sysmanage/login/index.html"
}

func (this *LoginController) Post() {
	ret := make(map[string]interface{})
	username := this.GetString("username")
	pwd := this.GetString("password")
	if username == "" {
		ret["msg"] = "用户名不能为空"
	} else if pwd == "" {
		ret["msg"] = "密码不能为空"
	} else if beego.BConfig.RunMode == "prod" && !GetCpt().VerifyReq(this.Ctx.Request) {
		ret["msg"] = "验证码错误"
	} else {
		if this.GetSession("loginprikey") == nil {
			ret["msg"] = "请刷新后再试"
		} else {
			prikey := this.GetSession("loginprikey").(string)
			pwdDecrypt := RsaDecrypt(pwd, prikey)
			//pwdDecrypt := pwd
			o := orm.NewOrm()
			admin := Admin{Username: username}
			err := o.Read(&admin, "Username")
			if err != nil {
				beego.Error("Login error", err)
				ret["msg"] = "用户名或密码错误"
			} else {
				cols := make([]string, 0)
				if admin.Enabled == 0 {
					ret["msg"] = "用户名或密码错误"
				} else if admin.Locked == 1 {
					ret["msg"] = "账号已被锁定，无法登录"
				} else if admin.Password != Md5(pwdDecrypt, Pubsalt, admin.Salt) {
					ret["msg"] = "用户名或密码错误"
					admin.LoginFailureCount += 1
					if admin.LoginFailureCount >= 5 {
						admin.Locked = 1
						cols = append(cols, "Locked")
					}
					cols = append(cols, "LoginFailureCount")
				} else {
					if admin.LoginVerify == 1 {
						go SendMailVerifyCode(admin.Email)
						ret["code"] = 2  // 需要二次验证
						ret["msg"] = "登录成功，请输入邮箱验证码"
					} else {
						token := GetGuid()
						SetCache(fmt.Sprintf("loginAdminId%d", admin.Id), token, 28800)
						this.SetSession("token", token)
						this.SetSession("loginAdminId", admin.Id)
						this.SetSession("loginAdminName", admin.Name)
						this.SetSession("loginAdminUsername", admin.Username)
						ret["code"] = 1
						ret["msg"] = "登录成功"
						ret["url"] = this.URLFor("BaseController.Index")
					}
					admin.LoginFailureCount = 0
					admin.LoginIp = this.Ctx.Input.IP()
					admin.LoginDate = time.Now()
					cols = append(cols, "LoginFailureCount", "LoginIp", "LoginDate")
				}
				if len(cols) > 0 {
					o.Update(&admin, cols...)
				}
			}
		}
	}
	beego.Warn("LoginController Post from ip:", this.Ctx.Input.IP(), "username:", username)
	this.Data["json"] = ret
	this.ServeJSON()
}

func (this *LoginController) LoginMailVerify() {
	ret := make(map[string]interface{})
	username := this.GetString("username")
	verifyCode := this.GetString("code")
	if verifyCode == "" {
		ret["msg"] = "邮箱验证码不能为空"
	} else {
		o := orm.NewOrm()
		admin := Admin{Username: username}
		if err := o.Read(&admin, "Username"); err != nil {
			ret["msg"] = "验证失败，请重试"
		} else {
			if isVerify := VerifyMailVerifyCode(admin.Email, verifyCode); isVerify {
				token := GetGuid()
				SetCache(fmt.Sprintf("loginAdminId%d", admin.Id), token, 28800)
				this.SetSession("token", token)
				this.SetSession("loginAdminId", admin.Id)
				this.SetSession("loginAdminName", admin.Name)
				this.SetSession("loginAdminUsername", admin.Username)
				ret["code"] = 1
				ret["msg"] = "验证成功"
				ret["url"] = this.URLFor("BaseController.Index")
			} else {
				ret["msg"] = "验证失败"
			}
		}
	}
	this.Data["json"] = ret
	this.ServeJSON()
}

func (this *LoginController) Logout() {
	DelCache(fmt.Sprintf("loginAdminId%v", this.GetSession("loginAdminId")))
	this.DelSession("loginAdminId")
	this.Redirect(this.URLFor("LoginController.Get"), 302)
}
