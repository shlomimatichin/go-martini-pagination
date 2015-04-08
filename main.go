package pagination

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"math"
	"net/http"
	"strconv"
)

type Pagination struct {
	Page    uint
	Offset  uint
	PerPage uint
	total   uint
	result  []interface{}
}

var DefaultPerPage uint = 20

func (pagination *Pagination) Append(obj interface{}) {
	pagination.result = append(pagination.result, obj)
}

func (pagination *Pagination) SetTotal(total uint) {
	pagination.total = total
}

func (pagination *Pagination) SetAbort() {
	pagination.total = math.MaxUint32 - 1
}

func Service(c martini.Context, req *http.Request, r render.Render) {
	var pagination Pagination
	pagination.Page = 0
	pagination.PerPage = DefaultPerPage
	pagination.total = math.MaxUint32
	pagination.result = make([]interface{}, 0)
	if len(req.URL.Query()["perpage"]) > 0 {
		if len(req.URL.Query()["perpage"]) > 1 {
			panic("More than one perpage parameter attached to get url")
		}
		perpage, err := strconv.ParseUint(req.URL.Query()["perpage"][0], 10, 32)
		if err != nil {
			panic(fmt.Sprintf("Error parsing 'perpage': %s", err))
		}
		pagination.PerPage = uint(perpage)
	}
	if len(req.URL.Query()["page"]) > 0 {
		if len(req.URL.Query()["page"]) > 1 {
			panic("More than one page parameter attached to get url")
		}
		page, err := strconv.ParseUint(req.URL.Query()["page"][0], 10, 32)
		if err != nil {
			panic(fmt.Sprintf("Error parsing 'page': %s", err))
		}
		pagination.Page = uint(page)
		pagination.Offset = uint(page * uint64(pagination.PerPage))
	}
	c.Map(&pagination)
	c.Next()
	if pagination.total == math.MaxUint32-1 {
		return
	}
	if pagination.total == math.MaxUint32 {
		panic("Must set 'SetTotal' on pagination.Pagination")
	}
	resultJSON := map[string]interface{}{
		"data":    pagination.result,
		"total":   pagination.total,
		"page":    pagination.Page,
		"perpage": pagination.PerPage,
	}
	r.JSON(200, resultJSON)
}
