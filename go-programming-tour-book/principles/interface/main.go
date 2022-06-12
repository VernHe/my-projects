package main

type foo interface {
	fooFunc()
}
type foo1 struct{}

func (f1 foo1) fooFunc() {}

// 测试陷阱
type Binary struct {
	uint64
}

type Stringer interface {
	String() string
}

// 通过指针调用方法
func (b *Binary) String() string {
	return "hello world"
}

func main() {
	//var f foo
	//f1 := foo1{}
	//f = foo(f1) // foo1{} escapes to heap
	//fmt.Println(f)
	//f.fooFunc() // 调用方法时，f发生逃逸，因为方法是动态分配的

	// 测试接口陷阱
	//a := Binary{54}
	// 通过指针调用方法 接口中存储的是值（会报错）
	//b := Stringer(a)
	//b.String()

	testParamCopy()
}

func testParamCopy() {
	a := 56
	add(&a)
	println(a)
}

// 传指针是修改原始数据，传值是修改副本
func add(p interface{}) {
	num := p.(*int)
	*num += 1
	println(*num)
}
