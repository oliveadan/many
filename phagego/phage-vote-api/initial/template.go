package initial

import (
	. "phagego/frameweb-v2/utils"
	"github.com/astaxie/beego"
	"phagego/phage-vote-api/utils"
)

func initTemplateFunc() {
	beego.AddFuncMap("getSiteConfigCodeMap", GetSiteConfigCodeMap)
	beego.AddFuncMap("getAnimal", utils.GetAnimal)
}

func GetSiteConfigCodeMap() map[string]string {
	m := map[string]string{
		"DIY":  "自定义",
		Scname: "站点名称"}
	return m
}
