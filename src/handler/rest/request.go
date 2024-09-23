package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
)

// Bind request body to struct using tag 'json'
func (r *rest) Bind(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindWith(obj, binding.Default(ctx.Request.Method, ctx.ContentType()))
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	return nil
}

// BindQuery bind all query params to struct using tag 'form'
func (r *rest) BindQuery(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindWith(obj, binding.Query)
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	return nil
}

// BindUri bind uri params to struct using tag 'uri'
func (r *rest) BindUri(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindUri(obj)
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	return nil
}

// BindParams bind all params (query and uri params) to struct using tag 'uri' and 'form'
func (r *rest) BindParams(ctx *gin.Context, obj interface{}) error {
	err := r.BindQuery(ctx, obj)
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	err = r.BindUri(ctx, obj)
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	return nil
}
