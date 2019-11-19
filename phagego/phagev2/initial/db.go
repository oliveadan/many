package initial

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	. "phagego/phagev2/models"
)

func InitSql() {
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	maxIdle := 30
	maxConn := 30
	var dataSource string
	dbDriver := beego.AppConfig.String("dbdriver")
	if dbDriver == "mysql" {
		user := beego.AppConfig.String("mysqluser")
		passwd := beego.AppConfig.String("mysqlpass")
		host := beego.AppConfig.String("mysqlurls")
		port, err := beego.AppConfig.Int("mysqlport")
		dbname := beego.AppConfig.String("mysqldb")

		if nil != err {
			port = 3306
		}
		dataSource = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=%s", user, passwd, host, port, dbname, "Asia%2FShanghai")
	} else if dbDriver == "sqlite3" {
		dataSource = "data.db"
	} else {
		panic("未知数据库驱动类型")
	}
	orm.RegisterDataBase("default", dbDriver, dataSource, maxIdle, maxConn)
}

func InitDbFrameData() {
	// 自动建表
	autoCreateDb := beego.AppConfig.DefaultInt("dbautocreate", 0)
	if autoCreateDb == 1 || autoCreateDb == 2 {
		beego.Info("Auto create db")
		isForce := false
		if autoCreateDb == 2 { // drop table 后再建表
			isForce = true
		}
		// 遇到错误立即返回
		err := orm.RunSyncdb("default", isForce, true)
		if err != nil {
			beego.Error("Auto create db error", err.Error())
			return
		}
		if err := initDbData(); err != nil {
			return
		}
	}
}

