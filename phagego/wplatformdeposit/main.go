package main

func main() {
	isValid := LicenseValidate("WPLATFORM_DEPOSIT")
	if isValid { // 验证通过
		mw := NewHomeGui()
		InitHomeGui(mw)
		mw.Run()
	}
}
