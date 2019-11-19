package initial

import (
	"github.com/astaxie/beego"
	. "phagego/phagev2/utils"

	"phagego/phage-check-web/utils"
)

func initTemplateFunc() {
	beego.AddFuncMap("getSiteConfigCodeMap", GetSiteConfigCodeMap)
	beego.AddFuncMap("getDialStatusMap", utils.GetDialStatusNameMap)

}

func GetSiteConfigCodeMap() map[string]string {
	m := map[string]string{
		Scname: "站点名称"}
	return m
}
