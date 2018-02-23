/*
	user micro-service root package
*/
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/adriendomoison/gobootapi/user/rest"
	"github.com/adriendomoison/gobootapi/user/repo"
	"github.com/adriendomoison/gobootapi/user/service"
	"github.com/adriendomoison/gobootapi/user/rest/jsonmodel"
	"github.com/adriendomoison/gobootapi/apicore/core/model"
)

// Make sure the interface is implemented correctly
var _ model.MicroServiceInterface = (*MicroService)(nil)

// Implement interface
type MicroService struct {
	rest jsonmodel.Interface
}

// New return a new micro service instance
func New() *MicroService {
	return &MicroService{rest.New(service.New(repo.New()))}
}

// Attach link the user micro-service with its dependencies to the system
func (ms *MicroService) Attach(group *gin.RouterGroup) {
	// TODO add middleware to check access token
	group.POST("/users", ms.rest.Post)
	group.GET("/users/:email", ms.rest.Get)
	group.PUT("/users/:email/email", ms.rest.PutEmail)
	group.PUT("/users/:email/password", ms.rest.PutPassword)
	group.DELETE("/users/:email", ms.rest.Delete)
}
