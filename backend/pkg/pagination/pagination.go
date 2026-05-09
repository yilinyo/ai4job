// pkg/pagination/pagination.go
package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Page     int
	PageSize int
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func FromContext(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return Params{Page: page, PageSize: pageSize}
}
