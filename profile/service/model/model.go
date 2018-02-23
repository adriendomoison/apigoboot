package model

import (
	"github.com/adriendomoison/gobootapi/profile/rest/jsonmodel"
	"github.com/adriendomoison/gobootapi/apicore/helpers/servicehelper"
)

type Interface interface {
	Add(jsonmodel.RequestDTO) (jsonmodel.ResponseDTO, *servicehelper.Error)
	Retrieve(string) (jsonmodel.ResponseDTO, *servicehelper.Error)
	Edit(jsonmodel.RequestDTO) (jsonmodel.ResponseDTO, *servicehelper.Error)
	Remove(string) (*servicehelper.Error)
}