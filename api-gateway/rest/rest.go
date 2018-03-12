// Package rest implement routes
package rest

import (
	"github.com/adriendomoison/apigoboot/api-gateway/config"
	"github.com/gin-gonic/gin"
)

// AppInfo print basic API info (API version, API name and used port)
func AppInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": "0.0.0 - Coco nut",
		"name":    "apigoboot",
		"port":    config.GPort,
	})
}
