package convert

import "strconv"

// 类型转换

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

// Int string转int
func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

// MustInt 转Int
func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

// UInt32 string转int32
func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

// MustUInt32 转 Int32
func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}
