package rdb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return db, nil
}

func CommonToJson[T any](t T) (string, error) {
	s, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func CommonFromJson[T any](t T, s string) error {
	err := json.Unmarshal([]byte(s), &t)
	if err != nil {
		return err
	}
	return nil
}

type JsonWrapper[T any] struct {
	Inner T
}

func NewJsonWrapper[T any](t T) JsonWrapper[T] {
	return JsonWrapper[T]{inner: t}
}

func (j *JsonWrapper[T]) Scan(src interface{}) error {
	if src == nil {
		*j = JsonWrapper[T]{}
		return nil
	}
	var source string
	switch t := src.(type) {
	case string:
		source = t
	case []byte:
		if len(t) == 0 {
			source = ""
		} else {
			source = string(t)
		}
	case nil:
		source = ""
	default:
		return fmt.Errorf("不支持的类型")
	}
	return CommonFromJson(&j.Inner, source)
}

func (j JsonWrapper[T]) Value() (driver.Value, error) {
	return CommonToJson(j.Inner)
}
