package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RootController struct {
}

func NewRootController() *RootController {
	return &RootController{}
}

func (r *RootController) Root(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
