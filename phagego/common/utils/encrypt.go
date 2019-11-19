package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strconv"
	"fmt"
	"sort"
	"strings"
)

const Pubsalt = "phage2"

//md5方法
func Md5(ss ...string) string {
	h := md5.New()
	for _, s := range ss {
		if _, err := h.Write([]byte(s)); err != nil {
			return ""
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

//Guid方法
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

func GenMd5Sign(key string, m *map[string]interface{}) string {
	var signData []string
	for k, v := range *m {
		if v != nil && v != "" && k != "sign" && k != "key" {
			switch v.(type) {
			case string:
				signData = append(signData, k+"="+v.(string))
			case float64:
				signData = append(signData, k+"="+strconv.FormatFloat(v.(float64), 'f', 0, 64))
			case float32:
				signData = append(signData, k+"="+strconv.FormatFloat(float64(v.(float32)), 'f', 0, 64))
			case int:
				signData = append(signData, k+"="+strconv.Itoa(v.(int)))
			case int64:
				signData = append(signData, k+"="+strconv.FormatInt(v.(int64), 10))
			default:
				signData = append(signData, fmt.Sprintf("%s=%v", k, v))
			}
		}
	}
	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key
	//fmt.Println(signStr)
	return strings.ToUpper(Md5(signStr))
}

func GenMd5Sign2(key string, m *map[string]string) string {
	var signData []string
	for k, v := range *m {
		if v != "" && k != "sign" && k != "key" {
			signData = append(signData, k+"="+v)
		}
	}
	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key
	//fmt.Println(signStr)
	return strings.ToUpper(Md5(signStr))
}
