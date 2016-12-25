package regtree

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/baa.v1"
)

type storeTreeEntry struct {
	pattern string
	params  int
}

func TestTreeAdd(t *testing.T) {
	tests := []struct {
		id      string
		entries []storeTreeEntry
	}{
		{
			"all static",
			[]storeTreeEntry{
				{"/gopher/bumper.png", 0},
				{"/gopher/bumper192x108.png", 0},
				{"/gopher/doc.png", 0},
				{"/gopher/bumper320x180.png", 0},
				{"/gopher/docpage.png", 0},
				{"/gopher/doc", 0},
			},
		},
		{
			"parametric",
			[]storeTreeEntry{
				{"/users/:id", 1},
				{"/users/:id/profile", 1},
				{"/users/:id/:accnt(\\d+)/address", 2},
				{"/users/:id/age", 1},
				{"/users/:id/:accnt(\\d+)", 2},
			},
		},
		{
			"corner cases",
			[]storeTreeEntry{
				{"/users/:id/test/:name", 2},
				{"/users/abc/:id/:name", 2},
			},
		},
	}

	Convey("tree add", t, func() {
		for _, test := range tests {
			h := NewTree("/", nil)
			Convey(test.id, func() {
				for _, entry := range test.entries {
					node := h.Add(entry.pattern, []baa.HandlerFunc{f})
					if node == nil {
						fmt.Printf("nil node: %#v", entry)
					}
					So(node, ShouldNotBeNil)
					if len(node.params) != entry.params {
						fmt.Printf("error node: %#v\n", node)
					}
					So(len(node.params), ShouldEqual, entry.params)
				}
			})
		}
	})
}

func TestStoreGet(t *testing.T) {
	pairs := []struct {
		pattern string
		handler baa.HandlerFunc
	}{
		{"/gopher/bumper.png", func(c *baa.Context) { c.SetParam("value", "1"); c.JSON(200, c.Params()) }},
		{"/gopher/bumper192x108.png", func(c *baa.Context) { c.SetParam("value", "2"); c.JSON(200, c.Params()) }},
		{"/gopher/doc.png", func(c *baa.Context) { c.SetParam("value", "3"); c.JSON(200, c.Params()) }},
		{"/gopher/bumper320x180.png", func(c *baa.Context) { c.SetParam("value", "4"); c.JSON(200, c.Params()) }},
		{"/gopher/docpage.png", func(c *baa.Context) { c.SetParam("value", "5"); c.JSON(200, c.Params()) }},
		{"/gopher/doc", func(c *baa.Context) { c.SetParam("value", "6"); c.JSON(200, c.Params()) }},
		{"/users/:id", func(c *baa.Context) { c.SetParam("value", "7"); c.JSON(200, c.Params()) }},
		{"/users/:id/profile", func(c *baa.Context) { c.SetParam("value", "8"); c.JSON(200, c.Params()) }},
		{"/users/:id/:account(\\d+)/address", func(c *baa.Context) { c.SetParam("value", "9"); c.JSON(200, c.Params()) }},
		{"/users/:id/age", func(c *baa.Context) { c.SetParam("value", "10"); c.JSON(200, c.Params()) }},
		{"/users/:id/:account(\\d+)", func(c *baa.Context) { c.SetParam("value", "11"); c.JSON(200, c.Params()) }},
		{"/users/:id/test/:name", func(c *baa.Context) { c.SetParam("value", "12"); c.JSON(200, c.Params()) }},
		{"/users/abc/:id/:name", func(c *baa.Context) { c.SetParam("value", "13"); c.JSON(200, c.Params()) }},
		{"/all/*", func(c *baa.Context) { c.SetParam("value", "14"); c.JSON(200, c.Params()) }},
	}
	for _, pair := range pairs {
		b.Get(pair.pattern, pair.handler)
	}

	tests := []struct {
		pattern string
		value   interface{}
		params  map[string]string
	}{
		{"/gopher/bumper.png", "1", nil},
		{"/gopher/bumper192x108.png", "2", nil},
		{"/gopher/doc.png", "3", nil},
		{"/gopher/bumper320x180.png", "4", nil},
		{"/gopher/docpage.png", "5", nil},
		{"/gopher/doc", "6", nil},
		{"/users/abc", "7", map[string]string{"id": "abc"}},
		{"/users/abc/profile", "8", map[string]string{"id": "abc"}},
		{"/users/abc/123/address", "13", map[string]string{"id": "123", "name": "address"}},
		{"/users/abcd/age", "10", map[string]string{"id": "abcd"}},
		{"/users/abc/123", "11", map[string]string{"id": "abc", "account": "123"}},
		{"/users/abc/test/123", "13", map[string]string{"id": "test", "name": "123"}},
		{"/users/abc/xyz/123", "13", map[string]string{"id": "xyz", "name": "123"}},
		{"/g", nil, nil},
		{"/all", nil, nil},
		{"/all/", "14", nil},
		{"/all/abc", "14", map[string]string{"": "abc"}},
		{"/users/abc/xyz", nil, nil},
	}
	Convey("tree get", t, func() {
		for _, test := range tests {
			resp := request("GET", test.pattern)
			if test.value == nil {
				So(resp.Code, ShouldEqual, 404)
				continue
			}
			body := resp.Body.String()
			values := make(map[string]string)
			err := json.Unmarshal([]byte(body), &values)
			if err != nil {
				t.Logf("json error: %v", err)
			}
			So(resp.Code, ShouldEqual, 200)
			So(values["value"], ShouldEqual, test.value)
			if len(test.params) > 0 {
				// for key, val := range test.params {
				// 	So(values[key], ShouldEqual, val)
				// }
			}
		}
	})
}
