package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	age  int
}

func (s Student) RefCallArgs(age int, name *string) error {
	fmt.Println(age, *name)
	return nil
}

func main() {
	testValueOf()
}

func testValueOf() {
	kk := 3
	v := reflect.ValueOf(kk) // flag = 间接
	v.Interface()
	reflect.ValueOf(&kk) // flag = 直接

}

func testStruct() {
	// 得到的结构体 struct { F1 int; F2 string; F3 []int }
	value := makeStruct(0, "", []int{})
	fmt.Println("type: ", value.Elem().Type(), "value: ", value)
	value.Elem().Field(0).SetInt(10)
	value.Elem().Field(1).SetString("hello")
	value.Elem().Field(2).Set(reflect.ValueOf([]int{1, 2, 3, 4}))
	fmt.Println(value)
}

// 根据入参类型返回一个零值结构体,即返回值的类型为PtrTo(typ)
func makeStruct(vals ...interface{}) reflect.Value {
	fields := make([]reflect.StructField, len(vals))
	for i, va := range vals {
		t := reflect.TypeOf(va)
		sf := reflect.StructField{
			Name: fmt.Sprintf("F%d", (i + 1)),
			Type: t,
		}
		fields[i] = sf
	}
	st := reflect.StructOf(fields)
	s := reflect.New(st)
	return s
}

func testFunc() {
	s := Student{
		Name: "Bob",
		age:  18,
	}
	v := reflect.ValueOf(s)
	fmt.Println(v.NumMethod()) //方法个数
	method := v.MethodByName("RefCallArgs")
	if !method.IsZero() {
		mt := method.Type()
		fmt.Println(mt)          //func(int, string) error
		fmt.Println(mt.NumIn())  //入参个数
		fmt.Println(mt.NumOut()) //返回值个数
		name := "Bob"
		// 调用方法
		res := method.Call([]reflect.Value{reflect.ValueOf(123), reflect.ValueOf(&name)})
		// 返回值
		for _, re := range res {
			fmt.Println(re)
		}
	}

}

func testStructSet() {
	type S struct {
		Num  int
		name string
	}

	s := S{
		// 未导出的属性无法Set
		Num:  56,
		name: "hello world",
	}
	value := reflect.ValueOf(&s).Elem()

	fmt.Println(value.Field(0).CanSet())
	value.Field(0).SetInt(123)
	fmt.Println(s)
}

func testSet() {
	var num float64 = 3.14
	pointer := reflect.ValueOf(&num)
	value := pointer.Elem()
	fmt.Println("value can set? ", value.CanSet())
	value.SetFloat(77)
	fmt.Println("new value is ", num)
}

func testTypeElem() {
	type A = [16]int16
	var c <-chan map[A][]byte // chan 存的是map
	tc := reflect.TypeOf(c)
	fmt.Println(tc.Kind())        // chan
	fmt.Println(tc.ChanDir())     // <-chan
	tm := tc.Elem()               // map[[16]int16][]uint8
	fmt.Println(tm.Kind())        // map
	fmt.Println(tm.Key())         // [16]int16
	fmt.Println(tm.Key().Kind())  // array
	fmt.Println(tm.Elem())        // []uint8
	fmt.Println(tm.Elem().Kind()) // slice

}

func testValueElem() {
	x := 56
	y := &x
	var z interface{} = y

	v := reflect.ValueOf(&z)

	vz := v.Elem()
	fmt.Println(vz.Kind()) // interface
	vy := vz.Elem()
	fmt.Println(vy.Kind()) // ptr
	vx := vy.Elem()
	fmt.Println(vx.Kind()) // int
}

func testPtr() {
	a := 56
	aa := 57
	println(reflect.ValueOf(a).Int())          // 如果把a换成&a会报错 reflect: call of reflect.Value.Int on ptr Value
	println(reflect.ValueOf(&aa).Elem().Int()) // Elem()方法会返回所指向的指针的值
}

func testReflect1() {
	student := Student{"Bob", 18}
	creatQuery(student)
}

// 根据结构体生成SQL
func creatQuery(q interface{}) string {
	v := reflect.ValueOf(q)
	t := reflect.TypeOf(q)
	// 是结构体类型
	if v.Kind() == reflect.Struct {
		// 获取结构体名称
		typeName := t.Name()
		query := fmt.Sprintf("insert into %s values(", typeName)
		// 获取属性
		numField := v.NumField()
		for i := 0; i < numField; i++ {
			switch v.Field(i).Kind() {
			case reflect.Int:
				if i == 0 {
					query = fmt.Sprintf("%s%d", query, v.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s,%d", query, v.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query = fmt.Sprintf("%s%s", query, v.Field(i).String())
				} else {
					query = fmt.Sprintf("%s,%s", query, v.Field(i).String())
				}
				//...
			}
		}
		query = query + ")"
		fmt.Println(query)
		return query
	}
	return ""
}
