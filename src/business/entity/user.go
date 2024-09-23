package entity

import (
	"github.com/reyhanmichiels/go-pkg/auth"
	"github.com/reyhanmichiels/go-pkg/null"
	"github.com/reyhanmichiels/go-pkg/query"
)

type User struct {
	ID           int64       `db:"id" json:"id"`
	RoleID       int64       `db:"fk_role_id" json:"roleID"`
	Name         string      `db:"name" json:"name"`
	Email        string      `db:"email" json:"email"`
	Password     string      `db:"password" json:"password"`
	RefreshToken null.String `db:"refresh_token" json:"refreshToken" swaggertype:"string"`
	Status       int64       `db:"status" json:"status"`
	Flag         int64       `db:"flag" json:"flag,omitempty"`
	Meta         null.String `db:"meta" json:"meta,omitempty" swaggertype:"string"`
	CreatedAt    null.Time   `db:"created_at" json:"createdAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	CreatedBy    null.String `db:"created_by" json:"createdBy" swaggertype:"string"`
	UpdatedAt    null.Time   `db:"updated_at" json:"updatedAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy    null.String `db:"updated_by" json:"updatedBy" swaggertype:"string"`
	DeletedAt    null.Time   `db:"deleted_at" json:"deletedAt,omitempty" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy    null.String `db:"deleted_by" json:"deletedBy,omitempty" swaggertype:"string"`
}

type UserInputParam struct {
	RoleID          int64       `db:"fk_role_id" json:"-"`
	Name            string      `db:"name" json:"name"`
	Email           string      `db:"email" json:"email"`
	Password        string      `db:"password" json:"password"`
	ConfirmPassword string      `db:"-" json:"confirmPassword"`
	CreatedAt       null.Time   `db:"created_at" json:"-"`
	CreatedBy       null.String `db:"created_by" json:"-"`
}

type UserUpdateParam struct {
	Name         string      `db:"name" json:"name"`
	RefreshToken string      `db:"refresh_token" json:"refreshToken"`
	UpdatedAt    null.Time   `db:"updated_at" json:""`
	UpdatedBy    null.String `db:"updated_by" json:""`
}

type UserParam struct {
	ID           int64  `db:"id" uri:"user_id" param:"id"`
	Email        string `db:"email" param:"email"`
	RefreshToken string `db:"refresh_token" param:"refresh_token"`
	PaginationParam
	QueryOption query.Option
	BypassCache bool
}

type UserLoginParam struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type UserLoginResponse struct {
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenParam struct {
	RefreshToken string `json:"refreshToken"`
}

func (u *User) ConvertToUserAuth() auth.User {
	return auth.User{
		ID:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		RoleID: u.RoleID,
	}
}
