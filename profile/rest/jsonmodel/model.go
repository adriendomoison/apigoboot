package jsonmodel

import "github.com/gin-gonic/gin"

type Interface interface {
	ValidateAccessToken(*gin.Context)
	Get(*gin.Context)
	Put(*gin.Context)
}

type RequestDTO struct {
	PublicId          string `json:"public_id"`
	FirstName         string `json:"first_name" binding:"required,min=2"`
	LastName          string `json:"last_name" binding:"required,min=2"`
	Email             string `json:"email" binding:"required,email"`
	ProfilePictureUrl string `json:"profile_picture_url"`
	Birthday          string `json:"birthday" binding:"required,min=10"`
	OrderAmount       uint   `json:"order_amount"`
}

type ResponseDTO struct {
	PublicId          string `json:"public_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	ProfilePictureUrl string `json:"profile_picture_url"`
	Birthday          string `json:"birthday"`
	OrderAmount       uint   `json:"order_amount"`
}
