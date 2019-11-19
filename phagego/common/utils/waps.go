package utils

import "strings"

func IsWap(ua string) bool {
	ua = strings.ToLower(ua)
	//常见的手机端UA判断，有即默认为手机端 pad、phone关键字替代ipad，iphone,apad等
	mobile := []string{"mobile", "phone", "android", "pad", "pod", "symbian", "wap", "smartphone", "apk", "ios"}
	for _, m := range mobile {
		if strings.Contains(ua, m) {
			return true
		}
	}
	//生僻的不常见的UA判断
	mbstr := "w3c,acs-,alav,alca,amoi,audi,avan,benq,bird,blac"
	mbstr += "blaz,brew,cell,cldc,cmd-,dang,doco,eric,hipt,inno"
	mbstr += "ipaq,java,jigs,kddi,keji,leno,lg-c,lg-d,lg-g,"
	mbstr += "maui,maxo,midp,mits,mmef,mobi,mot-,moto,mwbp,"
	mbstr += "newt,noki,oper,palm,pana,pant,phil,play,port,prox"
	mbstr += "qwap,sage,sams,sany,sch-,sec-,send,seri,sgh-,shar"
	mbstr += "sie-,siem,smal,smar,sony,sph-,symb,t-mo,teli,"
	mbstr += "tosh,tsm-,upg1,upsi,vk-v,voda,wap-,wapa,wapi,wapp"
	mbstr += "wapr,webc,winw,winw,xda,xda-,up.browser,up.link,mmp,midp,xoom"
	slice := strings.Split(mbstr, ",")
	for _, m := range slice {
		if strings.Contains(ua, m) {
			return true
		}
	}
	return false
}
