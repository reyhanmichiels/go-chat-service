package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

// @Summary Register
// @Description Register New User
// @Tags Auth
// @Param data body entity.UserInputParam true "User Data"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.User{}}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 404 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /auth/v1/register [POST]
func (r *rest) RegisterNewUser(ctx *gin.Context) {
	var param entity.UserInputParam

	err := r.Bind(ctx, &param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	authInfo, err := r.uc.User.Register(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, authInfo, nil)
}

// @Summary Sign In
// @Description Sign In With Email and Password
// @Tags Auth
// @Param data body entity.UserLoginParam true "Email And Password"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.UserLoginResponse{}}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 404 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /auth/v1/login [POST]
func (r *rest) SignInWithPassword(ctx *gin.Context) {
	var param entity.UserLoginParam

	err := r.Bind(ctx, &param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	authInfo, err := r.uc.User.SignIn(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, authInfo, nil)
}

// @Summary Refresh Token
// @Description Exchange Refresh Token with new Access Token
// @Tags Auth
// @Param data body entity.RefreshTokenParam true "RefreshToken"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.UserLoginResponse{}}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 404 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /auth/v1/token/refresh [POST]
func (r *rest) RefreshToken(ctx *gin.Context) {
	var param entity.RefreshTokenParam

	err := r.Bind(ctx, &param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	authInfo, err := r.uc.User.RefreshToken(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, authInfo, nil)
}

func (r *rest) DummyLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{})
}
