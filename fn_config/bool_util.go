package fn_config

import (
	"fmt"
	"reflect"
	"strconv"
)

func getBool(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if defVal, ok := field.Tag.Lookup("defVal"); ok {
			return reflect.ValueOf(getBoolOrDefaultProperty(key, defVal))
		}
		return reflect.ValueOf(getMustBoolProperty(key))
	}
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		return reflect.ValueOf(getBoolOrDefaultProperty(key, defVal))
	}
	return reflect.ValueOf(getBoolOrDefaultProperty(key, "false"))
}

// 获取 必须的字符串
func getMustBoolProperty(key string) bool {
	value := confViper.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("Need %s Config Property", key))
	}
	if ret, err := strconv.ParseBool(value); err != nil {
		panic(fmt.Sprintf("Format error %s Config Property", key))
	} else {
		return ret
	}
}

//获取字符串默认配置
func getBoolOrDefaultProperty(key string, defVal string) bool {
	value := confViper.GetString(key)
	if value == "" {
		if defVal == "" {
			return false
		} else {
			if ret, err := strconv.ParseBool(defVal); err != nil {
				panic(fmt.Sprintf("Format error %s Config Property", key))
			} else {
				return ret
			}
		}
	}

	if ret, err := strconv.ParseBool(value); err != nil {
		panic(fmt.Sprintf("Format error %s Config Property", key))
	} else {
		return ret
	}
}
