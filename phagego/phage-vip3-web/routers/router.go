package routers

import (
	"github.com/astaxie/beego"
	_ "phagego/frameweb-v2/routers"
	"phagego/phage-vip3-web/controllers/common/level"
	"phagego/phage-vip3-web/controllers/common/lucky"
	"phagego/phage-vip3-web/controllers/common/membersingle"
	"phagego/phage-vip3-web/controllers/common/membertotal"
	"phagego/phage-vip3-web/controllers/common/period"
	"phagego/phage-vip3-web/controllers/front"
	"phagego/phage-vip3-web/controllers/rewardlog"
)

func init() {
	// 后台管理系统
	var adminRouter string = beego.AppConfig.String("adminrouter")

	//前端
	beego.Router("/", &front.FrontIndexController{})
	beego.Router("/query", &front.FrontIndexController{}, "post:Query")
	beego.Router("/getgift", &front.FrontIndexController{}, "post:GetGift")
	beego.Router("/repairmembertotal", &front.FrontIndexController{}, "get:RepairMemberTotal")

	//vip等级
	beego.Router(adminRouter+"/level/index", &level.LevelController{})
	beego.Router(adminRouter+"/level/add", &level.LevelAddController{})
	beego.Router(adminRouter+"/level/edit", &level.LevelEditController{})
	beego.Router(adminRouter+"/level/delone", &level.LevelController{}, "post:Delone")
	//好运金
	beego.Router(adminRouter+"/lucky/index", &lucky.LuckyController{})
	beego.Router(adminRouter+"/lucky/add", &lucky.LuckyAddController{})
	beego.Router(adminRouter+"/lucky/edit", &lucky.LuckyEditController{})
	beego.Router(adminRouter+"/lucky/delone", &lucky.LuckyController{}, "post:Delone")
	//周期分类
	beego.Router(adminRouter+"/period/index", &period.PeriodIndexController{})
	beego.Router(adminRouter+"/period/add", &period.PeriodAddController{})
	beego.Router(adminRouter+"/period/edit", &period.PeriodEditController{})
	beego.Router(adminRouter+"/period/delone", &period.PeriodIndexController{}, "post:Delone")
	//单期投注
	beego.Router(adminRouter+"/membersingle/index", &membersingle.MembersingleIndexController{})
	beego.Router(adminRouter+"/membersingle/add", &membersingle.MembersingleAddController{})
	beego.Router(adminRouter+"/membersingle/edit", &membersingle.MembersingleEditController{})
	beego.Router(adminRouter+"/membersingle/delone", &membersingle.MembersingleIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/membersingle/delbatch", &membersingle.MembersingleIndexController{}, "post:DelBatch")
	beego.Router(adminRouter+"/membersingle/import", &membersingle.MembersingleIndexController{}, "post:Import")
	beego.Router(adminRouter+"/membersingle/countgift", &membersingle.MembersingleIndexController{}, "post:CountGift")
	beego.Router(adminRouter+"/membersingle/countluckygift", &membersingle.MembersingleIndexController{}, "post:CountLuckyGift")
	beego.Router(adminRouter+"/membersingle/export", &membersingle.MembersingleIndexController{}, "post:Export")
	//会员统计
	beego.Router(adminRouter+"/membertotal/index", &membertotal.MemberTotalIndexController{})
	beego.Router(adminRouter+"/membertotal/delbatch", &membertotal.MemberTotalIndexController{}, "post:Delbatch")
	beego.Router(adminRouter+"/membertotal/export", &membertotal.MemberTotalIndexController{}, "post:Export")
	beego.Router(adminRouter+"/membertotal/edit", &membertotal.MemberTotalEditController{})

	//中奖记录
	beego.Router(adminRouter+"/rewardlog/index", &rewardlog.RewardLogIndexController{})
	beego.Router(adminRouter+"/rewardlog/delone", &rewardlog.RewardLogIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/rewardlog/Delbatch", &rewardlog.RewardLogIndexController{}, "post:Delbatch")
	beego.Router(adminRouter+"/rewardlog/export", &rewardlog.RewardLogIndexController{}, "post:Export")
	beego.Router(adminRouter+"/rewardlog/delivered", &rewardlog.RewardLogIndexController{}, "post:Delivered")
	beego.Router(adminRouter+"/rewardlog/deliveredbatch", &rewardlog.RewardLogIndexController{}, "post:Deliveredbatch")
}