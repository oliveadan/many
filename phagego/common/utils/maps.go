package utils

import (
	"strings"
	"sort"
	"bytes"
	"net/url"
	"encoding/json"
	"reflect"
)

// 将map按key的字典序排序后拼接
// 参数顺序 ignoreKeys、ignoreEmpty（Y/N）、sep、urlencode
func JoinKv(m map[string]string, args ...string) string {
	var ignoreKeys string
	var ignoreEmpty bool = true
	var sep string = "&"
	var urlencode bool = false
	for i, v := range args {
		switch i {
		case 0:
			ignoreKeys = v
		case 1:
			if v == "N" {
				ignoreEmpty = false
			}
		case 2:
			sep = v
		case 3:
			if v == "Y" {
				urlencode = true
			}
		}
	}

	keys := make([]string, len(m))
	i := 0
	for k, v := range m {
		if ignoreKeys != "" && strings.Contains(ignoreKeys, k) {
			continue
		}
		if ignoreEmpty && v == "" {
			continue
		}
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	b := &bytes.Buffer{}
	if urlencode {
		for _, k := range keys {
			if k == "" {
				continue
			}
			b.WriteString(sep)
			b.WriteString(url.QueryEscape(k))
			b.WriteString("=")
			b.WriteString(url.QueryEscape(m[k]))
		}
	} else {
		for _, k := range keys {
			if k == "" {
				continue
			}
			b.WriteString(sep)
			b.WriteString(k)
			b.WriteString("=")
			b.WriteString(m[k])
		}
	}
	return strings.Replace(b.String(), sep, "", 1)
}

// 将map[string]string转化为结构体
func MapStringToStruct(m map[string]string, i interface{}) error {
	bin, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, i)
	if err != nil {
		return err
	}
	return nil
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}


// 根据2维数组获取map
func List2Map(list [][]string) map[string]string {
	var m = make(map[string]string)
	for _, v := range list {
		m[v[0]] = v[1]
	}
	return m
}
