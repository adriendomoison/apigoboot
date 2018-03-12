/*
	Package rest implement the callback required by the user package
*/
package rest

import (
	"github.com/adriendomoison/apigoboot/errorhandling/apihelper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ResponseDTOUserInfo struct {
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
}

type RequestDTOCheckCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	AuthType string `json:"auth_type"`
}

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

// GetByEmail allows to access the service to retrieve a user info when sending its email (private API)
func (r *rest) CheckCredentials(c *gin.Context) {
	var reqDTO RequestDTOCheckCredentials
	if err := c.BindJSON(&reqDTO); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
	} else {
		if resDTO, err := r.service.CheckCredentials(reqDTO); err != nil {
			c.JSON(apihelper.BuildResponseError(err))
		} else {
			c.JSON(http.StatusOK, resDTO)
		}
	}
}
