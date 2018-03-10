/*
	user component
*/
package user

import (
	"github.com/gin-gonic/gin"
)

type RestInterface interface {
	Post(c *gin.Context)
	Get(c *gin.Context)
	PutEmail(c *gin.Context)
	PutPassword(c *gin.Context)
	Delete(c *gin.Context)
	GetByEmail(c *gin.Context)
	GetById(c *gin.Context)
	CheckCredentials(c *gin.Context)
	ValidateAccessToken(c *gin.Context)
}

// Implement interface
type Component struct {
	rest RestInterface
}

// New return a new micro service instance
func New(rest RestInterface) *Component {
	return &Component{rest}
}

// AttachPublicAPI add the user micro-service public api with its dependencies
func (component *Component) AttachPublicAPI(group *gin.RouterGroup) {
	group.POST("/users", component.rest.Post)
	group.GET("/users/:email", component.rest.ValidateAccessToken, component.rest.Get)
	group.PUT("/users/:email/email", component.rest.ValidateAccessToken, component.rest.PutEmail)
	group.PUT("/users/:email/password", component.rest.ValidateAccessToken, component.rest.PutPassword)
	group.DELETE("/users/:email", component.rest.ValidateAccessToken, component.rest.Delete)
}

// AttachPrivateAPI add the user micro-service user api with its dependencies
func (component *Component) AttachPrivateAPI(group *gin.RouterGroup) {
	group.GET("/user/email/:email", component.rest.GetByEmail)
	group.GET("/user/id/:userId", component.rest.GetById)
	group.POST("/user/check-credentials", component.rest.CheckCredentials)
}
