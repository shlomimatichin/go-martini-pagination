package pagination

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func Test_Defaults(t *testing.T) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/foobar", Service, func(pagi *Pagination) {
		pagi.Append(pagi.Page)
		pagi.Append(pagi.Offset)
		pagi.Append(pagi.PerPage)
		pagi.SetTotal(3)
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(res.Body.String()), &result)
	if err != nil {
		panic(err)
	}
	expect(t, len(result["data"].([]interface{})), 3)
	expectInt(t, result["data"].([]interface{})[0], 0)
	expectInt(t, result["data"].([]interface{})[1], 0)
	expectInt(t, result["data"].([]interface{})[2], 20)
	expectInt(t, result["total"], 3)
}

func Test_PageSpecified(t *testing.T) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/foobar", Service, func(pagi *Pagination) {
		pagi.Append(pagi.Page)
		pagi.Append(pagi.Offset)
		pagi.Append(pagi.PerPage)
		pagi.SetTotal(300)
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar?page=8", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(res.Body.String()), &result)
	if err != nil {
		panic(err)
	}
	expect(t, len(result["data"].([]interface{})), 3)
	expectInt(t, result["data"].([]interface{})[0], 8)
	expectInt(t, result["data"].([]interface{})[1], 8*20)
	expectInt(t, result["data"].([]interface{})[2], 20)
	expectInt(t, result["total"], 300)
}

func Test_EmptyResult(t *testing.T) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/foobar", Service, func(pagi *Pagination) {
		pagi.SetTotal(0)
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(res.Body.String()), &result)
	if err != nil {
		panic(err)
	}
	expect(t, len(result["data"].([]interface{})), 0)
	expectInt(t, result["total"], 0)
}

func Test_ErrorAbortsPagination(t *testing.T) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/foobar", Service, func(pagi *Pagination, r render.Render) {
		r.JSON(500, map[string]interface{}{"Error": "internal error"})
		pagi.SetAbort()
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 500)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(res.Body.String()), &result)
	if err != nil {
		panic(err)
	}
	expectString(t, result["Error"], "internal error")
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func expectInt(t *testing.T, a interface{}, num int) {
	if int(a.(float64)) != num {
		t.Errorf("Expected %d - Got %f", num, a.(float64))
	}
}

func expectString(t *testing.T, a interface{}, expected string) {
	if a.(string) != expected {
		t.Errorf("Expected '%s' - Got '%s'", expected, a.(string))
	}
}