// 初始化数据库数据--框架数据
func initDbData() error {
	o := orm.NewOrm()
	// 通过 admin表判断数据库是否已初始化，已存在数据，则说明已初始化过。
	if isExist := o.QueryTable(new(Admin)).Exist(); isExist {
		return nil
	}
	beego.Info("Init frame data")
	o.Begin()
	// 系统配置
	sc := SiteConfig{Id: 1, Creator: 0, Modifior: 0, Version: 0, Code: "NAME", Value: "平台名称", IsSystem: 1}
	if num, err := o.Insert(&sc); err != nil || num != 1 {
		o.Rollback()
		beego.Error("Init SiteConfig data error", err)
		return err
	}
	// 管理员
	admins := []Admin{
		{Id: 1, Creator: 0, Modifior: 0, Version: 0, Enabled: 1, Locked: 0, IsSystem: 1, LoginFailureCount: 0, Salt: "17b007bdb8e7af362a1167bcce7277c9", Name: "超级管理员", Password: "9b16a6a8b524be91d0f440f61ed76fab", Username: "superadmin", LoginVerify: 0},
		{Id: 2, Creator: 0, Modifior: 0, Version: 0, Enabled: 1, Locked: 0, IsSystem: 0, LoginFailureCount: 0, Salt: "253da8a9583fccd5645690aa25a71d20", Name: "管理员", Password: "ec7d8fc2e0093ffec5f39fede8e0bdd6", Username: "admin", LoginVerify: 0},
	}
	if num, err := o.InsertMulti(len(admins), admins); err != nil || int(num) != len(admins) {
		o.Rollback()
		beego.Error("Init Admin data error", err)
		return err
	}
	// 角色
	roles := []Role{
		{Id: 1, Creator: 0, Modifior: 0, Version: 0, Enabled: 1, Description: "后台管理最高权限", IsSystem: 1, Name: "超级管理员"},
		{Id: 2, Creator: 0, Modifior: 0, Version: 0, Enabled: 1, Description: "后台管理权限", IsSystem: 0, Name: "管理员"},
		{Id: 3, Creator: 0, Modifior: 0, Version: 0, Enabled: 1, Description: "可以查看规定的字段", IsSystem: 0, Name: "客服"},
	}
	if num, err := o.InsertMulti(len(roles), roles); err != nil || int(num) != len(roles) {
		o.Rollback()
		beego.Error("Init Role data error", err)
		return err
	}
	// 管理员--角色关联
	adminRoles := []AdminRole{
		{Id: 1, AdminId: 1, RoleId: 1},
		{Id: 2, AdminId: 2, RoleId: 2},
	}
	if num, err := o.InsertMulti(len(adminRoles), adminRoles); err != nil || int(num) != len(adminRoles) {
		o.Rollback()
		beego.Error("Init AdminRole data error", err)
		return err
	}
	// 菜单权限配置
	permissions := []Permission{
		{Id: 1, Creator: 0, Modifior: 0, Version: 0, Pid: 0, Enabled: 1, Display: 0, Description: "系统框架", Url: "BaseController.Index", Name: "系统框架", Icon: "", Sort: 1},
		{Id: 2, Creator: 0, Modifior: 0, Version: 0, Pid: 0, Enabled: 1, Display: 0, Description: "修改密码", Url: "ChangePwdController.Get", Name: "修改密码", Icon: "", Sort: 2},
		{Id: 3, Creator: 0, Modifior: 0, Version: 0, Pid: 0, Enabled: 1, Display: 0, Description: "系统信息", Url: "SysIndexController.Get", Name: "系统信息", Icon: "", Sort: 3},
		{Id: 10, Creator: 0, Modifior: 0, Version: 0, Pid: 0, Enabled: 1, Display: 0, Description: "系统通用-文件上传", Url: "SyscommonController.Upload", Name: "系统通用-文件上传", Icon: "", Sort: 10},
		{Id: 20, Creator: 0, Modifior: 0, Version: 0, Pid: 0, Enabled: 1, Display: 1, Description: "系统设置", Url: "", Name: "系统设置", Icon: "#xe716;", Sort: 100},
		{Id: 21, Creator: 0, Modifior: 0, Version: 0, Pid: 20, Enabled: 1, Display: 1, Description: "管理员", Url: "AdminIndexController.Get", Name: "管理员", Icon: "", Sort: 100},
		{Id: 22, Creator: 0, Modifior: 0, Version: 0, Pid: 21, Enabled: 1, Display: 0, Description: "添加管理员", Url: "AdminAddController.Get", Name: "添加管理员", Icon: "", Sort: 100},
		{Id: 23, Creator: 0, Modifior: 0, Version: 0, Pid: 21, Enabled: 1, Display: 0, Description: "编辑管理员", Url: "AdminEditController.Get", Name: "编辑管理员", Icon: "", Sort: 100},
		{Id: 24, Creator: 0, Modifior: 0, Version: 0, Pid: 21, Enabled: 1, Display: 0, Description: "删除管理员", Url: "AdminIndexController.Delone", Name: "删除管理员", Icon: "", Sort: 100},
		{Id: 25, Creator: 0, Modifior: 0, Version: 0, Pid: 21, Enabled: 1, Display: 0, Description: "锁定解锁管理员", Url: "AdminIndexController.Locked", Name: "锁定解锁管理员", Icon: "", Sort: 100},
		{Id: 26, Creator: 0, Modifior: 0, Version: 0, Pid: 21, Enabled: 1, Display: 0, Description: "管理员登录验证", Url: "AdminIndexController.LoginVerify", Name: "管理员登录验证", Icon: "", Sort: 100},
		{Id: 30, Creator: 0, Modifior: 0, Version: 0, Pid: 20, Enabled: 1, Display: 1, Description: "角色管理", Url: "RoleIndexController.Get", Name: "角色管理", Icon: "", Sort: 100},
		{Id: 31, Creator: 0, Modifior: 0, Version: 0, Pid: 30, Enabled: 1, Display: 0, Description: "添加角色", Url: "RoleAddController.Get", Name: "添加角色", Icon: "", Sort: 100},
		{Id: 32, Creator: 0, Modifior: 0, Version: 0, Pid: 30, Enabled: 1, Display: 0, Description: "编辑角色", Url: "RoleEditController.Get", Name: "编辑角色", Icon: "", Sort: 100},
		{Id: 33, Creator: 0, Modifior: 0, Version: 0, Pid: 30, Enabled: 1, Display: 0, Description: "删除角色", Url: "RoleIndexController.Delone", Name: "删除角色", Icon: "", Sort: 100},
		{Id: 40, Creator: 0, Modifior: 0, Version: 0, Pid: 20, Enabled: 1, Display: 1, Description: "菜单管理", Url: "PermissionIndexController.Get", Name: "菜单管理", Icon: "", Sort: 100},
		{Id: 41, Creator: 0, Modifior: 0, Version: 0, Pid: 40, Enabled: 1, Display: 0, Description: "添加菜单", Url: "PermissionAddController.Get", Name: "添加菜单", Icon: "", Sort: 100},
		{Id: 42, Creator: 0, Modifior: 0, Version: 0, Pid: 40, Enabled: 1, Display: 0, Description: "编辑菜单", Url: "PermissionEditController.Get", Name: "编辑菜单", Icon: "", Sort: 100},
		{Id: 43, Creator: 0, Modifior: 0, Version: 0, Pid: 40, Enabled: 1, Display: 0, Description: "删除菜单", Url: "PermissionIndexController.Delone", Name: "删除菜单", Icon: "", Sort: 100},
		{Id: 50, Creator: 0, Modifior: 0, Version: 0, Pid: 20, Enabled: 1, Display: 1, Description: "站点配置", Url: "SiteConfigIndexController.Get", Name: "站点配置", Icon: "", Sort: 100},
		{Id: 51, Creator: 0, Modifior: 0, Version: 0, Pid: 50, Enabled: 1, Display: 0, Description: "添加站点配置", Url: "SiteConfigAddController.Get", Name: "添加站点配置", Icon: "", Sort: 100},
		{Id: 52, Creator: 0, Modifior: 0, Version: 0, Pid: 50, Enabled: 1, Display: 0, Description: "编辑站点配置", Url: "SiteConfigEditController.Get", Name: "编辑站点配置", Icon: "", Sort: 100},
		{Id: 53, Creator: 0, Modifior: 0, Version: 0, Pid: 50, Enabled: 1, Display: 0, Description: "删除站点配置", Url: "SiteConfigIndexController.Delone", Name: "删除站点配置", Icon: "", Sort: 100},
		{Id: 60, Creator: 0, Modifior: 0, Version: 0, Pid: 20, Enabled: 1, Display: 1, Description: "快捷导航", Url: "QuickNavIndexController.Get", Name: "快捷导航", Icon: "", Sort: 100},
		{Id: 61, Creator: 0, Modifior: 0, Version: 0, Pid: 60, Enabled: 1, Display: 0, Description: "添加快捷导航", Url: "QuickNavAddController.Get", Name: "添加快捷导航", Icon: "", Sort: 100},
		{Id: 62, Creator: 0, Modifior: 0, Version: 0, Pid: 60, Enabled: 1, Display: 0, Description: "编辑快捷导航", Url: "QuickNavEditController.Get", Name: "编辑快捷导航", Icon: "", Sort: 100},
		{Id: 63, Creator: 0, Modifior: 0, Version: 0, Pid: 60, Enabled: 1, Display: 0, Description: "删除快捷导航", Url: "QuickNavIndexController.Delone", Name: "删除快捷导航", Icon: "", Sort: 100},
		{Id: 64, Creator: 0, Modifior: 0, Version: 0, Pid: 20, Enabled: 1, Display: 1, Description: "客服权限", Url: "CustomerIndexController.Get", Name: "客服权限", Icon: "", Sort: 100},
		{Id: 65, Creator: 0, Modifior: 0, Version: 0, Pid: 64, Enabled: 1, Display: 0, Description: "修改客服权限", Url: "CustomerEditController.Get", Name: "修改客服权限", Icon: "", Sort: 100},
	}
	if num, err := o.InsertMulti(len(permissions), permissions); err != nil || int(num) != len(permissions) {
		o.Rollback()
		beego.Error("Init Permission data error", err)
		return err
	}
	// 角色--权限关联
	rolePermissions := []RolePermission{
		{Id: 1, RoleId: 1, PermissionId: 1},
		{Id: 2, RoleId: 1, PermissionId: 2},
		{Id: 3, RoleId: 1, PermissionId: 3},
		{Id: 20, RoleId: 1, PermissionId: 20},
		{Id: 21, RoleId: 1, PermissionId: 21},
		{Id: 22, RoleId: 1, PermissionId: 22},
		{Id: 23, RoleId: 1, PermissionId: 23},
		{Id: 24, RoleId: 1, PermissionId: 24},
		{Id: 25, RoleId: 1, PermissionId: 25},
		{Id: 26, RoleId: 1, PermissionId: 26},
		{Id: 30, RoleId: 1, PermissionId: 30},
		{Id: 31, RoleId: 1, PermissionId: 31},
		{Id: 32, RoleId: 1, PermissionId: 32},
		{Id: 33, RoleId: 1, PermissionId: 33},
		{Id: 40, RoleId: 1, PermissionId: 40},
		{Id: 41, RoleId: 1, PermissionId: 41},
		{Id: 42, RoleId: 1, PermissionId: 42},
		{Id: 43, RoleId: 1, PermissionId: 43},
		{Id: 50, RoleId: 1, PermissionId: 50},
		{Id: 51, RoleId: 1, PermissionId: 51},
		{Id: 52, RoleId: 1, PermissionId: 52},
		{Id: 53, RoleId: 1, PermissionId: 53},
		{Id: 60, RoleId: 1, PermissionId: 60},
		{Id: 61, RoleId: 1, PermissionId: 61},
		{Id: 62, RoleId: 1, PermissionId: 62},
		{Id: 63, RoleId: 1, PermissionId: 63},
		{Id: 64, RoleId: 1, PermissionId: 64},
		{Id: 65, RoleId: 1, PermissionId: 65},
		{Id: 101, RoleId: 2, PermissionId: 1},
		{Id: 102, RoleId: 2, PermissionId: 2},
		{Id: 103, RoleId: 2, PermissionId: 3},
		{Id: 110, RoleId: 2, PermissionId: 10},
		{Id: 120, RoleId: 2, PermissionId: 20},
		{Id: 121, RoleId: 2, PermissionId: 21},
		{Id: 122, RoleId: 2, PermissionId: 22},
		{Id: 123, RoleId: 2, PermissionId: 23},
		{Id: 124, RoleId: 2, PermissionId: 24},
		{Id: 125, RoleId: 2, PermissionId: 25},
		{Id: 126, RoleId: 2, PermissionId: 26},
		{Id: 164, RoleId: 2, PermissionId: 64},
		{Id: 166, RoleId: 2, PermissionId: 65},
		{Id: 167, RoleId: 3, PermissionId: 1},
		{Id: 168, RoleId: 3, PermissionId: 2},
		{Id: 169, RoleId: 3, PermissionId: 3},

	}
	if num, err := o.InsertMulti(len(rolePermissions), rolePermissions); err != nil || int(num) != len(rolePermissions) {
		o.Rollback()
		beego.Error("Init RolePermission data error", err)
		return err
	}
	o.Commit()
	return nil
}
