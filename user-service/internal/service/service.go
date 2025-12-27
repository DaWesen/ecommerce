package service

import (
	"context"
	"ecommerce/user-service/kitex_gen/api"
)

type UserService interface {
	Register(ctx context.Context, req *api.RegisterReq) (*api.RegisterResp, error)
}
