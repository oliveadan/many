package index

import (
	"fmt"
	"net"
	"phagego/phagev2/controllers/sysmanage"
	"phagego/common/utils"
	"time"
)

type SysIndexController struct {
	sysmanage.BaseController
}

func (this *SysIndexController) Prepare() {
	this.EnableXSRF = false
}

func (this *SysIndexController) Get() {
	this.TplName = "sysmanage/index/index.html"
}

func (this *SysIndexController) Systeminfo() {
	var code int
	var msg string
	var data = make([]string, 0)
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &data)
	token := this.GetString("token")
	if token == "" {
		return
	}
	t := time.Now().Format("2006-01-02")
	if token != utils.Md5(t, utils.Pubsalt) {
		return
	}

	netInterfaces, err := net.Interfaces()
	if err != nil {
		msg = fmt.Sprintf("fail to get net interfaces: %v", err)
		return
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		data = append(data, fmt.Sprintf("%d,%s,%d,%s,%s", netInterface.MTU, netInterface.Flags.String(), netInterface.Index, netInterface.HardwareAddr.String(), netInterface.Name))
	}
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		msg = fmt.Sprintf("fail to get net InterfaceAddrs addrs: %v", err)
		return
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				data = append(data, ipNet.IP.String())
			}
		}
	}
	code = 1
}
