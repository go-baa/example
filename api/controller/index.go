package controller

import (
	"gopkg.in/baa.v1"
)

type index struct{}

// IndexController ...
var IndexController = index{}

// Index list articles
func (index) Index(c *baa.Context) {

}

// Show show article
func (index) Show(c *baa.Context) {

}
