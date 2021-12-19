package fn_config_test

import (
	"github.com/liyanze888/funny-core/fn_factory"
	"log"
	"testing"
)

type GlobalConfig struct {
	prefix     string `value:"robot"`
	UserIds    []int64
	UserIdsStr []string
	UserStatus []bool
	UserDatas  []string `defVal:"[hello,world,ttt]"`
}

type UserConfig struct {
	Age             int
	Name            string
	RefTex          map[string]int64
	RefTexBool      map[string]bool
	RefIntSlice     map[string][]int64
	RefTxtSlice     map[string][]string
	RefBoolSlice    map[string][]bool
	DefIntSlice     map[int64]bool     `defVal:"{1:false,5:true}"`
	DefStringSlice  map[int64]string   `defVal:"{1:\"user\",5:\"name\"}"`
	DefMapSlice     map[int64][]string `defVal:"{1:[\"user\" , name],5:[\"name\" , user]}"`
	DefMapBoolSlice map[string][]bool  `defVal:"{user:[true , false],\"name\":[false , true]}"`
}

type User struct {
	Config     *GlobalConfig `autowire:""`
	UserConfig *UserConfig   `autowire:""`
}

func TestConfig(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.Ltime)
	fn_factory.BeanFactory.RegisterConfigBean(&GlobalConfig{})
	fn_factory.BeanFactory.RegisterConfigBean(&UserConfig{})
	u := &User{}
	fn_factory.BeanFactory.RegisterBean(u)
	fn_factory.BeanFactory.StartUp()

	log.Printf("u.config = %v", u.Config)
	log.Printf("u.UserConfig = %v", u.UserConfig)
}
