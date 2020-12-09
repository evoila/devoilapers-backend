package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


// ShowHello godoc
// @Summary Show a hello message
// @Description Get hello message
// @Produce  plain
// @Success 200 {string} string
// @Router / [get]
func ShowHello(c *gin.Context) {
	c.String(http.StatusOK, "Hello!")
}
