package utils

import "github.com/astaxie/beego"

func RetServeJson(c *beego.Controller, msg *string, code *int, data ...interface{}) {
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	ret["data"] = data
	c.Data["json"] = &ret
	c.ServeJSON()
}
