package main

import (
	"encoding/json"
	"github.com/allegro/bigcache"
	"github.com/go-programming-tour-book/cache/fast"
	"strconv"
	"time"
)

type Value struct {
	A string
	B int
	C time.Time
	D []byte
	E float32
	F *string
	T T
}

type T struct {
	H int
	I int
	K int
	L int
	M int
	N int
}

func main() {

	test3()
}

func test1() {
	m := make(map[int]*Value, 10000000)
	for i := 0; i < 10000000; i++ {
		m[i] = &Value{}
	}
	for i := 0; i < 10000000; i++ {
		delete(m, i)
		m[1000000+i] = &Value{}
		time.Sleep(5 * time.Millisecond)
	}
}

func test2() {
	// 1.5以后，GC会忽略map中K和V都是基本数据类型的
	/*
		$ GODEBUG=gctrace=1 ./main1
		gc 1 @0.026s 14%: 0.50+35+0 ms clock, 4.0+0.11/71/140+0 ms cpu, 306->313->313 MB, 307 MB goal, 8 P
		gc 2 @0.431s 9%: 0+173+0 ms clock, 0+48/346/855+0 ms cpu, 611->659->654 MB, 626 MB goal, 8 P
		gc 3 @1.488s 11%: 0+678+0 ms clock, 0+140/1356/3358+0 ms cpu, 1243->1415->1403 MB, 1308 MB goal, 8 P
		GC forced
		gc 4 @122.172s 0%: 0+931+0 ms clock, 0+0/1842/5296+0 ms cpu, 1686->1686->1676 MB, 2807 MB goal, 8 P
		GC forced
		gc 5 @243.106s 0%: 0+809+0 ms clock, 0+0/1615/4772+0 ms cpu, 1679->1679->1673 MB, 3353 MB goal, 8 P
		GC forced
		gc 6 @363.926s 0%: 0+862+0 ms clock, 0+0/1712/5057+0 ms cpu, 1676->1676->1670 MB, 3347 MB goal, 8 P
	*/
	m := make(map[int]int, 10000000)
	for i := 0; i < 10000000; i++ {
		m[i] = i
	}
	for i := 0; i < 10000000; i++ {
		delete(m, i)
		m[1000000+i] = i
		time.Sleep(5 * time.Millisecond)
	}
}

// fastCache
func test3() {
	/**
	$ GODEBUG=gctrace=1 ./gc
	gc 1 @0.009s 2%: 0+1.8+0 ms clock, 0+0/2.0/3.0+0 ms cpu, 4->4->3 MB, 5 MB goal, 8 P
	gc 2 @0.017s 4%: 0+3.5+0 ms clock, 0+0/6.1/6.1+0 ms cpu, 7->7->7 MB, 8 MB goal, 8 P
	gc 3 @0.032s 9%: 0.50+8.3+0 ms clock, 4.0+1.6/16/8.3+0 ms cpu, 14->16->15 MB, 15 MB goal, 8 P
	gc 4 @0.065s 9%: 0+13+0.53 ms clock, 0+1.9/25/61+4.2 ms cpu, 29->32->30 MB, 31 MB goal, 8 P
	gc 5 @0.127s 9%: 0+25+0 ms clock, 0+0/49/124+0 ms cpu, 56->62->58 MB, 60 MB goal, 8 P
	gc 6 @0.250s 8%: 0+46+0 ms clock, 0+0/93/211+0 ms cpu, 106->120->112 MB, 116 MB goal, 8 P
	gc 7 @0.480s 8%: 0+96+0 ms clock, 0+0.99/191/470+0 ms cpu, 203->225->208 MB, 224 MB goal, 8 P
	gc 8 @0.930s 9%: 0+259+0 ms clock, 0+0.99/517/1261+0 ms cpu, 380->427->394 MB, 416 MB goal, 8 P
	gc 9 @1.909s 10%: 0+624+0 ms clock, 0+1.9/1241/2995+0 ms cpu, 716->819->752 MB, 788 MB goal, 8 P
	gc 10 @3.968s 11%: 0+1300+0 ms clock, 0+2.9/2595/5975+0 ms cpu, 1354->1582->1450 MB, 1505 MB goal, 8 P
	gc 11 @8.203s 11%: 0+2500+0 ms clock, 0+6.9/4996/12863+0 ms cpu, 2572->2923->2660 MB, 2900 MB goal, 8 P

	*/
	cache := fast.NewFastCache(0, 1024, nil)
	for i := 0; i < 10000000; i++ {
		cache.Set(strconv.Itoa(i), &Value{})
	}
	for i := 0; i < 10000000; i++ {
		cache.Del(strconv.Itoa(i))
		cache.Set(strconv.Itoa(1000000+i), &Value{})
		time.Sleep(5 * time.Millisecond)
	}
}

// bigCache
func test4() {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	for i := 0; i < 10000000; i++ {
		bytes, _ := json.Marshal(&Value{})
		cache.Set(strconv.Itoa(i), bytes)
	}
	for i := 0; i < 10000000; i++ {
		cache.Delete(strconv.Itoa(i))
		bytes, _ := json.Marshal(&Value{})
		cache.Set(strconv.Itoa(1000000+i), bytes)
		time.Sleep(5 * time.Millisecond)
	}
}
