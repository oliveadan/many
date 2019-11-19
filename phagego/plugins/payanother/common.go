package payanother

import (
	"bytes"
	"sort"
	"strings"
	comm "phagego/common/utils"
	"net/url"
)

const (
	PayAnother_huitao = "huitao"
)

func GetPayAnotherList() [][]string {
	return [][]string{
		{PayAnother_huitao, "汇淘支付"},
	}
}

type CommVo struct {

}

func GetPayAnotherVo(name string) interface{} {
	switch name {
	case PayAnother_huitao:
		return &CommVo{}
	}
	return nil
}

func GetBankMap(name string) map[string]string {
	switch name {
	case PayAnother_huitao:
		return map[string]string{"ABC": "BANK_ABC", "BOC": "BANK_BOC", "COMM": "BANK_BOCOM", "CCB": "BANK_CCB", "ICBC": "BANK_ICBC", "PSBC": "BANK_PSBC", "CMB": "BANK_CMB", "SPDB": "BANK_SPDB", "CEB": "BANK_CEB", "CITIC": "BANK_CITIC", "PAB": "BANK_PAB", "CMBC": "BANK_CMBC", "HXBANK": "BANK_HXBC", "GDB": "BANK_GDB", "BJBANK": "BANK_BOBJ", "SHBANK": "BANK_BOS", "CIB": "BANK_CIB"}
	}
	return map[string]string{}
}

func GetPayAnotherMap() map[string]string {
	var m = make(map[string]string)
	for _, v := range GetPayAnotherList() {
		m[v[0]] = v[1]
	}
	return m
}

func GenMd5Sign(key string, m *map[string]string, ignores ... string) string {
	var signData []string
	for k, v := range *m {
		isIgnore := false
		for _, tmp := range ignores {
			if tmp == k {
				isIgnore = true
				break
			}
		}
		if isIgnore {
			continue
		}
		if v != "" && k != "sign" && k != "key" {
			signData = append(signData, k+"="+v)
		}
	}
	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key
	//fmt.Println(signStr)
	return strings.ToUpper(comm.Md5(signStr))
}

func MapToUrl(v map[string]string, encode bool) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		if encode {
			k = url.QueryEscape(k)
			vs = url.QueryEscape(v[k])
		}
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		prefix := k + "="
		buf.WriteString(prefix)
		buf.WriteString(vs)
	}
	return buf.String()
}
