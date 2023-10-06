package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindByPhone(ctx context.Context, phone string) (*User, error)
		FindById(ctx context.Context, id int64) (*User, error)
		FindByUsername(ctx context.Context, username string) (*User, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

// FindByPhone finds a user by phone.
func (m *customUserModel) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	err := m.QueryRowNoCacheCtx(ctx, &user, fmt.Sprintf("select %s from %s where `phone` = ? limit 1", userRows, m.table), phone)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindById finds a user by id.
func (m *customUserModel) FindById(ctx context.Context, id int64) (*User, error) {
	user, err := m.FindOne(ctx, id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// FindByUsername finds a user by username.
func (m *customUserModel) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := m.QueryRowNoCacheCtx(ctx, &user, fmt.Sprintf("select %s from %s where `username` = ? limit 1", userRows, m.table), username)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
