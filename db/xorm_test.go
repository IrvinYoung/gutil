package db

import (
	"testing"
	"xorm.io/xorm"
)
import _ "github.com/go-sql-driver/mysql"

const (
	DSN = "root:jiushini@tcp(localhost:3306)/test"
)

type Person struct {
	Name string `xorm:"varchar(8) notnull default ''"`
	Age  int64  `xorm:"int(11) notnull default 0"`
}

func (p *Person)TableName()string{
	return "person"
}

func (p *Person) CreateTable(sess *xorm.Session) (err error) {
	err = DefaultCreateXormTable(sess, p.TableName(), p)
	return
}

func (p *Person) Update(sess *xorm.Session, cols ...string) (err error) {
	err = DefaultXormUpdate(sess, p.TableName(), 1, p, cols...)
	return
}

func (p *Person) Store(sess *xorm.Session) (err error) {
	err = DefaultXormStore(sess, p.TableName(), p)
	return
}

func TestXormCreateTable(t *testing.T) {
	eng, err := xorm.NewEngine("mysql", DSN)
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	sess := eng.NewSession()
	if err = sess.Begin(); err != nil {
		t.Fatal(err)
	}

	var p Person
	if err = p.CreateTable(sess); err != nil {
		sess.Rollback()
		t.Fatal(err)
	}

	if err = sess.Commit(); err != nil {
		t.Fatal(err)
	}
	t.Log("done")
}
