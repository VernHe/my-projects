package main

import "testing"

/**

$ go test main.go iface_test.go -bench=. -benchmem
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkDirect-8       1000000000               0.3157 ns/op          0 B/op          0 allocs/op
BenchmarkInterface-8    768032440                1.492 ns/op           0 B/op          0 allocs/op
PASS
ok      command-line-arguments  1.682s

*/

type Shape interface {
	Add(a, b int32) int32
}

type Rectangle struct {
	a int
}

func (r Rectangle) Add(a, b int32) int32 {
	return a + b
}

func BenchmarkDirect(b *testing.B) {
	adder := Rectangle{a: 6379}
	for i := 0; i < b.N; i++ {
		adder.Add(10, 32)
	}
}

func BenchmarkInterface(b *testing.B) {
	adder := Rectangle{a: 6379}
	for i := 0; i < b.N; i++ {
		Shape(adder).Add(10, 32)
	}
}
