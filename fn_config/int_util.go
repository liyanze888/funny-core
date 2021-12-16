package fn_config

import (
	"fmt"
	"reflect"
	"strconv"
)

func getInt(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(int(getMustInt64Property(key)))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(int(getInt64PropertyOrDefault(key, defV)))
}

func getInt8(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(int8(getMustInt64Property(key)))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(int8(getInt64PropertyOrDefault(key, defV)))
}

func getInt16(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(int16(getMustInt64Property(key)))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(int16(getInt64PropertyOrDefault(key, defV)))
}

func getInt32(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(int32(getMustInt64Property(key)))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(int32(getInt64PropertyOrDefault(key, defV)))
}

func getInt64(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(getMustInt64Property(key))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(getInt64PropertyOrDefault(key, defV))
}

func getUint(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(uint(getMustUint64Property(key)))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(uint(getUint64PropertyOrDefault(key, defV)))
}

func getUint8(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(uint8(getMustUint64Property(key)))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(uint8(getUint64PropertyOrDefault(key, defV)))
}

func getUint16(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(uint16(getMustUint64Property(key)))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(uint16(getUint64PropertyOrDefault(key, defV)))
}

func getUint32(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(uint32(getMustUint64Property(key)))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(uint32(getUint64PropertyOrDefault(key, defV)))
}

func getUint64(key string, field reflect.StructField) reflect.Value {
	if _, ok := field.Tag.Lookup("notnull"); ok {
		if _, ok := field.Tag.Lookup("defVal"); !ok {
			return reflect.ValueOf(getMustUint64Property(key))
		}
	}
	defV := "0"
	if defVal, ok := field.Tag.Lookup("defVal"); ok {
		if defVal != "" {
			defV = defVal
		}
	}
	return reflect.ValueOf(getUint64PropertyOrDefault(key, defV))
}

func getMustInt64Property(key string) int64 {
	value := confViper.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("Need %s Config Property", key))
	}

	if atoi, err := strconv.ParseInt(value, 10, 64); err != nil {
		panic(fmt.Sprintf("Format Error %s Config Property", key))
	} else {
		return atoi
	}
	return 0
}

func getInt64PropertyOrDefault(key string, defVal string) int64 {
	value := confViper.GetString(key)
	if value == "" {
		if defVal == "" {
			return 0
		}
		if atoi, err := strconv.ParseInt(defVal, 10, 64); err != nil {
			panic(fmt.Sprintf("Format Error %s Config Property", key))
		} else {
			return atoi
		}
		return 0
	}
	if atoi, err := strconv.ParseInt(value, 10, 64); err != nil {
		panic(fmt.Sprintf("Format Error %s Config Property", key))
	} else {
		return atoi
	}
}

func getMustUint64Property(key string) uint64 {
	value := confViper.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("Need %s Config Property", key))
	}

	if atoi, err := strconv.ParseUint(value, 10, 64); err != nil {
		panic(fmt.Sprintf("Format Error %s Config Property", key))
	} else {
		return atoi
	}
	return 0
}

func getUint64PropertyOrDefault(key string, defVal string) uint64 {
	value := confViper.GetString(key)
	if value == "" {
		if defVal == "" {
			return 0
		}
		if atoi, err := strconv.ParseUint(defVal, 10, 64); err != nil {
			panic(fmt.Sprintf("Format Error %s Config Property", key))
		} else {
			return atoi
		}
		return 0
	}
	if atoi, err := strconv.ParseUint(value, 10, 64); err != nil {
		panic(fmt.Sprintf("Format Error %s Config Property", key))
	} else {
		return atoi
	}
}
