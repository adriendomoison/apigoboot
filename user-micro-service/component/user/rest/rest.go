// Package rest implement the callback required by the user package
package rest

import (
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ServiceInterface is the model for the service package of user
type ServiceInterface interface {
	GetResourceOwnerId(email string) uint
	Add(RequestDTO) (ResponseDTO, *servicehelper.Error)
	Retrieve(string) (ResponseDTO, *servicehelper.Error)
	EditEmail(RequestDTOPutEmail) (ResponseDTO, *servicehelper.Error)
	EditPassword(RequestDTOPutPassword) (ResponseDTO, *servicehelper.Error)
	Remove(string) *servicehelper.Error
	CheckCredentials(RequestDTOCheckCredentials) (ResponseDTOUserInfo, *servicehelper.Error)
	AddWithProfile(profile RequestDTOWithProfile) (ResponseDTOWithProfile, *servicehelper.Error)
	RetrieveWithProfile(email string) (ResponseDTOWithProfile, *servicehelper.Error)
	IsThatTheUserId(email string, userIdToCheck uint) (bool, *servicehelper.Error)
	RetrieveUserInfoByEmail(email string) (resDTO ResponseDTOUserInfo, error *servicehelper.Error)
	RetrieveUserInfoByUserId(userId uint) (resDTO ResponseDTOUserInfo, error *servicehelper.Error)
}

// RequestDTO is the object to map JSON request body
type RequestDTO struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RequestDTOWithProfile is the object to map JSON request body when a profile is added to the request
type RequestDTOWithProfile struct {
	Username  string `json:"username"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	Birthday  string `json:"birthday" binding:"required,min=10"`
}

// RequestDTOPutEmail is the object to map JSON request body for requests to edit the email of a user
type RequestDTOPutEmail struct {
	Email    string `json:"email" binding:"required,email"`
	NewEmail string `json:"new_email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RequestDTOPutPassword is the object to map JSON request body for requests to edit the email of a user
type RequestDTOPutPassword struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResponseDTO is the object to map JSON response body
type ResponseDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ResponseDTOWithProfile is the object to map JSON response body when a profile is added to the response
type ResponseDTOWithProfile struct {
	PublicId  string `json:"profile_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Birthday  string `json:"birthday"`
}

// Make sure the interface is implemented correctly
var _ user.RestInterface = (*rest)(nil)

// Implement interface
type rest struct {
	service ServiceInterface
}

// New return a new rest instance
func New(service ServiceInterface) *rest {
	return &rest{service}
}

// Post allows to access the service to create a user
// add create_profile=true query parameter to create the profile alongside the user
func (r *rest) Post(c *gin.Context) {
	var createProfile bool
	if !apihelper.GetBoolQueryParam(c, &createProfile, "create_profile", false) {
		return
	}

	if !createProfile {
		var reqDTO RequestDTO
		if err := c.BindJSON(&reqDTO); err != nil {
			c.JSON(apihelper.BuildRequestError(err))
		} else {
			if resDTO, err := r.service.Add(reqDTO); err != nil {
				c.JSON(apihelper.BuildResponseError(err))
			} else {
				c.JSON(http.StatusCreated, resDTO)
			}
		}
	} else {
		var reqDTO RequestDTOWithProfile
		if err := c.BindJSON(&reqDTO); err != nil {
			c.JSON(apihelper.BuildRequestError(err))
		} else {
			if resDTO, err := r.service.AddWithProfile(reqDTO); err != nil {
				c.JSON(apihelper.BuildResponseError(err))
			} else {
				c.JSON(http.StatusCreated, resDTO)
			}
		}
	}
}

// Get allows to access the service to retrieve a user when sending its email
// add get_profile=true query parameter to get the profile alongside the user
func (r *rest) Get(c *gin.Context) {

	var getProfile bool
	if !apihelper.GetBoolQueryParam(c, &getProfile, "get_profile", false) {
		return
	}

	if resDTO, err := r.service.RetrieveWithProfile(c.Param("email")); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, resDTO)
	}
}

// PutEmail allows to access the service to update the email of a user
func (r *rest) PutEmail(c *gin.Context) {
	var reqDTO RequestDTOPutEmail
	if err := c.BindJSON(&reqDTO); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
	} else {
		if resDTO, err := r.service.EditEmail(reqDTO); err != nil {
			c.JSON(apihelper.BuildResponseError(err))
		} else {
			c.JSON(http.StatusOK, resDTO)
		}
	}
}

// PutPassword allows to access the service to update the password of a user
func (r *rest) PutPassword(c *gin.Context) {
	var reqDTO RequestDTOPutPassword
	if err := c.BindJSON(&reqDTO); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
	} else {
		if reqDTO, err := r.service.EditPassword(reqDTO); err != nil {
			c.JSON(apihelper.BuildResponseError(err))
		} else {
			c.JSON(http.StatusOK, reqDTO)
		}
	}
}

// Delete allows to access the service to remove a user from the records
func (r *rest) Delete(c *gin.Context) {
	if err := r.service.Remove(c.Param("email")); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user has been deleted successfully"})
	}
}
