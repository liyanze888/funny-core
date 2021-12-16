package fn_config

import (
	"fmt"
	"reflect"
)

func getString(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if defVal, ok := field.Tag.Lookup("defVal"); ok {
			return reflect.ValueOf(getStringConfigOrDefault(key, defVal))
		}
		return reflect.ValueOf(getMustString(key))
	}
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		return reflect.ValueOf(getStringConfigOrDefault(key, defVal))
	}
	return reflect.ValueOf(getStringConfigOrDefault(key, ""))
}

// 获取 必须的字符串
func getMustString(key string) string {
	value := confViper.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("Need %s Config Property", key))
	}
	return value
}

//获取字符串默认配置
func getStringConfigOrDefault(key string, defVal string) string {
	value := confViper.GetString(key)
	if value == "" {
		return defVal
	}
	return value
}
