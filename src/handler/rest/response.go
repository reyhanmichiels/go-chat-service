package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/appcontext"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
	"github.com/reyhanmichiels/go-pkg/header"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

func (r *rest) httpRespSuccess(ctx *gin.Context, code codes.Code, data interface{}, p *entity.Pagination) {
	c := ctx.Request.Context()

	displayMessage := codes.Compile(code, appcontext.GetAcceptLanguage(c))

	// build meta data
	meta := entity.Meta{
		Path:       r.ginConfig.Meta.Host + ctx.Request.URL.String(),
		StatusCode: displayMessage.StatusCode,
		Status:     http.StatusText(displayMessage.StatusCode),
		Message:    fmt.Sprintf("%s %s [%d] %s", ctx.Request.Method, ctx.Request.URL.RequestURI(), displayMessage.StatusCode, http.StatusText(displayMessage.StatusCode)),
		Timestamp:  time.Now().Format(time.RFC3339),
		RequestID:  appcontext.GetRequestId(c),
	}

	reqStartTime := appcontext.GetRequestStartTime(c)
	if !time.Time.IsZero(reqStartTime) {
		meta.TimeElapsed = fmt.Sprintf("%dms", int64(time.Since(reqStartTime)/time.Millisecond))
	}
	// build response structure
	resp := &entity.HTTPResp{
		Message: entity.HTTPMessage{
			Title: displayMessage.Title,
			Body:  displayMessage.Body,
		},
		Meta:       meta,
		Data:       data,
		Pagination: p,
	}

	raw, err := r.json.Marshal(&resp)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeInternalServerError, err.Error()))
		return
	}

	c = appcontext.SetAppResponseCode(c, code)

	c = appcontext.SetResponseHttpCode(c, displayMessage.StatusCode)

	ctx.Request = ctx.Request.WithContext(c)

	ctx.Header(header.KeyRequestID, appcontext.GetRequestId(c))

	ctx.Data(displayMessage.StatusCode, header.ContentTypeJSON, raw)
}

func (r *rest) httpRespError(ctx *gin.Context, err error) {
	c := ctx.Request.Context()

	// handle if request time out
	if errors.Is(c.Err(), context.DeadlineExceeded) {
		err = errors.NewWithCode(codes.CodeContextDeadlineExceeded, "Context Deadline Exceeded")
	}

	httpStatus, displayError := errors.Compile(err, appcontext.GetAcceptLanguage(c))

	statusStr := http.StatusText(httpStatus)

	reqStartTime := appcontext.GetRequestStartTime(c)
	timeElapsed := ""
	if !time.Time.IsZero(reqStartTime) {
		timeElapsed = fmt.Sprintf("%dms", int64(time.Since(reqStartTime)/time.Millisecond))
	}

	errResp := &entity.HTTPResp{
		Message: entity.HTTPMessage{
			Title: displayError.Title,
			Body:  displayError.Body,
		},
		Meta: entity.Meta{
			Path:       r.ginConfig.Meta.Host + ctx.Request.URL.String(),
			StatusCode: httpStatus,
			Status:     statusStr,
			Message:    fmt.Sprintf("%s %s [%d] %s", ctx.Request.Method, ctx.Request.URL.RequestURI(), httpStatus, statusStr),
			Error: &entity.MetaError{
				Code:    int(displayError.Code),
				Message: err.Error(),
			},
			Timestamp:   time.Now().Format(time.RFC3339),
			TimeElapsed: timeElapsed,
			RequestID:   appcontext.GetRequestId(c),
		},
	}

	r.log.Error(c, err)

	c = appcontext.SetAppResponseCode(c, displayError.Code)

	c = appcontext.SetAppErrorMessage(c, fmt.Sprintf("%s - %s", displayError.Title, displayError.Body))

	c = appcontext.SetResponseHttpCode(c, httpStatus)

	ctx.Request = ctx.Request.WithContext(c)

	ctx.Header(header.KeyRequestID, appcontext.GetRequestId(c))
	ctx.AbortWithStatusJSON(httpStatus, errResp)
}
