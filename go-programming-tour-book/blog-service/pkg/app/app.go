package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

// Pager 分页信息
type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

// NewResponse 创建一个响应（封装gin.Context）
func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

// ToResponse 响应，将响应内容保存至gin的上下文
func (r Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}

// ToResponseList 分页查询的结果
func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	})
}

// ToErrorResponse 含报错信息的响应
func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}
