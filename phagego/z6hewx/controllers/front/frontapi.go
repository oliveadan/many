package front

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math/big"
	. "phagego/z6hewx/models"
	"strconv"

	"crypto/rand"
	"strings"
)

type FrontController struct {
	beego.Controller
}

func (this *FrontController) Get() {
	o := orm.NewOrm()
	//微信
	var wechats []Wechat
	var wechat Wechat
	num, err := o.QueryTable(new(Wechat)).Filter("Enabled", 1).All(&wechats)
	//获取所有的微信账号，根据总数随机出一个
	rand, err := rand.Int(rand.Reader, big.NewInt(num))
	rands, _ := strconv.Atoi(fmt.Sprintf("%d", rand))
	if err != nil {
		beego.Error("随机数失败", err)
	}
	//遍历所有微信账号，根据上面的随机数返回一个
	for i, v := range wechats {
		if i == rands {
			wechat = v
			break
		}
	}

	if err != nil {
		beego.Error("前台查询微信失败", err)
	}
	//中奖号码
	var winningmumbers []WinningNumbers
	_, err1 := o.QueryTable(new(WinningNumbers)).OrderBy("-Period").Limit(4).All(&winningmumbers)

	if err1 != nil {
		beego.Error("前台查询中奖号码失败", err)
	}

	for i, v := range winningmumbers {
		winningmumbers[i].Numbers = strings.Replace(v.Numbers, ",", "、", -1)

	}


	this.Data["wechat"] = wechat
	this.Data["winningmumbers"] = winningmumbers
	this.TplName = "front/index.html"
}


