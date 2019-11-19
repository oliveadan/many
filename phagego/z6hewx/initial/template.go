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
		Scname: "站点名称"}
	return m
}
