package user

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
	"github.com/reyhanmichiels/go-pkg/null"
	"github.com/reyhanmichiels/go-pkg/redis"
	libsql "github.com/reyhanmichiels/go-pkg/sql"
	mock_log "github.com/reyhanmichiels/go-pkg/tests/mock/log"
	mock_parser "github.com/reyhanmichiels/go-pkg/tests/mock/parser"
	mock_redis "github.com/reyhanmichiels/go-pkg/tests/mock/redis"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_user_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	type args struct {
		ctx        context.Context
		inputParam entity.UserInputParam
	}

	mockRedis := mock_redis.NewMockInterface(ctrl)
	mockJson := mock_parser.NewMockJSONInterface(ctrl)

	type mockFields struct {
		redis *mock_redis.MockInterface
		json  *mock_parser.MockJSONInterface
	}

	mockField := mockFields{
		redis: mockRedis,
		json:  mockJson,
	}

	mockTime := time.Now()

	mockArgsInputParam := entity.UserInputParam{
		RoleID:    1,
		Name:      "my name",
		Email:     "test@mail.com",
		CreatedAt: null.TimeFrom(mockTime),
		CreatedBy: null.StringFrom("1"),
	}

	mockResult := entity.User{
		ID:        1,
		RoleID:    mockArgsInputParam.RoleID,
		Name:      mockArgsInputParam.Name,
		Email:     mockArgsInputParam.Email,
		Status:    1,
		CreatedAt: mockArgsInputParam.CreatedAt,
		CreatedBy: mockArgsInputParam.CreatedBy,
	}

	query := regexp.QuoteMeta(`
	INSERT INTO user
		(
		 	fk_role_id,
		 	name,
		 	email,
		 	password,
		 	created_at,
		 	created_by
		)
		VALUES
		(
		 	?,
		 	?,
		 	?,
		 	?,
		 	?,
		 	?
		)
	`)

	tests := []struct {
		name string
		args
		prepSqlMock func() (*sql.DB, error)
		mockFunc    func(mock mockFields, ctx context.Context, param entity.UserInputParam)
		wantErr     bool
		want        entity.User
	}{
		{
			name: "failed begin transaction",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin().WillReturnError(assert.AnError)

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
			},
			wantErr: true,
		},
		{
			name: "failed exec query",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnError(assert.AnError)

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
			},
			wantErr: true,
		},
		{
			name: "duplicate value",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnError(errors.NewWithCode(codes.CodeSQLUniqueConstraint, entity.DuplicateEntryErrMessage))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
			},
			wantErr: true,
		},
		{
			name: "failed get affected rows",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.ResultNoRows)

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
			},
			wantErr: true,
		},
		{
			name: "no user created",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(0))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
			},
			wantErr: true,
		},
		{
			name: "failed get last insert id",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(1))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
			},
			wantErr: true,
		},
		{
			name: "failed to commit",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit().WillReturnError(assert.AnError)

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
			},
			wantErr: true,
		},
		{
			name: "success - but failed del redis",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
				mock.redis.EXPECT().Del(ctx, deleteUserKeysPattern).Return(assert.AnError)
			},
			wantErr: false,
			want:    mockResult,
		},
		{
			name: "success",
			args: args{
				ctx:        context.Background(),
				inputParam: mockArgsInputParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserInputParam) {
				mock.redis.EXPECT().Del(ctx, deleteUserKeysPattern).Return(nil)
			},
			wantErr: false,
			want:    mockResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockField, tt.args.ctx, tt.args.inputParam)
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient := libsql.Init(libsql.Config{
				Driver: "sqlmock",
				Leader: libsql.ConnConfig{
					MockDB: sqlServer,
				},
				Follower: libsql.ConnConfig{
					MockDB: sqlServer,
				},
			}, logger)

			d := Init(InitParam{Db: sqlClient, Log: logger, Redis: mockRedis, Json: mockJson})
			got, err := d.Create(tt.args.ctx, tt.args.inputParam)
			if (err != nil) && !tt.wantErr {
				t.Errorf("User.Create() err %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_user_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockRedis := mock_redis.NewMockInterface(ctrl)
	mockJson := mock_parser.NewMockJSONInterface(ctrl)

	type mockFields struct {
		redis *mock_redis.MockInterface
		json  *mock_parser.MockJSONInterface
	}

	mockField := mockFields{
		redis: mockRedis,
		json:  mockJson,
	}

	type args struct {
		ctx         context.Context
		updateParam entity.UserUpdateParam
		selectParam entity.UserParam
	}

	mockTime := time.Now()

	mockUpdateParam := entity.UserUpdateParam{
		Name:         "my-name",
		RefreshToken: "refresh-token",
		UpdatedAt:    null.TimeFrom(mockTime),
		UpdatedBy:    null.StringFrom("1"),
	}

	mockSelectParam := entity.UserParam{
		ID: 1,
	}

	updateQuery := " SET name=?, refresh_token=?, updated_at=?, updated_by=?"
	query := regexp.QuoteMeta(updateUser + updateQuery)

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		mockFunc    func(mock mockFields, ctx context.Context)
		wantErr     bool
	}{
		{
			name: "failed begin tx",
			args: args{
				ctx:         context.Background(),
				updateParam: mockUpdateParam,
				selectParam: mockSelectParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin().WillReturnError(errors.NewWithCode(codes.CodeSQLTxBegin, "failed to begin tx"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context) {
			},
			wantErr: true,
		},
		{
			name: "failed to exec cause duplicate entry",
			args: args{
				ctx:         context.Background(),
				updateParam: mockUpdateParam,
				selectParam: mockSelectParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnError(errors.NewWithCode(codes.CodeInvalidValue, "Duplicate entry"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context) {
			},
			wantErr: true,
		},
		{
			name: "failed to exec",
			args: args{
				ctx:         context.Background(),
				updateParam: mockUpdateParam,
				selectParam: mockSelectParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnError(errors.NewWithCode(codes.CodeSQLTxExec, "failed to exec"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context) {
			},
			wantErr: true,
		},
		{
			name: "failed to get rows affected",
			args: args{
				ctx:         context.Background(),
				updateParam: mockUpdateParam,
				selectParam: mockSelectParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.ResultNoRows)

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context) {
			},
			wantErr: true,
		},
		{
			name: "no rows affected",
			args: args{
				ctx:         context.Background(),
				updateParam: mockUpdateParam,
				selectParam: mockSelectParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(0))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context) {
			},
			wantErr: true,
		},
		{
			name: "failed to commit",
			args: args{
				ctx:         context.Background(),
				updateParam: mockUpdateParam,
				selectParam: mockSelectParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(1))
				sqlMock.ExpectCommit().WillReturnError(errors.NewWithCode(codes.CodeSQLTxCommit, "failed to commit"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context) {
			},
			wantErr: true,
		},
		{
			name: "success - but failed to delete cache",
			args: args{
				ctx:         context.Background(),
				updateParam: mockUpdateParam,
				selectParam: mockSelectParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(1))
				sqlMock.ExpectCommit()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context) {
				mock.redis.EXPECT().Del(ctx, deleteUserKeysPattern).Return(assert.AnError)
			},
			wantErr: false,
		},
		{
			name: "success",
			args: args{
				ctx:         context.Background(),
				updateParam: mockUpdateParam,
				selectParam: mockSelectParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(1))
				sqlMock.ExpectCommit()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context) {
				mock.redis.EXPECT().Del(ctx, deleteUserKeysPattern).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockField, tt.args.ctx)
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient := libsql.Init(libsql.Config{
				Driver: "sqlmock",
				Leader: libsql.ConnConfig{
					MockDB: sqlServer,
				},
				Follower: libsql.ConnConfig{
					MockDB: sqlServer,
				},
			}, logger)

			d := Init(InitParam{Db: sqlClient, Log: logger, Redis: mockRedis, Json: mockJson})
			err = d.Update(tt.args.ctx, tt.args.updateParam, tt.args.selectParam)
			if (err != nil) && !tt.wantErr {
				t.Errorf("User.Update() err %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_user_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockRedis := mock_redis.NewMockInterface(ctrl)
	mockJson := mock_parser.NewMockJSONInterface(ctrl)

	type mockFields struct {
		redis *mock_redis.MockInterface
		json  *mock_parser.MockJSONInterface
	}

	mockField := mockFields{
		redis: mockRedis,
		json:  mockJson,
	}

	queryExt := `WHERE 1=1 AND id=?;`
	query := regexp.QuoteMeta(readUser + queryExt)

	mockParam := entity.UserParam{
		ID: 1,
	}

	mockMarshal, _ := json.Marshal(mockParam)

	mockTime := time.Now()

	mockResult := entity.User{
		ID:           1,
		RoleID:       1,
		Name:         "my name",
		Email:        "test@mail.com",
		RefreshToken: null.StringFrom("refresh-token"),
		Status:       1,
		CreatedAt:    null.TimeFrom(mockTime),
		CreatedBy:    null.StringFrom("1"),
	}

	marshalledResult, _ := json.Marshal(mockResult)

	type args struct {
		ctx   context.Context
		param entity.UserParam
	}

	expectedColumn := []string{"id", "fk_role_id", "name", "email", "refresh_token", "status", "created_at", "created_by"}
	expectedRowResult := []driver.Value{1, 1, "my name", "test@mail.com", "refresh-token", 1, mockTime, 1}

	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		mockFunc    func(mock mockFields, ctx context.Context, param entity.UserParam)
		want        entity.User
		wantErr     bool
	}{
		{
			name: "failed to marshal",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed get from redis",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return("", redis.Nil)
			},
			wantErr: true,
		},
		{
			name: "failed unmarshal",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return(string(marshalledResult), nil)
				mock.json.EXPECT().Unmarshal(marshalledResult, &entity.User{}).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success get data from redis",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return(string(marshalledResult), nil)
				mock.json.EXPECT().Unmarshal(marshalledResult, &entity.User{}).SetArg(1, mockResult).Return(nil)
			},
			wantErr: false,
			want:    mockResult,
		},
		{
			name: "failed to query",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectQuery(query).WillReturnError(errors.NewWithCode(codes.CodeSQLRead, "failed to query"))

				return sqlServer, err
			},
			wantErr: true,
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return(string(marshalledResult), assert.AnError)
			},
		},
		{
			name: "not found",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectQuery(query).WillReturnError(libsql.ErrNotFound)

				return sqlServer, err
			},
			wantErr: true,
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return(string(marshalledResult), assert.AnError)
			},
		},
		{
			name: "structScan failed",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				row := sqlMock.NewRows([]string{"id", "wrong"})
				sqlMock.ExpectQuery(query).WillReturnRows(row)

				return sqlServer, err
			},
			wantErr: true,
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return(string(marshalledResult), assert.AnError)
			},
		},
		{
			name: "success but failed marshal when upsert cache",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				row := sqlMock.NewRows(expectedColumn)
				row.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(row)

				return sqlServer, err
			},
			wantErr: false,
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return(string(marshalledResult), assert.AnError)
				mock.redis.EXPECT().GetDefaultTTL(context.Background()).Return(time.Minute)
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.json.EXPECT().Marshal(mockResult).Return(nil, assert.AnError)
			},
			want: mockResult,
		},
		{
			name: "success but failed set redis when upsert cache",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				row := sqlMock.NewRows(expectedColumn)
				row.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(row)

				return sqlServer, err
			},
			wantErr: false,
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return(string(marshalledResult), assert.AnError)
				mock.redis.EXPECT().GetDefaultTTL(context.Background()).Return(time.Minute)
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.json.EXPECT().Marshal(mockResult).Return(marshalledResult, nil)
				mock.redis.EXPECT().SetEX(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal)), string(marshalledResult), time.Minute).Return(assert.AnError)
			},
			want: mockResult,
		},
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				row := sqlMock.NewRows(expectedColumn)
				row.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(row)

				return sqlServer, err
			},
			wantErr: false,
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal))).Return(string(marshalledResult), assert.AnError)
				mock.redis.EXPECT().GetDefaultTTL(context.Background()).Return(time.Minute)
				mock.json.EXPECT().Marshal(param).Return(mockMarshal, nil)
				mock.json.EXPECT().Marshal(mockResult).Return(marshalledResult, nil)
				mock.redis.EXPECT().SetEX(ctx, fmt.Sprintf(getUserByKey, string(mockMarshal)), string(marshalledResult), time.Minute).Return(nil)
			},
			want: mockResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockField, tt.args.ctx, tt.args.param)
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient := libsql.Init(libsql.Config{
				Driver: "sqlmock",
				Leader: libsql.ConnConfig{
					MockDB: sqlServer,
				},
				Follower: libsql.ConnConfig{
					MockDB: sqlServer,
				},
			}, logger)

			d := Init(InitParam{Db: sqlClient, Log: logger, Redis: mockRedis, Json: mockJson})
			got, err := d.Get(tt.args.ctx, tt.args.param)
			if (err != nil) && !tt.wantErr {
				t.Errorf("User.Get() err %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_user_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockRedis := mock_redis.NewMockInterface(ctrl)
	mockJson := mock_parser.NewMockJSONInterface(ctrl)

	type mockFields struct {
		redis *mock_redis.MockInterface
		json  *mock_parser.MockJSONInterface
	}

	mockField := mockFields{
		redis: mockRedis,
		json:  mockJson,
	}

	mockParam := entity.UserParam{
		PaginationParam: entity.PaginationParam{
			IncludePagination: true,
		},
	}

	mockMarshalledParam, _ := json.Marshal(mockParam)

	mockTime := time.Now()

	mockResult := []entity.User{
		{
			ID:           1,
			RoleID:       1,
			Name:         "my name",
			Email:        "test@mail.com",
			RefreshToken: null.StringFrom("refresh-token"),
			Status:       1,
			CreatedAt:    null.TimeFrom(mockTime),
			CreatedBy:    null.StringFrom("1"),
		},
	}

	mockMarshalledResult, _ := json.Marshal(mockResult)

	mockPaginationResult := entity.Pagination{
		CurrentPage:     1,
		CurrentElements: 1,
		TotalPages:      1,
		TotalElements:   1,
		SortBy:          []string{},
	}

	mockMarshalledPaginationResult, _ := json.Marshal(mockPaginationResult)

	queryExt := `WHERE 1=1 LIMIT 0, 10;`
	query := regexp.QuoteMeta(readUser + queryExt)
	queryCountExt := " WHERE 1=1"
	queryCount := regexp.QuoteMeta(countUser + queryCountExt)

	expectedColumn := []string{"id", "fk_role_id", "name", "email", "refresh_token", "status", "created_at", "created_by"}
	expectedRowResult := []driver.Value{1, 1, "my name", "test@mail.com", "refresh-token", 1, mockTime, 1}

	type args struct {
		ctx   context.Context
		param entity.UserParam
	}

	tests := []struct {
		name         string
		args         args
		prepSqlMock  func() (*sql.DB, error)
		mockFunc     func(mock mockFields, ctx context.Context, param entity.UserParam)
		want         []entity.User
		wantPaginate *entity.Pagination
		wantErr      bool
	}{
		{
			name: "get cache list - failed marshal params",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectQuery(query).WillReturnError(errors.NewWithCode(codes.CodeSQLRead, "failed to query"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(gomock.Any()).Return(nil, assert.AnError)
			},
			want:         []entity.User{},
			wantPaginate: nil,
			wantErr:      true,
		},
		{
			name: "get cache list - failed get user from cache",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectQuery(query).WillReturnError(errors.NewWithCode(codes.CodeSQLRead, "failed to query"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam))).Return("", redis.Nil)
			},
			want:         []entity.User{},
			wantPaginate: nil,
			wantErr:      true,
		},
		{
			name: "get cache list - failed unmarshal user",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectQuery(query).WillReturnError(errors.NewWithCode(codes.CodeSQLRead, "failed to query"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam))).Return(string(mockMarshalledResult), nil)
				mock.json.EXPECT().Unmarshal(mockMarshalledResult, &[]entity.User{}).Return(assert.AnError)
			},
			want:    []entity.User{},
			wantErr: true,
		},
		{
			name: "get cache list - failed get pagination",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectQuery(query).WillReturnError(errors.NewWithCode(codes.CodeSQLRead, "failed to query"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam))).Return(string(mockMarshalledResult), nil)
				mock.json.EXPECT().Unmarshal(mockMarshalledResult, &[]entity.User{}).SetArg(1, mockResult).Return(nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByPaginationKey, string(mockMarshalledParam))).Return("", assert.AnError)
			},
			want:    []entity.User{},
			wantErr: true,
		},
		{
			name: "get cache list - failed unmarshal pagination",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectQuery(query).WillReturnError(errors.NewWithCode(codes.CodeSQLRead, "failed to query"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam))).Return(string(mockMarshalledResult), nil)
				mock.json.EXPECT().Unmarshal(mockMarshalledResult, &[]entity.User{}).SetArg(1, mockResult).Return(nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByPaginationKey, string(mockMarshalledParam))).Return(string(mockMarshalledPaginationResult), nil)
				mock.json.EXPECT().Unmarshal(mockMarshalledPaginationResult, &entity.Pagination{}).Return(assert.AnError)
			},
			want:    []entity.User{},
			wantErr: true,
		},
		{
			name: "get cache list - success",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam))).Return(string(mockMarshalledResult), nil)
				mock.json.EXPECT().Unmarshal(mockMarshalledResult, &[]entity.User{}).SetArg(1, mockResult).Return(nil)
				mock.redis.EXPECT().Get(ctx, fmt.Sprintf(getUserByPaginationKey, string(mockMarshalledParam))).Return(string(mockMarshalledPaginationResult), nil)
				mock.json.EXPECT().Unmarshal(mockMarshalledPaginationResult, &entity.Pagination{}).SetArg(1, mockPaginationResult).Return(nil)
			},
			wantPaginate: &mockPaginationResult,
			want:         mockResult,
			wantErr:      false,
		},
		{
			name: "failed to query",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				sqlMock.ExpectQuery(query).WillReturnError(errors.NewWithCode(codes.CodeSQLRead, "failed to query"))

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
			},
			wantErr: true,
			want:    []entity.User{},
		},
		{
			name: "structScan failed",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				rows := sqlMock.NewRows([]string{"id"})
				rows.AddRow("wrong")
				sqlMock.ExpectQuery(query).WillReturnRows(rows)

				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
			},
			wantErr: true,
			want:    []entity.User{},
		},
		{
			name: "failed to query count",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				rows := sqlMock.NewRows(expectedColumn)
				rows.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(rows)
				sqlMock.ExpectQuery(queryCount).WillReturnError(errors.NewWithCode(codes.CodeSQLRead, "failed to query count"))
				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
			},
			wantErr: true,
			want:    mockResult,
		},
		{
			name: "success - but failed marshal param",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				rows := sqlMock.NewRows(expectedColumn)
				rows.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(rows)
				rowCount := sqlMock.NewRows([]string{"COUNT(*)"}).AddRow(1)
				sqlMock.ExpectQuery(queryCount).WillReturnRows(rowCount)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
				mock.redis.EXPECT().GetDefaultTTL(ctx).Return(time.Minute)
				mock.json.EXPECT().Marshal(param).Return(nil, assert.AnError)
			},
			wantErr:      false,
			want:         mockResult,
			wantPaginate: &mockPaginationResult,
		},
		{
			name: "success - but failed marshal user",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				rows := sqlMock.NewRows(expectedColumn)
				rows.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(rows)
				rowCount := sqlMock.NewRows([]string{"COUNT(*)"}).AddRow(1)
				sqlMock.ExpectQuery(queryCount).WillReturnRows(rowCount)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
				mock.redis.EXPECT().GetDefaultTTL(ctx).Return(time.Minute)
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.json.EXPECT().Marshal(mockResult).Return(nil, assert.AnError)
			},
			wantErr:      false,
			want:         mockResult,
			wantPaginate: &mockPaginationResult,
		},
		{
			name: "success - but failed set user to redis",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				rows := sqlMock.NewRows(expectedColumn)
				rows.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(rows)
				rowCount := sqlMock.NewRows([]string{"COUNT(*)"}).AddRow(1)
				sqlMock.ExpectQuery(queryCount).WillReturnRows(rowCount)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.json.EXPECT().Marshal(mockResult).Return(mockMarshalledResult, nil)
				mock.redis.EXPECT().GetDefaultTTL(ctx).Return(time.Minute)
				mock.redis.EXPECT().SetEX(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam)), string(mockMarshalledResult), time.Minute).Return(assert.AnError)
			},
			wantErr:      false,
			want:         mockResult,
			wantPaginate: &mockPaginationResult,
		},
		{
			name: "success - but failed marshall pagination",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				rows := sqlMock.NewRows(expectedColumn)
				rows.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(rows)
				rowCount := sqlMock.NewRows([]string{"COUNT(*)"}).AddRow(1)
				sqlMock.ExpectQuery(queryCount).WillReturnRows(rowCount)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.json.EXPECT().Marshal(mockResult).Return(mockMarshalledResult, nil)
				mock.redis.EXPECT().GetDefaultTTL(ctx).Return(time.Minute)
				mock.redis.EXPECT().SetEX(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam)), string(mockMarshalledResult), time.Minute).Return(nil)
				mock.json.EXPECT().Marshal(mockPaginationResult).Return(nil, assert.AnError)
			},
			wantErr:      false,
			want:         mockResult,
			wantPaginate: &mockPaginationResult,
		},
		{
			name: "success - but failed set pagination to redis",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				rows := sqlMock.NewRows(expectedColumn)
				rows.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(rows)
				rowCount := sqlMock.NewRows([]string{"COUNT(*)"}).AddRow(1)
				sqlMock.ExpectQuery(queryCount).WillReturnRows(rowCount)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.json.EXPECT().Marshal(mockResult).Return(mockMarshalledResult, nil)
				mock.redis.EXPECT().GetDefaultTTL(ctx).Return(time.Minute)
				mock.redis.EXPECT().SetEX(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam)), string(mockMarshalledResult), time.Minute).Return(nil)
				mock.json.EXPECT().Marshal(mockPaginationResult).Return(mockMarshalledPaginationResult, nil)
				mock.redis.EXPECT().SetEX(ctx, fmt.Sprintf(getUserByPaginationKey, string(mockMarshalledParam)), string(mockMarshalledPaginationResult), time.Minute).Return(assert.AnError)
			},
			wantErr:      false,
			want:         mockResult,
			wantPaginate: &mockPaginationResult,
		},
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()

				rows := sqlMock.NewRows(expectedColumn)
				rows.AddRow(expectedRowResult...)
				sqlMock.ExpectQuery(query).WillReturnRows(rows)
				rowCount := sqlMock.NewRows([]string{"COUNT(*)"}).AddRow(1)
				sqlMock.ExpectQuery(queryCount).WillReturnRows(rowCount)
				return sqlServer, err
			},
			mockFunc: func(mock mockFields, ctx context.Context, param entity.UserParam) {
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, assert.AnError)
				mock.json.EXPECT().Marshal(param).Return(mockMarshalledParam, nil)
				mock.json.EXPECT().Marshal(mockResult).Return(mockMarshalledResult, nil)
				mock.redis.EXPECT().GetDefaultTTL(ctx).Return(time.Minute)
				mock.redis.EXPECT().SetEX(ctx, fmt.Sprintf(getUserByQueryKey, string(mockMarshalledParam)), string(mockMarshalledResult), time.Minute).Return(nil)
				mock.json.EXPECT().Marshal(mockPaginationResult).Return(mockMarshalledPaginationResult, nil)
				mock.redis.EXPECT().SetEX(ctx, fmt.Sprintf(getUserByPaginationKey, string(mockMarshalledParam)), string(mockMarshalledPaginationResult), time.Minute).Return(nil)
			},
			wantErr:      false,
			want:         mockResult,
			wantPaginate: &mockPaginationResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockField, tt.args.ctx, tt.args.param)
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient := libsql.Init(libsql.Config{
				Driver: "sqlmock",
				Leader: libsql.ConnConfig{
					MockDB: sqlServer,
				},
				Follower: libsql.ConnConfig{
					MockDB: sqlServer,
				},
			}, logger)

			d := Init(InitParam{Db: sqlClient, Log: logger, Redis: mockRedis, Json: mockJson})
			got, pg, err := d.GetList(tt.args.ctx, tt.args.param)
			if (err != nil) && !tt.wantErr {
				t.Errorf("User.GetList() err %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantPaginate, pg)
		})
	}
}
