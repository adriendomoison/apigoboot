package model

import (
	"github.com/adriendomoison/gobootapi/user/rest/jsonmodel"
	"github.com/adriendomoison/gobootapi/apicore/helpers/servicehelper"
)

type Interface interface {
	GetResourceOwnerId(email string) uint
	Add(jsonmodel.RequestDTO) (jsonmodel.ResponseDTO, *servicehelper.Error)
	Retrieve(string) (jsonmodel.ResponseDTO, *servicehelper.Error)
	EditEmail(jsonmodel.RequestDTOPutEmail) (jsonmodel.ResponseDTO, *servicehelper.Error)
	EditPassword(password jsonmodel.RequestDTOPutPassword) (jsonmodel.ResponseDTO, *servicehelper.Error)
	Remove(string) (*servicehelper.Error)
}