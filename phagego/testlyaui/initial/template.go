package initial

import (
	"github.com/astaxie/beego"
	. "phagego/frameweb-v2/utils"
)

func initTemplateFunc() {
	beego.AddFuncMap("getSiteConfigCodeMap", GetSiteConfigCodeMap)
}

func GetSiteConfigCodeMap() map[string]string {
	m := map[string]string{
		"DIY":  "自定义",
		Scname: "站点名称"}
	return m
}
