package fn_factory

import (
	"fmt"
	"reflect"
	"strings"
)

var BeanFactory = newBeanFactory()

// Bean 的初始化状态值
//
const (
	Uninitialized = iota
	Initializing
	Initialized
)

type MyBean interface{}

//Bean 初始化接口
type MyBeanInitialization interface {
	InitBean(b BeanContext)
}

type MyBeanPostInitialization interface {
	PostInitilization()
}

type MyBeanDefinition struct {
	Init  int
	Bean  MyBean
	Name  string
	Type  reflect.Type
	Value reflect.Value
}

type BeanContext interface {
	RegisterBean(bean MyBean)
	RegisterBeanByName(name string, bean MyBean)
	AutoWireBeans() error
	FindBeanDefinitionByName(name string) *MyBeanDefinition
	FindBeansByType(i interface{})
	FindBeanDefinitionsByType(t reflect.Type) []*MyBeanDefinition
	FindBeanByName(name string) MyBean
	PostInitialization()
	StartUp()
}

type BeanFacotry struct {
	beans map[string]*MyBeanDefinition
	BeanContext
}

func (b *BeanFacotry) StartUp() {
	b.AutoWireBeans()
	b.PostInitialization()
}
func (b *BeanFacotry) PostInitialization() {
	for _, bean := range b.beans {
		// 执行 Bean 的初始化接口
		if c, ok := bean.Bean.(MyBeanPostInitialization); ok {
			c.PostInitilization()
		}
	}
}

//注册bean
func (b *BeanFacotry) RegisterBean(bean MyBean) {
	t, v := getBeanType(bean)
	// 是不是应该设置一些自己变量
	b.registerBeanDefinition(t.Elem().Name(), bean, t, v)
}

func (b *BeanFacotry) getAlias(t reflect.Type) []string {
	t1 := t.Elem()
	// 遍历 SpringBean 所有的字段
	for i := 0; i < t1.NumField(); i++ {
		f := t1.Field(i)
		if f.Name == "Describe" {
			if alias, ok := f.Tag.Lookup("alias"); ok {
				if strings.TrimSpace(alias) != "" {
					split := strings.Split(strings.TrimSpace(alias), ",")
					return split
				}
			}
		}
	}
	return []string{}
}

// 注册 Bean 使用指定的 Bean 名称
func (b *BeanFacotry) RegisterNameBean(name string, bean MyBean) {
	t, v := getBeanType(bean)
	b.registerBeanDefinition(name, bean, t, v)
}

// 自动绑定所有的 Bean
func (b *BeanFacotry) AutoWireBeans() error {
	for _, beanDefinition := range b.beans {
		if err := b.wireBeanByDefinition(beanDefinition); err != nil {
			return err
		}
	}
	return nil
}

func (b *BeanFacotry) wiredSliceOrArrayDefinition(f reflect.StructField, v reflect.Value) {
	definitions := b.FindBeanDefinitionsByType(f.Type.Elem())
	result := reflect.Indirect(reflect.New(f.Type))

	for _, def := range definitions {
		result.Set(reflect.Append(result, def.Value))
	}
	v.Set(result)
}

func (b *BeanFacotry) wiredMapDefinition(f reflect.StructField, v reflect.Value) {
	if beanName, ok := f.Tag.Lookup("beanFiledName"); ok {
		key := f.Type.Key()
		definitions := b.FindBeanDefinitionsByType(f.Type.Elem())
		if len(definitions) == 0 {
			panic(fmt.Sprintf("not exists type bean = %v ", f.Type.Name()))
		}
		//反射创建出一个map
		result := reflect.MakeMap(reflect.MapOf(key, f.Type.Elem()))

		for _, definition := range definitions {
			b.wireBeanByDefinition(definition)
			vt := definition.Type.Elem()
			_, exist := vt.FieldByName(beanName)
			if !exist {
				continue
			}
			switch key.Kind() {
			case reflect.Int:
				i := definition.Value.Elem().FieldByName(beanName).Int()
				result.SetMapIndex(reflect.ValueOf(int(i)), definition.Value)
			case reflect.Int32:
				i := definition.Value.Elem().FieldByName(beanName).Int()
				result.SetMapIndex(reflect.ValueOf(int(i)), definition.Value)
			case reflect.Int64:
				i := definition.Value.Elem().FieldByName(beanName).Int()
				result.SetMapIndex(reflect.ValueOf(int(i)), definition.Value)
			case reflect.String:
				i := definition.Value.Elem().FieldByName(beanName).String()
				result.SetMapIndex(reflect.ValueOf(i), definition.Value)
			}
		}
		v.Set(result)
	} else {
		panic(fmt.Sprintf("not exists beanFieldName = %v", beanName))
	}
}

