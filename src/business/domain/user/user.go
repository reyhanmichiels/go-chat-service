package user

import (
	"context"
	"fmt"

	"github.com/reyhanmichiels/go-pkg/errors"
	"github.com/reyhanmichiels/go-pkg/log"
	"github.com/reyhanmichiels/go-pkg/parser"
	"github.com/reyhanmichiels/go-pkg/redis"
	"github.com/reyhanmichiels/go-pkg/sql"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

type Interface interface {
	GetList(ctx context.Context, param entity.UserParam) ([]entity.User, *entity.Pagination, error)
	Get(ctx context.Context, param entity.UserParam) (entity.User, error)
	Create(ctx context.Context, inputParam entity.UserInputParam) (entity.User, error)
	Update(ctx context.Context, updateParam entity.UserUpdateParam, selectParam entity.UserParam) error
}

type user struct {
	db    sql.Interface
	log   log.Interface
	redis redis.Interface
	json  parser.JSONInterface
}

type InitParam struct {
	Db    sql.Interface
	Log   log.Interface
	Redis redis.Interface
	Json  parser.JSONInterface
}

func Init(param InitParam) Interface {
	return &user{
		db:    param.Db,
		log:   param.Log,
		redis: param.Redis,
		json:  param.Json,
	}
}

func (u *user) GetList(ctx context.Context, param entity.UserParam) ([]entity.User, *entity.Pagination, error) {
	if !param.BypassCache {
		user, pg, err := u.getCacheList(ctx, param)
		switch {
		case errors.Is(err, redis.Nil):
			u.log.Error(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			u.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return user, &pg, nil
		}
	}

	user, pg, err := u.getListSQL(ctx, param)
	if err != nil {
		return user, pg, err
	}

	err = u.upsertCacheList(ctx, param, user, *pg, u.redis.GetDefaultTTL(ctx))
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return user, pg, nil
}

func (u *user) Get(ctx context.Context, param entity.UserParam) (entity.User, error) {
	user := entity.User{}

	marshalledParam, err := u.json.Marshal(param)
	if err != nil {
		return user, err
	}

	if !param.BypassCache {
		user, err = u.getCache(ctx, fmt.Sprintf(getUserByKey, string(marshalledParam)))
		switch {
		case errors.Is(err, redis.Nil):
			u.log.Error(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			u.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return user, nil
		}
	}

	user, err = u.getSQL(ctx, param)
	if err != nil {
		return user, err
	}

	err = u.upsertCache(ctx, fmt.Sprintf(getUserByKey, string(marshalledParam)), user, u.redis.GetDefaultTTL(ctx))
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return user, nil
}

func (u *user) Create(ctx context.Context, inputParam entity.UserInputParam) (entity.User, error) {
	user, err := u.createSQL(ctx, inputParam)
	if err != nil {
		return user, err
	}

	err = u.deleteCache(ctx)
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return user, nil
}

func (u *user) Update(ctx context.Context, updateParam entity.UserUpdateParam, selectParam entity.UserParam) error {
	err := u.updateSQL(ctx, updateParam, selectParam)
	if err != nil {
		return err
	}

	err = u.deleteCache(ctx)
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}
