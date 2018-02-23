package model

import (
	"github.com/adriendomoison/go-boot-api/profile/rest/jsonmodel"
	"github.com/adriendomoison/go-boot-api/apicore/helpers/servicehelper"
)

type Interface interface {
	Add(jsonmodel.RequestDTO) (jsonmodel.ResponseDTO, *servicehelper.Error)
	Retrieve(string) (jsonmodel.ResponseDTO, *servicehelper.Error)
	Edit(jsonmodel.RequestDTO) (jsonmodel.ResponseDTO, *servicehelper.Error)
	Remove(string) (*servicehelper.Error)
}