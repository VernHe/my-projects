package entity

// Data 解析接口返回解析结果
type Data struct {
	Cover string
	Text  string
	Video string
}

// Result 解析接口响应的结果
type Result struct {
	Code int
	Desc string
	Data *Data
	Succ bool
}
