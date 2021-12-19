package fn_config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func getMap(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return getMustMapProperty(key, field.Type)
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal[1 : len(defVal)-1]
		}
	}
	return getMapProperty(key, field.Type, defV)
}

func getMustMapProperty(key string, fieldType reflect.Type) reflect.Value {
	stringMap := confViper.GetStringMapString(key)
	if len(stringMap) == 0 {
		panic(fmt.Sprintf("Need %s Config Property", key))
	}
	return getMapProperty(key, fieldType, "")
}

func getMapProperty(key string, fieldType reflect.Type, defVal string) reflect.Value {
	switch fieldType.Elem().Kind() {
	case reflect.Slice:
		return getStringMapSliceValues(key, fieldType, defVal)
	}
	stringMap := confViper.GetStringMapString(key)
	if len(stringMap) == 0 && len(defVal) == 0 {
		return reflect.MakeMap(fieldType)
	}

	if len(stringMap) == 0 && len(defVal) > 0 {
		stringMap = parseBaseKeyValue(key, defVal, fieldType.Elem())
	}

	rMap := reflect.MakeMap(fieldType)
	for keyV, value := range stringMap {
		rMap.SetMapIndex(getKeyVal(key, keyV, fieldType.Key()), getValueVal(key, value, fieldType.Elem()))
	}
	return rMap
}

func getStringMapSliceValues(key string, fieldType reflect.Type, defVal string) reflect.Value {
	sliceMap := confViper.GetStringMapStringSlice(key)
	if len(sliceMap) == 0 {
		sliceMap = parseSliceDefalutMap(key, defVal, fieldType.Elem())
	}

	rMap := reflect.MakeMap(fieldType)
	for keyVal, value := range sliceMap {
		rMap.SetMapIndex(getKeyVal(key, keyVal, fieldType.Key()), getSliceVal(key, value, fieldType.Elem()))
	}
	return rMap
}

