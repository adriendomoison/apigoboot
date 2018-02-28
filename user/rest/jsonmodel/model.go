package jsonmodel

import "github.com/gin-gonic/gin"

type Interface interface {
	ValidateAccessToken(*gin.Context)
	Post(*gin.Context)
	Get(*gin.Context)
	PutEmail(*gin.Context)
	PutPassword(*gin.Context)
	Delete(*gin.Context)
}

type RequestDTO struct {
	Username  string `json:"username"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type RequestDTOPost struct {
	Username  string `json:"username"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	Birthday  string `json:"birthday" binding:"required,min=10"`
}

type RequestDTOPutEmail struct {
	Email    string `json:"email" binding:"required,email"`
	NewEmail string `json:"new_email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RequestDTOPutPassword struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ResponseDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
