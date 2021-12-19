package fn_config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func getSlice(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		// 需要改造  如果默认值 则must 里面不报错
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return getMustSliceProperty(key, field.Type)
		}
	}

	defV := "[]"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal[1 : len(defVal)-1]
		}
	}
	return getSliceProperty(key, field.Type, defV)
}

func getMustSliceProperty(key string, fieldType reflect.Type) reflect.Value {
	switch fieldType.Elem().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return getMustIntSliceProperty(key, fieldType)
	case reflect.String:
		return getMustStringSliceProperty(key, fieldType)
	case reflect.Bool:
		return getMustBoolSliceProperty(key, fieldType)

	}
	return reflect.MakeSlice(fieldType, 0, 0)
}

func getMustIntSliceProperty(key string, fieldType reflect.Type) reflect.Value {
	confIntSlice := confViper.GetIntSlice(key)
	if len(confIntSlice) == 0 {
		panic(fmt.Sprintf("Need %s Config Property", key))
	}

	slice := reflect.MakeSlice(fieldType, 0, 0)
	for _, v := range confIntSlice {
		slice = reflect.Append(slice, getTargetValue(v, fieldType.Elem()))
	}
	return slice
}

func getMustStringSliceProperty(key string, fieldType reflect.Type) reflect.Value {
	confStringSlice := confViper.GetStringSlice(key)
	if len(confStringSlice) == 0 {
		panic(fmt.Sprintf("Need %s Config Property", key))
	}
	slice := reflect.MakeSlice(fieldType, 0, 0)
	for _, v := range confStringSlice {
		slice = reflect.Append(slice, reflect.ValueOf(v))
	}
	return slice
}

func getMustBoolSliceProperty(key string, fieldType reflect.Type) reflect.Value {
	confStringSlice := confViper.GetStringSlice(key)
	if len(confStringSlice) == 0 {
		panic(fmt.Sprintf("Need %s Config Property", key))
	}
	slice := reflect.MakeSlice(fieldType, 0, 0)
	for _, v := range confStringSlice {
		bVal, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("Format Error %s Config Property", key))
		}
		slice = reflect.Append(slice, reflect.ValueOf(bVal))
	}
	return slice
}

func getSliceProperty(key string, fieldType reflect.Type, defVal string) reflect.Value {
	switch fieldType.Elem().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return getIntSliceProperty(key, fieldType, defVal)
	case reflect.String:
		return getStringSliceProperty(key, fieldType, defVal)
	case reflect.Bool:
		return getBoolSliceProperty(key, fieldType, defVal)
	}
	return reflect.MakeSlice(fieldType, 0, 0)
}

func getIntSliceProperty(key string, fieldType reflect.Type, defVal string) reflect.Value {
	confIntSlice := confViper.GetIntSlice(key)
	if len(confIntSlice) == 0 {
		if len(defVal) > 0 {
			return parseIntSliceDefault(key, fieldType, defVal)
		}
		return reflect.MakeSlice(fieldType, 0, 0)
	}
	slice := reflect.MakeSlice(fieldType, 0, 0)
	for _, v := range confIntSlice {
		slice = reflect.Append(slice, getTargetValue(v, fieldType.Elem()))
	}
	return slice
}

func getStringSliceProperty(key string, fieldType reflect.Type, defVal string) reflect.Value {
	slice := reflect.MakeSlice(fieldType, 0, 0)
	confStringSlice := confViper.GetStringSlice(key)
	if len(confStringSlice) == 0 {
		if len(defVal) > 0 {
			return parseStringSliceDefault(fieldType, defVal)
		}
		return reflect.MakeSlice(fieldType, 0, 0)
	}

	for _, v := range confStringSlice {
		slice = reflect.Append(slice, reflect.ValueOf(v))
	}
	return slice
}

func getBoolSliceProperty(key string, fieldType reflect.Type, defVal string) reflect.Value {
	confStringSlice := confViper.GetStringSlice(key)
	if len(confStringSlice) == 0 {
		if len(defVal) > 0 {
			return parseBoolSliceDefault(key, fieldType, defVal)
		}
		return reflect.MakeSlice(fieldType, 0, 0)
	}

	slice := reflect.MakeSlice(fieldType, 0, 0)
	for _, v := range confStringSlice {
		bVal, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("Format Error %s Config Property", key))
		}
		slice = reflect.Append(slice, reflect.ValueOf(bVal))
	}
	return slice
}

func parseBoolSliceDefault(key string, fieldType reflect.Type, defVal string) reflect.Value {
	slice := reflect.MakeSlice(fieldType, 0, 0)
	confDatas := strings.Split(defVal, ",")
	for _, data := range confDatas {
		bVal, err := strconv.ParseBool(strings.TrimSpace(data))
		if err != nil {
			panic(fmt.Sprintf("Format Error %s Config Property", key))
		}
		slice = reflect.Append(slice, reflect.ValueOf(bVal))
	}
	return slice
}

func parseIntSliceDefault(key string, fieldType reflect.Type, defVal string) reflect.Value {
	confDatas := strings.Split(defVal, ",")
	slice := reflect.MakeSlice(fieldType, 0, 0)
	for _, data := range confDatas {
		atoi, err := strconv.Atoi(data)
		if err != nil {
			panic(fmt.Sprintf("Format Error %s Config Property", key))
		}
		reflect.Append(slice, getTargetValue(atoi, fieldType.Elem()))
	}
	return slice
}

func parseStringSliceDefault(fieldType reflect.Type, defVal string) reflect.Value {
	confDatas := strings.Split(defVal, ",")
	slice := reflect.MakeSlice(fieldType, 0, 0)
	for _, data := range confDatas {
		slice = reflect.Append(slice, reflect.ValueOf(data))
	}
	return slice
}

func getTargetValue(s int, elemType reflect.Type) reflect.Value {
	switch elemType.Kind() {
	case reflect.Int:
		return reflect.ValueOf(s)
	case reflect.Int8:
		return reflect.ValueOf(int8(s))
	case reflect.Int16:
		return reflect.ValueOf(int16(s))
	case reflect.Int32:
		return reflect.ValueOf(int32(s))
	case reflect.Int64:
		return reflect.ValueOf(int64(s))
	case reflect.Uint:
		return reflect.ValueOf(uint(s))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(s))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(s))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(s))
	case reflect.Uint64:
		return reflect.ValueOf(uint64(s))
	}
	panic(fmt.Sprintf("Config Unkonw Type Error %v", elemType))
}
