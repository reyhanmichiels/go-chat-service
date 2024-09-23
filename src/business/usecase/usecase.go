package usecase

import (
	"github.com/reyhanmichiels/go-pkg/auth"
	"github.com/reyhanmichiels/go-pkg/hash"
	"github.com/reyhanmichiels/go-pkg/log"
	"github.com/reyhanmichiels/go-pkg/parser"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/domain"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/usecase/user"
)

type Usecases struct {
	User user.Interface
}

type InitParam struct {
	Dom  *domain.Domains
	Json parser.JSONInterface
	Log  log.Interface
	Hash hash.Interface
	Auth auth.Interface
}

func Init(param InitParam) *Usecases {
	return &Usecases{
		User: user.Init(user.InitParam{UserDomain: param.Dom.User, Auth: param.Auth, Hash: param.Hash}),
	}
}
