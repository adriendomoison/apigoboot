/*
	coremodel of the API
*/
package model

import "github.com/gin-gonic/gin"

// MicroServiceInterface is the interface that need to be implemented to plug a micro-service into the system
type MicroServiceInterface interface {
	Attach(group *gin.RouterGroup)
}

