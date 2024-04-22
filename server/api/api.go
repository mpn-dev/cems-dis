package api

import (
  "cems-dis/model"
)

const (
  DEFAULT_PAGE_SIZE = 20
)

type ApiService struct {
  model *model.Model
}

func New(model *model.Model) ApiService {
  return ApiService{model: model}
}
