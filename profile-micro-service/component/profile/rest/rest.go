// Package rest implement the callback required by the profile package
package rest

import (
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ServiceInterface is the model for the service package of profile
type ServiceInterface interface {
	GetResourceOwnerId(email string) uint
	Add(creation RequestDTOCreation) (ResponseDTO, *servicehelper.Error)
	Retrieve(string) (ResponseDTO, *servicehelper.Error)
	Edit(RequestDTO) (ResponseDTO, *servicehelper.Error)
	Remove(string) *servicehelper.Error
	IsThatTheUserId(string, uint) (bool, *servicehelper.Error)
}

// RequestDTOCreation is the object to map JSON request body of a profile creation request
type RequestDTOCreation struct {
	PublicId          string `json:"profile_id"`
	FirstName         string `json:"first_name" binding:"required,min=2"`
	LastName          string `json:"last_name" binding:"required,min=2"`
	Email             string `json:"email" binding:"required,email"`
	ProfilePictureUrl string `json:"profile_picture_url"`
	Birthday          string `json:"birthday" binding:"required,min=10"`
}

// RequestDTO is the object to map JSON request body
type RequestDTO struct {
	PublicId          string `json:"profile_id" binding:"required,min=16"`
	FirstName         string `json:"first_name" binding:"required,min=2"`
	LastName          string `json:"last_name" binding:"required,min=2"`
	ProfilePictureUrl string `json:"profile_picture_url" binding:"required,url"`
	Birthday          string `json:"birthday" binding:"required,min=10"`
}

// ResponseDTO is the object to map JSON response body
type ResponseDTO struct {
	PublicId          string `json:"profile_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	ProfilePictureUrl string `json:"profile_picture_url"`
	Birthday          string `json:"birthday"`
	OrderAmount       uint   `json:"order_amount"`
}

// Make sure the interface is implemented correctly
var _ profile.RestInterface = (*rest)(nil)

// Implement interface
type rest struct {
	service ServiceInterface
}

// New return a new rest instance
func New(service ServiceInterface) *rest {
	return &rest{service}
}

// Post allows to access the service to create a profile
func (r *rest) Post(c *gin.Context) {
	var reqDTO RequestDTOCreation
	if err := c.BindJSON(&reqDTO); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
	} else {
		if resDTO, err := r.service.Add(reqDTO); err != nil {
			c.JSON(apihelper.BuildResponseError(err))
		} else {
			c.JSON(http.StatusCreated, resDTO)
		}
	}
}

// Get allows to access the service to retrieve a profile when sending its profile public id
func (r *rest) Get(c *gin.Context) {
	if resDTO, err := r.service.Retrieve(c.Param("profileId")); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, resDTO)
	}
}

// Put allows to access the service to update the properties of a profile
func (r *rest) Put(c *gin.Context) {
	var reqDTO RequestDTO
	if err := c.BindJSON(&reqDTO); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
	} else {
		if resDTO, err := r.service.Edit(reqDTO); err != nil {
			c.JSON(apihelper.BuildResponseError(err))
		} else {
			c.JSON(http.StatusOK, resDTO)
		}
	}
}

// Delete allows to access the service to remove a profile from the records
func (r *rest) Delete(c *gin.Context) {
	if err := r.service.Remove(c.Param("profileId")); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "profile has been deleted successfully"})
	}
}
