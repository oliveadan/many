package routers

import (
	"github.com/astaxie/beego"
	_ "phagego/frameweb-v2/routers"
	"phagego/phage-vote-api/controllers/common/setvote"
	"phagego/phage-vote-api/controllers/common/votedetail"
	"phagego/phage-vote-api/controllers/front/voteapi"
)

func init() {
	//前端
	beego.Router("/", &voteapi.IndexVoteApiIndexController{})
	beego.Router("/vote", &voteapi.IndexVoteApiIndexController{}, "post:Vote")
	beego.Router("/queryvote", &voteapi.IndexVoteApiIndexController{}, "post:AllVote")
	// 后台管理系统
	var adminRouter string = beego.AppConfig.String("adminrouter")
	beego.Router(adminRouter+"/votedetail", &votedetail.IndexVoteDetailController{})
	beego.Router(adminRouter+"/delallvotedetail", &votedetail.IndexVoteDetailController{}, "post:DelBtch")
	beego.Router(adminRouter+"/setvoteindex", &setvote.IndexSetVoteController{})
	beego.Router(adminRouter+"/setvoteedit", &setvote.EditSetVoteController{})
}
