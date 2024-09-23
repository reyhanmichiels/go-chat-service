package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/reyhanmichiels/go-pkg/appcontext"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
	"github.com/reyhanmichiels/go-pkg/header"
	"github.com/reyhanmichiels/go-pkg/query"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

const (
	infoRequest  string = `httpclient Sent Request: uri=%v method=%v`
	infoResponse string = `httpclient Received Response: uri=%v method=%v resp_code=%v`
)

func (r *rest) CustomRecovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok { // nolint: errorlint
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}
			if brokenPipe {
				// If the connection is dead, we can't write a status to it.
				ctx.Error(err.(error)) // nolint: errcheck
				ctx.Abort()
			} else {
				r.httpRespError(ctx, errors.NewWithCode(codes.CodeInternalServerError, http.StatusText(http.StatusInternalServerError)))
			}
			r.log.Panic(err)
		}
	}()
	ctx.Next()
}

// SetTimeout timeout middleware wraps the request context with a timeout
func (r *rest) SetTimeout(ctx *gin.Context) {
	// wrap the request context with a timeout
	c, cancel := context.WithTimeout(ctx.Request.Context(), 1*time.Second)

	// cancel to clear resources after finished
	defer cancel()

	c = appcontext.SetRequestStartTime(c, time.Now())

	// replace request with context wrapped request
	ctx.Request = ctx.Request.WithContext(c)

	// use a channel to signal completion
	finish := make(chan bool, 1)

	// go to next handler with routine so it can be cancelled
	go func() {
		ctx.Next()
		finish <- true
	}()

	select {
	case <-c.Done():
		ctx.Abort()
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeContextDeadlineExceeded, "Request Timeout"))
	case <-finish:
	}
}

func (r *rest) BodyLogger(ctx *gin.Context) {
	if r.ginConfig.LogRequest {
		r.log.Info(ctx.Request.Context(), fmt.Sprintf(infoRequest, ctx.Request.RequestURI, ctx.Request.Method))
	}

	ctx.Next()
	if r.ginConfig.LogResponse {
		if ctx.Writer.Status() < 300 {
			r.log.Info(ctx.Request.Context(), fmt.Sprintf(infoResponse, ctx.Request.RequestURI, ctx.Request.Method, ctx.Writer.Status()))
		} else {
			r.log.Error(ctx.Request.Context(), fmt.Sprintf(infoResponse, ctx.Request.RequestURI, ctx.Request.Method, ctx.Writer.Status()))
		}
	}
}

func (r *rest) addFieldsToContext(ctx *gin.Context) {
	reqid := ctx.GetHeader(header.KeyRequestID)
	if reqid == "" {
		reqid = uuid.New().String()
	}

	c := ctx.Request.Context()
	c = appcontext.SetRequestId(c, reqid)
	c = appcontext.SetUserAgent(c, ctx.Request.Header.Get(header.KeyUserAgent))
	c = appcontext.SetAcceptLanguage(c, ctx.Request.Header.Get(header.KeyAcceptLanguage))
	c = appcontext.SetServiceVersion(c, r.ginConfig.Meta.Version)
	c = appcontext.SetDeviceType(c, ctx.Request.Header.Get(header.KeyDeviceType))
	c = appcontext.SetCacheControl(c, ctx.Request.Header.Get(header.KeyCacheControl))
	c = appcontext.SetServiceName(c, ctx.Request.Header.Get(header.KeyServiceName))
	ctx.Request = ctx.Request.WithContext(c)
	ctx.Next()
}

func (r *rest) VerifyUser(ctx *gin.Context) {
	userID, err := r.verifyUserToken(ctx)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	err = r.setUserInfo(ctx, userID)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	ctx.Next()
}

func (r *rest) verifyUserToken(ctx *gin.Context) (int64, error) {
	var userID int64

	// get token from context
	token := ctx.Request.Header.Get(header.KeyAuthorization)
	if token == "" {
		return userID, errors.NewWithCode(codes.CodeUnauthorized, "empty token")
	}

	_, err := fmt.Sscanf(token, "Bearer %v", &token)
	if err != nil {
		return userID, errors.NewWithCode(codes.CodeUnauthorized, "invalid token format: %s with err:%v", token, err)
	}

	// verify token
	userID, err = r.auth.ValidateAccessToken(token)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

func (r *rest) setUserInfo(ctx *gin.Context, userID int64) error {
	// set user info to context
	user, err := r.uc.User.Get(ctx.Request.Context(), entity.UserParam{
		ID: userID,
		QueryOption: query.Option{
			IsActive: true,
		},
	})
	if err != nil {
		return err
	}

	c := ctx.Request.Context()
	c = r.auth.SetUserAuthInfo(c, user.ConvertToUserAuth())
	c = appcontext.SetUserId(c, int(userID))
	ctx.Request = ctx.Request.WithContext(c)

	return nil
}
