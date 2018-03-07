package rest

import (
	"github.com/adriendomoison/gobootapi/errorhandling/apihelper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetByEmail allows to access the service to retrieve a user info when sending its email (private API)
func (r *rest) GetByEmail(c *gin.Context) {
	if resDTO, err := r.service.RetrieveUserInfoByEmail(c.Param("email")); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, resDTO)
	}
}

// GetByEmail allows to access the service to retrieve a user info when sending its email (private API)
func (r *rest) GetById(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(apihelper.BuildRequestError(err))
	}
	if resDTO, err := r.service.RetrieveUserInfoByUserId(uint(userId)); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, resDTO)
	}
}