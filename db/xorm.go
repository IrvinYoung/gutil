package db

import (
	"fmt"
	"xorm.io/xorm"
)
import "errors"

type XormBase interface {
	CreateTable(sess *xorm.Session) error
	Update(sess *xorm.Session, cols ...string) error
	Store(sess *xorm.Session) error
}

func DefaultXormStore(sess *xorm.Session, tableName string, d interface{}) (err error) {
	n, err := sess.Table(tableName).InsertOne(d)
	if err != nil {
		return
	}
	if n != 1 {
		err = fmt.Errorf("%s insert failed,count:%d", tableName, n)
		return
	}
	//do something else
	return
}

func DefaultXormUpdate(sess *xorm.Session, tableName string, id int64, d interface{}, cols ...string) (err error) {
	affected, err := sess.Table(tableName).ID(id).Cols(cols...).Update(d)
	if err != nil {
		return
	}
	if affected > 1 {
		err = fmt.Errorf("%s update row affected not one,but=%d", tableName, affected)
		return
	}
	return
}

func DefaultCreateXormTable(sess *xorm.Session, tableName string, d interface{}) error {
	if tableName == "" {
		return errors.New("table name is invalid")
	}
	return sess.Table(tableName).Sync2(d)
}
