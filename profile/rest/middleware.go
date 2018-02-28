package rest

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/adriendomoison/gobootapi/oauthmanagement/repo"
	"github.com/adriendomoison/gobootapi/oauthmanagement/service"
)

func (r *rest) ValidateAccessToken(c *gin.Context) {

	authorizationCode := c.Request.Header.Get("authorization")

	if authorizationCode != "" && len(authorizationCode) > 7 &&
		service.New(repo.New()).GetResourceOwnerId(authorizationCode[7:]) == r.service.GetResourceOwnerId(c.Param("profile_public_id")) {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusForbidden)
	}
}
