package jsonmodel

import "github.com/gin-gonic/gin"

type Interface interface {
	AppAuthorize(c *gin.Context)
	AppToken(c *gin.Context)
	AppInfo(c *gin.Context)
	AppAuthCode(c *gin.Context)
	AppAuthToken(c *gin.Context)
	AppAuthPassword(c *gin.Context)
	AppAuthClientCredentials(c *gin.Context)
	AppAuthAssertion(c *gin.Context)
	AppAuthRefresh(c *gin.Context)
	AppAuthInfo(c *gin.Context)
}
type ManagerInterface interface {

}