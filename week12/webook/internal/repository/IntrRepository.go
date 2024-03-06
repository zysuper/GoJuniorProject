package repository

import (
	intrv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/intr/v1"
)

type IntrRepository interface {
	intrv1.InteractiveServiceClient
}

type intrRepository struct {
	intrv1.InteractiveServiceClient
}

func NewIntrRepository(interactiveClient intrv1.InteractiveServiceClient) IntrRepository {
	return &intrRepository{InteractiveServiceClient: interactiveClient}
}
