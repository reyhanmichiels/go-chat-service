package user

import (
	"context"
	"fmt"
	"time"

	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

const (
	getUserByKey           = "boilerplate:user:get:%s"
	getUserByQueryKey      = "boilerplate:user:get:q:%s"
	getUserByPaginationKey = "boilerplate:user:get:p:%s"
	deleteUserKeysPattern  = "boilerplate:user*"
)

func (u *user) upsertCache(ctx context.Context, key string, user entity.User, ttl time.Duration) error {
	marshalledUser, err := u.json.Marshal(user)
	if err != nil {
		return errors.NewWithCode(codes.CodeMarshal, err.Error())
	}

	err = u.redis.SetEX(ctx, key, string(marshalledUser), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	return nil
}

func (u *user) getCache(ctx context.Context, key string) (entity.User, error) {
	user := entity.User{}

	marshalledUser, err := u.redis.Get(ctx, key)
	if err != nil {
		return user, err
	}

	err = u.json.Unmarshal([]byte(marshalledUser), &user)
	if err != nil {
		return user, errors.NewWithCode(codes.CodeUnmarshal, err.Error())
	}

	return user, nil
}

func (u *user) upsertCacheList(ctx context.Context, param entity.UserParam, users []entity.User, pg entity.Pagination, ttl time.Duration) error {
	keyValue, err := u.json.Marshal(param)
	if err != nil {
		return errors.NewWithCode(codes.CodeMarshal, err.Error())
	}

	// set user to cache
	marshalledUser, err := u.json.Marshal(users)
	if err != nil {
		return errors.NewWithCode(codes.CodeMarshal, err.Error())
	}

	err = u.redis.SetEX(ctx, fmt.Sprintf(getUserByQueryKey, string(keyValue)), string(marshalledUser), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	// set pagination to cache
	marshalledPagination, err := u.json.Marshal(pg)
	if err != nil {
		return errors.NewWithCode(codes.CodeMarshal, err.Error())
	}

	err = u.redis.SetEX(ctx, fmt.Sprintf(getUserByPaginationKey, string(keyValue)), string(marshalledPagination), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	return nil
}

func (u *user) getCacheList(ctx context.Context, param entity.UserParam) ([]entity.User, entity.Pagination, error) {
	var (
		users = []entity.User{}
		pg    = entity.Pagination{}
	)

	keyValue, err := u.json.Marshal(param)
	if err != nil {
		return users, pg, errors.NewWithCode(codes.CodeMarshal, err.Error())
	}

	// get user from redis
	marshalledUser, err := u.redis.Get(ctx, fmt.Sprintf(getUserByQueryKey, string(keyValue)))
	if err != nil {
		return users, pg, err
	}

	err = u.json.Unmarshal([]byte(marshalledUser), &users)
	if err != nil {
		return users, pg, errors.NewWithCode(codes.CodeUnmarshal, err.Error())
	}

	// get pagination from redis
	marshalledPagination, err := u.redis.Get(ctx, fmt.Sprintf(getUserByPaginationKey, string(keyValue)))
	if err != nil {
		return users, pg, err
	}

	err = u.json.Unmarshal([]byte(marshalledPagination), &pg)
	if err != nil {
		return users, pg, errors.NewWithCode(codes.CodeUnmarshal, err.Error())
	}

	return users, pg, nil
}

func (u *user) deleteCache(ctx context.Context) error {
	err := u.redis.Del(ctx, deleteUserKeysPattern)
	if err != nil {
		return err
	}

	return nil
}
