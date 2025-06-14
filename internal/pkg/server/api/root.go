package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Root(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
