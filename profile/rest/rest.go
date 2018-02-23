package rest

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/adriendomoison/gobootapi/profile/rest/jsonmodel"
	"github.com/adriendomoison/gobootapi/profile/service/model"
	"github.com/adriendomoison/gobootapi/apicore/helpers/apihelper"
)

// Make sure the interface is implemented correctly
var _ jsonmodel.Interface = (*rest)(nil)

// Implement interface
type rest struct {
	service model.Interface
}

// New return a new rest instance
func New(service model.Interface) *rest {
	return &rest{service}
}

// Post allows to access the service to create a profile
func (r *rest) Post(c *gin.Context) {
	var reqDTO jsonmodel.RequestDTO
	if err := c.Bind(&reqDTO); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
	} else {
		if resDTO, Err := r.service.Add(reqDTO); Err != nil {
			c.JSON(apihelper.BuildResponseError(Err))
		} else {
			c.JSON(http.StatusCreated, resDTO)
		}
	}
}

// Get allows to access the service to retrieve a profile when sending its public_id
func (r *rest) Get(c *gin.Context) {
	if resDTO, err := r.service.Retrieve(c.Param("profile_public_id")); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, resDTO)
	}
}

// Put allows to access the service to update the properties of a profile
func (r *rest) Put(c *gin.Context) {
	var reqDTO jsonmodel.RequestDTO
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
	if err := r.service.Remove(c.Param("profile_public_id")); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "profile has been deleted successfully"})
	}
}
