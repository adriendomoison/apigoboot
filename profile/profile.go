package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/adriendomoison/gobootapi/profile/rest"
	_ "github.com/adriendomoison/gobootapi/oauth/rest"
	"github.com/adriendomoison/gobootapi/profile/repo"
	"github.com/adriendomoison/gobootapi/profile/service"
	"github.com/adriendomoison/gobootapi/profile/rest/jsonmodel"
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

// Attach link the profile micro-service with its dependencies to the system
func (ms *MicroService) Attach(group *gin.RouterGroup) {
	// TODO add middleware to check access token
	group.POST("/profiles", ms.rest.Post)
	group.GET("/profiles/:public_id", ms.rest.Get)
	group.PUT("/profiles/:email", ms.rest.Put)
	group.DELETE("/profiles/:public_id", ms.rest.Delete)
}
