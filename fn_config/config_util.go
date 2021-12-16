package fn_config

import (
	"github.com/spf13/viper"
	"os"
	"reflect"
	"sync"
)

var confViper = viper.New()
var once sync.Once

var alreadyInitProperty = make(map[string]interface{})

func InitConfig(cf interface{}) {
	confViper.AutomaticEnv()
	if _, err := os.Open("config.yaml"); err == nil {
		confViper.SetConfigFile(getStringConfigOrDefault("config_file", "config.yaml"))
		err = confViper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}
	initProperties(cf)
}

func initProperties(i interface{}) {
	iType := reflect.TypeOf(i)
	iVal := reflect.ValueOf(i)
	tType := iType
	if iType.Kind() == reflect.Ptr {
		tType = iType.Elem()
	}

	if iVal.Kind() == reflect.Ptr {
		iVal = iVal.Elem()
	}

	fieldSize := tType.NumField()
	for i := 0; i < fieldSize; i++ {
		tField := tType.Field(i)
		key := tField.Name

		if tVal, ok := tField.Tag.Lookup("value"); ok {
			if tVal != "" {
				key = tVal
			}
		}
		vField := iVal.Field(i)
		property := getProperty(key, tField)
		if property != nil {
			if _, ok := alreadyInitProperty[key]; !ok {
				vField.Set(*property)
				alreadyInitProperty[key] = (*property).String()
			}
		}
	}

}

func getProperty(key string, field reflect.StructField) *reflect.Value {
	switch field.Type.Kind() {
	case reflect.String:
		ret := getString(key, field)
		return &ret
	case reflect.Int:
		ret := getInt(key, field)
		return &ret
	case reflect.Int8:
		ret := getInt8(key, field)
		return &ret
	case reflect.Int16:
		ret := getInt16(key, field)
		return &ret
	case reflect.Int32:
		ret := getInt32(key, field)
		return &ret
	case reflect.Int64:
		ret := getInt64(key, field)
		return &ret
	case reflect.Uint:
		ret := getUint(key, field)
		return &ret
	case reflect.Uint8:
		ret := getUint8(key, field)
		return &ret
	case reflect.Uint16:
		ret := getUint16(key, field)
		return &ret
	case reflect.Uint32:
		ret := getUint32(key, field)
		return &ret
	case reflect.Uint64:
		ret := getUint64(key, field)
		return &ret
	case reflect.Bool:
		ret := getBool(key, field)
		return &ret

	}
	return nil
}
