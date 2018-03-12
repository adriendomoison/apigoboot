// Package profile is a component to manage the profile of a user
package profile

import (
	"github.com/gin-gonic/gin"
)

type RestInterface interface {
	ValidateAccessToken(*gin.Context)
	Post(*gin.Context)
	Get(*gin.Context)
	Put(*gin.Context)
	Delete(*gin.Context)
}

// Implement interface
type Component struct {
	rest RestInterface
}

// New return a new micro service instance
func New(rest RestInterface) *Component {
	return &Component{rest}
}

// AttachPublicAPI add the profile micro-service public api with its dependencies
func (ms *Component) AttachPublicAPI(group *gin.RouterGroup) {
	group.GET("/profiles/:profileId", ms.rest.ValidateAccessToken, ms.rest.Get)
	group.PUT("/profiles/:profileId", ms.rest.ValidateAccessToken, ms.rest.Put)
}

// AttachPrivateAPI add the profile micro-service private api with its dependencies
func (ms *Component) AttachPrivateAPI(group *gin.RouterGroup) {
	group.POST("/profiles", ms.rest.Post)
	group.DELETE("/profiles/:profileId", ms.rest.Delete)
}
