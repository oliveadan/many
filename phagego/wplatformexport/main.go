package main

import "phagego/framewin/license"

const softName = "WPLATFORM_EXPORTDATA"
const licenseSalt = ""

func main() {
	if license.Validate(softName, licenseSalt) { // 验证通过
		mw := NewHomeGui()
		InitHomeGui(mw)
		mw.Run()
	}
}