func parseSliceDefalutMap(key string, defVal string, valType reflect.Type) map[string][]string {
	m := make(map[string][]string)
	for {
		index := strings.Index(defVal, ":")
		if index < 0 {
			return m
		}
		keyVal := defVal[:index]
		if strings.HasPrefix(keyVal, "\"") {
			keyVal = keyVal[1 : len(keyVal)-1]
		}
		defVal = defVal[index+1:]

		sliceIndex := strings.Index(defVal, "],")
		sliceVal := ""
		if sliceIndex >= 0 {
			sliceVal = defVal[:sliceIndex]
		} else {
			sliceVal = defVal[:len(defVal)-1]
		}
		defVal = defVal[sliceIndex+1:]
		sliceVal = sliceVal[1:]
		switch valType.Elem().Kind() {
		case reflect.String:
			values := make([]string, 0)
			for {
				if strings.HasPrefix(sliceVal, "\"") {
					sliceVal = sliceVal[1:]
					idx := strings.Index(sliceVal, "\"")
					if idx >= 0 {
						values = append(values, strings.TrimSpace(sliceVal[:idx]))
					} else {
						panic(fmt.Sprintf("Unsupport key type Error key = %s", key))
					}
					sliceVal = sliceVal[idx+1:]
					if len(sliceVal) <= 0 {
						m[keyVal] = values
						break
					} else {
						dotIdx := strings.Index(sliceVal, ",")
						if dotIdx >= 0 {
							sliceVal = sliceVal[dotIdx+1:]
						} else {
							m[keyVal] = values
							break
						}
					}
				} else {
					idx := strings.Index(sliceVal, ",")
					if idx > 0 {
						values = append(values, sliceVal[:idx])
					} else {
						values = append(values, strings.TrimSpace(sliceVal))
					}
					m[keyVal] = values
					break
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Bool:
			split := strings.Split(sliceVal, ",")
			m[keyVal] = split
		default:
			panic(fmt.Sprintf("Unsupport key type Error key = %s", key))
		}

		dotIdx := strings.Index(defVal, ",")
		if dotIdx >= 0 {
			defVal = strings.TrimSpace(defVal[dotIdx+1:])
		}
	}
	return m
}

func parseBaseKeyValue(key string, defVal string, valType reflect.Type) map[string]string {
	m := map[string]string{}
	for {
		index := strings.Index(defVal, ":")
		if index < 0 {
			return m
		}
		keyVal := defVal[:index]
		if strings.HasPrefix(keyVal, "\"") {
			keyVal = keyVal[1 : len(keyVal)-1]
		}

		defVal = defVal[index+1:]

		indexSubStr := ","
		subIndex := 1
		prefix := ""
		switch valType.Kind() {
		case reflect.String:
			if strings.HasPrefix(defVal, "\"") {
				indexSubStr = "\","
				prefix = "\""
				subIndex = 2
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Bool:
		default:
			panic(fmt.Sprintf("Unsupport key type Error key = %s", key))
		}

		valueValIndex := strings.Index(defVal, indexSubStr)
		valueVal := ""
		if valueValIndex < 0 {
			if len(prefix) == 0 {
				valueVal = defVal
			} else {
				valueVal = defVal[1 : len(defVal)-1]
			}
		} else {
			if len(prefix) != 0 {
				valueVal = defVal[1:valueValIndex]
			} else {
				valueVal = defVal[:valueValIndex]
			}
		}
		//if len(prefix) != 0 {
		//	valueVal = valueVal[1 : len(valueVal)-1]
		//}
		defVal = defVal[valueValIndex+subIndex:]
		m[keyVal] = valueVal

	}
	return m
}

func getSliceVal(key string, value []string, elemType reflect.Type) reflect.Value {
	slice := reflect.MakeSlice(elemType, 0, 0)
	for _, data := range value {
		val := getValueVal(key, data, elemType.Elem())
		slice = reflect.Append(slice, val)
	}
	return slice
}

func getKeyVal(key string, keyVal string, keyType reflect.Type) reflect.Value {
	val, b := getBaseTypeVal(key, keyVal, keyType)
	if b {
		return val
	}
	panic(fmt.Sprintf("Unsupport key type Error key = %s", key))
}

func getValueVal(key string, valueVal string, valueType reflect.Type) reflect.Value {
	val, b := getBaseTypeVal(key, valueVal, valueType)
	if b {
		return val
	}
	switch valueType.Kind() {
	case reflect.Bool:
		parseBool, err := strconv.ParseBool(strings.TrimSpace(valueVal))
		if err != nil {
			panic(fmt.Sprintf("Unsupport value type Error key = %s", key))
		}
		return reflect.ValueOf(parseBool)
	}
	panic(fmt.Sprintf("Unsupport value type Error key = %s", key))
}

func getBaseTypeVal(key string, val string, vType reflect.Type) (reflect.Value, bool) {
	switch vType.Kind() {
	case reflect.String:
		return reflect.ValueOf(strings.TrimSpace(val)), true
	case reflect.Int:
		return reflect.ValueOf(int(convertInt64(key, val))), true
	case reflect.Int8:
		return reflect.ValueOf(int8(convertInt64(key, val))), true
	case reflect.Int16:
		return reflect.ValueOf(int16(convertInt64(key, val))), true
	case reflect.Int32:
		return reflect.ValueOf(int32(convertInt64(key, val))), true
	case reflect.Int64:
		return reflect.ValueOf(convertInt64(key, val)), true
	case reflect.Uint:
		return reflect.ValueOf(uint(convertInt64(key, val))), true
	case reflect.Uint8:
		return reflect.ValueOf(uint8(convertInt64(key, val))), true
	case reflect.Uint16:
		return reflect.ValueOf(uint16(convertInt64(key, val))), true
	case reflect.Uint32:
		return reflect.ValueOf(uint32(convertInt64(key, val))), true
	case reflect.Uint64:
		return reflect.ValueOf(uint64(convertInt64(key, val))), true
	}
	return reflect.ValueOf(""), false
}

func convertInt64(key string, value string) int64 {
	if atoi, err := strconv.ParseInt(value, 10, 64); err != nil {
		panic(fmt.Sprintf("Format Error %s Config Property", key))
	} else {
		return atoi
	}
}
