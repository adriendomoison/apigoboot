/*
	user component
*/
package usercomponent

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
func (ms *Component) AttachPublicAPI(group *gin.RouterGroup) {
	group.POST("/users", ms.rest.Post)
	group.GET("/users/:email", ms.rest.ValidateAccessToken, ms.rest.Get)
	group.PUT("/users/:email/email", ms.rest.ValidateAccessToken, ms.rest.PutEmail)
	group.PUT("/users/:email/password", ms.rest.ValidateAccessToken, ms.rest.PutPassword)
	group.DELETE("/users/:email", ms.rest.ValidateAccessToken, ms.rest.Delete)
}


// AttachPrivateAPI add the user micro-service user api with its dependencies
func (ms *Component) AttachPrivateAPI(group *gin.RouterGroup) {
	group.GET("/user/email/:email", ms.rest.ValidateAccessToken, ms.rest.GetByEmail)
	group.GET("/user/id/:userId", ms.rest.ValidateAccessToken, ms.rest.GetById)
}
