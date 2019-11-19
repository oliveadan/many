package initial

import (
	"github.com/astaxie/beego"
	. "phagego/frameweb-v2/utils"
	"phagego/phage-vip-web/utils"
)

func initTemplateFunc() {
	beego.AddFuncMap("getSiteConfigCodeMap", GetSiteConfigCodeMap)
}

func GetSiteConfigCodeMap() map[string]string {
	m := map[string]string{
		Scname: "站点名称",
		utils.Scofficial: "官网网址",
		utils.Scranking : "排行榜网址",
	    utils.Scregister: "官网注册",
        utils.Sccust    : "在线客服",
	    utils.Scfqa     : "博彩责任",
	    utils.Scpromotion: "优惠活动"}
	return m
}
