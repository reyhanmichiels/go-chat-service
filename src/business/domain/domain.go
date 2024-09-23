package domain

import (
	"github.com/reyhanmichiels/go-pkg/log"
	"github.com/reyhanmichiels/go-pkg/parser"
	"github.com/reyhanmichiels/go-pkg/redis"
	"github.com/reyhanmichiels/go-pkg/sql"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/domain/user"
)

type Domains struct {
	User user.Interface
}

type InitParam struct {
	Log   log.Interface
	Db    sql.Interface
	Redis redis.Interface
	Json  parser.JSONInterface
	// TODO: add audit
}

func Init(param InitParam) *Domains {
	return &Domains{
		User: user.Init(user.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
	}
}
