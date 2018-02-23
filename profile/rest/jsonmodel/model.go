package jsonmodel

import "github.com/gin-gonic/gin"

type Interface interface {
	Post(*gin.Context)
	Get(*gin.Context)
	Put(*gin.Context)
	Delete(*gin.Context)
}

type RequestDTO struct {
	PublicId          string `json:"public_id"`
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	ProfilePictureUrl string `json:"profile_picture_url"`
	Birthday          string `json:"birthday" binding:"required"`
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
