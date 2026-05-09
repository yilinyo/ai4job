// pkg/response/response.go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: 20000, Message: "成功", Data: data})
}

func SuccessPage(c *gin.Context, list any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "成功",
		Data:    PageData{List: list, Total: total, Page: page, PageSize: pageSize},
	})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message, Data: nil})
}
