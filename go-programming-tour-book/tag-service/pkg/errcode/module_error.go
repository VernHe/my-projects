package errcode

var (
	ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")
	ErrorGetTokenFail   = NewError(20010002, "获取Token失败")
)
