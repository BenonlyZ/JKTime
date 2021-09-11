package dao

import (
	"database/sql"

	"github.com/pkg/errors"
)

type User struct {
	Id   string
	Name string
	Age  int64
}

func GetUserById(i string) (*User, error) {
	if i != "" {
		return &User{Id: i, Name: "张三", Age: 16}, nil
	}
	//直接使用 fmt.Errorf(" %v", sql.ErrNoRows) 得到的错误不再是sql.ErrNoRows
	//使用 fmt.Errorf(" %w", sql.ErrNoRows) 也能包装，但需要用errors.Is()
	return nil, errors.Wrap(sql.ErrNoRows, "查询字段错误: ")
}
