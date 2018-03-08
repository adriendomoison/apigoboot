/*
	core of the API
*/
package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/adriendomoison/apigoboot/apicore/config"
)

// AppInfo print basic API info (API version, API name and used port)
func AppInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": "0.0.0 - Coco nut",
		"name": "apigoboot",
		"port": config.GPort,
	})
}