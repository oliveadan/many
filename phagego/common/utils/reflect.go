package utils

import "reflect"

//定义注册结构map
var registerStructMaps = make(map[string]reflect.Type)

//根据name初始化结构
//在这里根据结构的成员注解进行DI注入，这里没有实现，只是简单都初始化
func ReflectNew(name string) (interface{}, bool) {
	if v, ok := registerStructMaps[name]; ok {
		c := reflect.New(v).Interface()
		return c, true
	} else {
		return nil, false
	}
}

//根据名字注册实例
func ReflectRegister(name string, c interface{}) {
	registerStructMaps[name] = reflect.TypeOf(c).Elem()
}
