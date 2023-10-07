package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LikeModel = (*customLikeModel)(nil)

type (
	// LikeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeModel.
	LikeModel interface {
		likeModel
	}

	customLikeModel struct {
		*defaultLikeModel
	}
)

// NewLikeModel returns a model for the database table.
func NewLikeModel(conn sqlx.SqlConn) LikeModel {
	return &customLikeModel{
		defaultLikeModel: newLikeModel(conn),
	}
}