// 绑定 BeanDefinition 指定的 Bean
func (b *BeanFacotry) wireBeanByDefinition(beanDefinition *MyBeanDefinition) error {

	// 确保 SpringBean 还未初始化
	if beanDefinition.Init != Uninitialized {
		return nil
	}

	fmt.Println("wire bean " + beanDefinition.Name)
	defer func() {
		fmt.Println("success wire bean " + beanDefinition.Name)
	}()

	beanDefinition.Init = Initializing
	t := beanDefinition.Type.Elem()
	v := beanDefinition.Value.Elem()
	// 遍历 SpringBean 所有的字段
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// 查找依赖绑定的标签
		if beanName, ok := f.Tag.Lookup("autowire"); ok {
			var definition *MyBeanDefinition
			if len(beanName) > 0 {
				definition = b.FindBeanDefinitionByName(beanName)
			} else {
				//如果 值是map
				switch f.Type.Kind() {
				case reflect.Map:
					b.wiredMapDefinition(f, v.Field(i))
					continue
				case reflect.Slice, reflect.Array:
					b.wiredSliceOrArrayDefinition(f, v.Field(i))
					continue
				default:
				}
				definitions := b.FindBeanDefinitionsByType(f.Type)
				if len(definitions) > 0 {
					definition = definitions[0]
				}
			}

			if definition != nil {
				b.wireBeanByDefinition(definition)
				v.Field(i).Set(definition.Value)
			}

			continue
		}

		// 查找属性绑定的标签
		//if value, ok := f.Tag.Lookup("value"); ok && len(value) > 0 {
		//	if strings.HasPrefix(value, "${") {
		//		str := value[2 : len(value)-1]
		//		ss := strings.Split(str, ":=")
		//
		//		var (
		//			propName  string
		//			propValue interface{}
		//		)
		//
		//		propName = ss[0]
		//		if len(ss) > 1 {
		//			propValue = ss[1]
		//		}
		//
		//		//if prop, ok := b.GetDefaultProperties(propName, ""); ok {
		//		//	propValue = prop
		//		//} else {
		//		//	if len(ss) < 2 {
		//		//		return errors.New("properties " + propName + " not config!")
		//		//	}
		//		//}
		//
		//		vf := v.Field(i)
		//
		//		switch vf.Kind() {
		//		case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		//			u := cast.ToUint64(propValue)
		//			vf.SetUint(u)
		//		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		//			i := cast.ToInt64(propValue)
		//			vf.SetInt(i)
		//		case reflect.String:
		//			s := cast.ToString(propValue)
		//			vf.SetString(s)
		//		case reflect.Bool:
		//			b := cast.ToBool(propValue)
		//			vf.SetBool(b)
		//		default:
		//			return errors.New("unsupported type " + vf.Type().String())
		//		}
		//	}
		//
		//	continue
		//}
	}

	// 执行 Bean 的初始化接口
	if c, ok := beanDefinition.Bean.(MyBeanInitialization); ok {
		c.InitBean(b)
	}

	beanDefinition.Init = Initialized
	return nil
}

// 根据 Bean 名称查找 BeanDefinition
func (b *BeanFacotry) FindBeanDefinitionByName(name string) *MyBeanDefinition {
	return b.beans[name]
}

// 根据 Bean 类型查找 BeanDefinition 数组
func (b *BeanFacotry) FindBeanDefinitionsByType(t reflect.Type) []*MyBeanDefinition {
	result := make([]*MyBeanDefinition, 0)
	for _, beanDefinition := range b.beans {
		if beanDefinition.Type.AssignableTo(t) {
			result = append(result, beanDefinition)
		}
	}
	return result
}

func (b *BeanFacotry) registerBeanDefinition(name string, bean MyBean, t reflect.Type, v reflect.Value) {
	b.beans[name] = &MyBeanDefinition{
		Init:  Uninitialized,
		Name:  name,
		Bean:  bean,
		Type:  t,
		Value: v,
	}
	alias := b.getAlias(t)
	if len(alias) != 0 {
		for _, alia := range alias {
			b.beans[alia] = b.beans[name]
		}
	}

}

func (b *BeanFacotry) FindBeansByType(i interface{}) {
	it := reflect.TypeOf(i)
	et := it.Elem()

	if it.Kind() != reflect.Ptr {
		panic("bean must be pointer")
	}

	v := reflect.New(et).Elem()
	t0 := et.Elem()

	for _, beanDefinition := range b.beans {
		if beanDefinition.Type.AssignableTo(t0) {
			v = reflect.Append(v, beanDefinition.Value)
		}
	}

	reflect.ValueOf(i).Elem().Set(v)
}

// 根据 Bean 名称查找 Beans
func (b *BeanFacotry) FindBeanByName(name string) MyBean {
	if beanDefinition, ok := b.beans[name]; ok {
		return beanDefinition.Bean
	}
	return nil
}

func getBeanType(bean MyBean) (t reflect.Type, v reflect.Value) {
	t = reflect.TypeOf(bean)
	if t.Kind() != reflect.Ptr {
		panic("bean must be pointer")
	}
	v = reflect.ValueOf(bean)
	return
}

func newBeanFactory() BeanContext {
	return &BeanFacotry{
		beans: make(map[string]*MyBeanDefinition),
	}
}
