package initial

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	frame "phagego/frameweb-v2/initial"
	. "phagego/frameweb-v2/models"
)

// 本项目数据初始化
func InitDbProjectData() {
	// 初始化系统基础数据
	frame.InitDbFrameData()
	// 初始化项目数据
	if beego.AppConfig.DefaultInt("dbautocreate", 0) == 0 {
		return
	}
	beego.Info("Init project data")
	permisions := []Permission{
		//vip设置
		{Id: 100, Pid: 0, Enabled: 1, Display: 1, Description: "vip设置", Url: "", Name: "vip设置", Icon: "#xe614;", Sort: 100},
		{Id: 101, Pid: 100, Enabled: 1, Display: 1, Description: "vip等级", Url: "LevelController.Get", Name: "vip等级", Icon: "", Sort: 100},
		{Id: 102, Pid: 100, Enabled: 1, Display: 0, Description: "添加vip等级", Url: "LevelAddController.Get", Name: "添加vip等级", Icon: "", Sort: 100},
		{Id: 103, Pid: 100, Enabled: 1, Display: 0, Description: "修改vip等级", Url: "LevelEditController.Get", Name: "修改vip等级", Icon: "", Sort: 100},
		{Id: 104, Pid: 100, Enabled: 1, Display: 0, Description: "删除vip等级", Url: "LevelController.Delone", Name: "删除vip等级", Icon: "", Sort: 100},
		//周期分类
		{Id: 110, Pid: 100, Enabled: 1, Display: 1, Description: "周期分类", Url: "PeriodIndexController.get", Name: "周期分类", Icon: "", Sort: 100},
		{Id: 111, Pid: 100, Enabled: 1, Display: 0, Description: "添加周期分类名称", Url: "PeriodAddController.get", Name: "添加周期分类名称", Icon: "", Sort: 100},
		{Id: 112, Pid: 100, Enabled: 1, Display: 0, Description: "修改周期分类", Url: "PeriodEditController.get", Name: "修改周期分类", Icon: "", Sort: 100},
		{Id: 113, Pid: 100, Enabled: 1, Display: 0, Description: "删除周期分类", Url: "PeriodIndexController.Delone", Name: "删除周期分类", Icon: "", Sort: 100},
        //会员投注
		{Id: 120, Pid: 100, Enabled: 1, Display: 1, Description: "会员投注", Url: "MembersingleIndexController.get", Name: "会员投注", Icon: "", Sort: 100},
		{Id: 121, Pid: 100, Enabled: 1, Display: 0, Description: "添加会员投注", Url: "MembersingleAddController.get", Name: "添加会员投注", Icon: "", Sort: 100},
		{Id: 122, Pid: 100, Enabled: 1, Display: 0, Description: "修改会员投注", Url: "MembersingleEditController.get", Name: "修改会员投注", Icon: "", Sort: 100},
		{Id: 123, Pid: 100, Enabled: 1, Display: 0, Description: "删除会员投注", Url: "MembersingleIndexController.Delone", Name: "删除会员投注", Icon: "", Sort: 100},
		{Id: 124, Pid: 100, Enabled: 1, Display: 0, Description: "导入会员投注", Url: "MembersingleIndexController.Import", Name: "导入会员投注", Icon: "", Sort: 100},
		{Id: 125, Pid: 100, Enabled: 1, Display: 0, Description: "删除一期会员投注", Url: "MembersingleIndexController.DelBatch", Name: "删除一期会员投注", Icon: "", Sort: 100},
		{Id: 126, Pid: 100, Enabled: 1, Display: 0, Description: "计算本期彩金", Url: "MembersingleIndexController.CountGift", Name: "计算本期彩金", Icon: "", Sort: 100},
		{Id: 127, Pid: 100, Enabled: 1, Display: 0, Description: "导出会员投注", Url: "MembersingleIndexController.Export", Name: "导出会员投注", Icon: "", Sort: 100},
		{Id: 128, Pid: 100, Enabled: 1, Display: 0, Description: "计算本期好运彩金", Url: "MembersingleIndexController.CountLuckyGift", Name: "计算本期好运彩金", Icon: "", Sort: 100},
        //会员统计
		{Id: 130, Pid: 100, Enabled: 1, Display: 0, Description: "会员统计", Url: "MemberTotalIndexController.get", Name: "会员统计", Icon: "", Sort: 100},
		{Id: 131, Pid: 100, Enabled: 1, Display: 0, Description: "删除所有会员统计", Url: "MemberTotalIndexController.Delbatch", Name: "删除所有会员统计", Icon: "", Sort: 100},
		{Id: 132, Pid: 100, Enabled: 1, Display: 0, Description: "导出会员统计", Url: "MemberTotalIndexController.Export", Name: "导出会员统计", Icon: "", Sort: 100},
		{Id: 133, Pid: 100, Enabled: 1, Display: 0, Description: "修改会员统计", Url: "MemberTotalEditController.get", Name: "修改会员统计", Icon: "", Sort: 100},
		//好运金配置
		{Id: 140, Pid: 100, Enabled: 1, Display: 1, Description: "好运金", Url: "LuckyController.Get", Name: "好运金", Icon: "", Sort: 99},
		{Id: 141, Pid: 100, Enabled: 1, Display: 0, Description: "添加好运金配置", Url: "LuckyAddController.Get", Name: "添加好运金配置", Icon: "", Sort: 99},
		{Id: 142, Pid: 100, Enabled: 1, Display: 0, Description: "修改好运金配置", Url: "LuckyEditController.Get", Name: "修改好运金配置", Icon: "", Sort: 99},
		{Id: 143, Pid: 100, Enabled: 1, Display: 0, Description: "删除好运金配置", Url: "LuckyController.Delone", Name: "删除好运金配置", Icon: "", Sort: 99},
        //中奖记录
		{Id: 149, Pid: 0, Enabled: 1, Display: 1, Description: "中奖记录管理", Url: "", Name: "中奖记录管理", Icon: "#xe614", Sort: 100},
		{Id: 150, Pid: 149, Enabled: 1, Display: 1, Description: "中奖记录", Url: "RewardLogIndexController.get", Name: "中奖记录", Icon: "", Sort: 100},
		{Id: 151, Pid: 149, Enabled: 1, Display: 0, Description: "删除中奖记录", Url: "RewardLogIndexController.Delone", Name: "删除中奖记录", Icon: "", Sort: 100},
		{Id: 152, Pid: 149, Enabled: 1, Display: 0, Description: "删除所有中奖记录", Url: "RewardLogIndexController.Delbatch", Name: "删除所有中奖记录", Icon: "", Sort: 100},
		{Id: 153, Pid: 149, Enabled: 1, Display: 0, Description: "派送中奖记录", Url: "RewardLogIndexController.Delivered", Name: "派送中奖记录", Icon: "", Sort: 100},
		{Id: 154, Pid: 149, Enabled: 1, Display: 0, Description: "批量派送中奖记录", Url: "RewardLogIndexController.Deliveredbatch", Name: "批量派送中奖记录", Icon: "", Sort: 100},


	}
	rolePermissions := []RolePermission{
		{Id: 200, RoleId: 2, PermissionId: 100},
		{Id: 201, RoleId: 2, PermissionId: 101},
		{Id: 202, RoleId: 2, PermissionId: 102},
		{Id: 203, RoleId: 2, PermissionId: 103},
		{Id: 204, RoleId: 2, PermissionId: 104},
		{Id: 205, RoleId: 2, PermissionId: 110},
		{Id: 206, RoleId: 2, PermissionId: 111},
		{Id: 207, RoleId: 2, PermissionId: 112},
		{Id: 208, RoleId: 2, PermissionId: 113},
		{Id: 209, RoleId: 2, PermissionId: 120},
		{Id: 210, RoleId: 2, PermissionId: 121},
		{Id: 211, RoleId: 2, PermissionId: 122},
		{Id: 212, RoleId: 2, PermissionId: 123},
		{Id: 213, RoleId: 2, PermissionId: 124},
		{Id: 214, RoleId: 2, PermissionId: 125},
		{Id: 215, RoleId: 2, PermissionId: 126},
		{Id: 216, RoleId: 2, PermissionId: 127},
		{Id: 217, RoleId: 2, PermissionId: 130},
		{Id: 218, RoleId: 2, PermissionId: 131},
		{Id: 219, RoleId: 2, PermissionId: 132},
		{Id: 220, RoleId: 2, PermissionId: 140},
		{Id: 221, RoleId: 2, PermissionId: 141},
		{Id: 222, RoleId: 2, PermissionId: 142},
		{Id: 223, RoleId: 2, PermissionId: 143},
		{Id: 229, RoleId: 2, PermissionId: 149},
		{Id: 230, RoleId: 2, PermissionId: 150},
		{Id: 231, RoleId: 2, PermissionId: 151},
		{Id: 232, RoleId: 2, PermissionId: 152},
		{Id: 233, RoleId: 2, PermissionId: 153},
		{Id: 234, RoleId: 2, PermissionId: 154},
		{Id: 235, RoleId: 2, PermissionId: 128},
		{Id: 236, RoleId: 2, PermissionId: 133},
	}
	o := orm.NewOrm()
	for _, v := range permisions {
		if _, _, err := o.ReadOrCreate(&v, "Id"); err != nil {
			beego.Error("InitProjectData Permission error", err)
		}
	}
	for _, v := range rolePermissions {
		if _, _, err := o.ReadOrCreate(&v, "Id"); err != nil {
			beego.Error("InitProjectData RolePermission error", err)
		}
	}
}
