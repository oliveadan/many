package initial

import (
	. "phagego/phagev2/initial"
	// . "phagego/phagev2/models"
)

func init() {
	InitLog()
	InitSql()
	InitBeeCache()
	InitFilter()
	InitMailConf()
	InitSysTemplateFunc()
	initTemplateFunc()
}
