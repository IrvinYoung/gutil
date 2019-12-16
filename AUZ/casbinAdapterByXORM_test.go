package AUZ

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
	"xorm.io/xorm"
	"testing"
)

type Permission struct {
	Id         int64
	PolicyType string `xorm:"char(4) index notnull default ''"`
	Role       string `xorm:"char(16) index notnull default ''"`
	Parent     string `xorm:"char(16) index notnull default ''"`
	Action     string `xorm:"char(16) index notnull default ''"`

	Ctime int64 `xorm:"int(11) notnull default 0"`

	Path string `xorm:"varchar(255) index notnull default ''"`
	//... others
}

func (p *Permission) TableName() string {
	return "backend_permission"
}

func (p *Permission) CompareFields() CasbinPolicy {
	var cp CasbinPolicy
	switch p.PolicyType {
	case "p":
		cp = &Permission{
			PolicyType: p.PolicyType,
			Role:       p.Role,
			Action:     p.Action,
			Path:       p.Path,
		}
	case "g":
		cp = &Permission{
			PolicyType: p.PolicyType,
			Role:       p.Role,
			Parent:     p.Parent,
		}
	}
	return cp
}

func (p *Permission) ToPolicy() string {
	var fields []string
	switch p.PolicyType {
	case "p":
		fields = []string{p.PolicyType, p.Role, p.Path, p.Action}
	case "g":
		fields = []string{p.PolicyType, p.Role, p.Parent}
	default:
		fields = make([]string, 0)
	}
	return strings.Join(fields, ",")
}

func (p *Permission) FromPolicy(sec string, ptype string, rule []string) CasbinPolicy {
	if len(rule) < 3 {
		rule = append(rule, make([]string, 3-len(rule))...) //required
	}
	newPermission := &Permission{}
	switch ptype {
	case "p":
		newPermission.PolicyType, newPermission.Role, newPermission.Path, newPermission.Action = ptype, rule[0], rule[1], rule[2]
	case "g":
		newPermission.PolicyType, newPermission.Role, newPermission.Parent = ptype, rule[0], rule[1]
	default:
		return nil
	}
	newPermission.Ctime = time.Now().Unix()
	return newPermission
}

func TestAll(t *testing.T) {
	e, err := testInit(t)
	if err != nil {
		t.Fatal(err)
	}

	testAddPolicy(t, e)
	testSavePolicy(t, e)
	//testRemovePolicy(t, e)
	testRemoveFilteredPolicy(t, e)
}

func testRemoveFilteredPolicy(t *testing.T, e *casbin.Enforcer) {
	ok,err := e.RemoveFilteredPolicy(1, "data1","write")
	if err!=nil{
		t.Error("remove policy failed",err)
	}
	t.Log("remove policy, result=", ok)

	ok,err = e.RemoveFilteredGroupingPolicy(0, "admin")
	if err!=nil{
		t.Error("remove group policy failed",err)
	}
	t.Log("remove group policy, result=", ok)

	return
}

func testRemovePolicy(t *testing.T, e *casbin.Enforcer) {
	ok, err := e.RemovePolicy("alice", "data1", "write")
	if err != nil {
		t.Error(err)
	}
	t.Log("delete policy, result=", ok)

	ok, err = e.RemoveGroupingPolicy("admin", "root")
	if err != nil {
		t.Error(err)
	}
	t.Log("delete group policy, result=", ok)
}

func testSavePolicy(t *testing.T, e *casbin.Enforcer) {
	if err := e.SavePolicy(); err != nil {
		t.Error(err)
	}
}

func testAddPolicy(t *testing.T, e *casbin.Enforcer) {
	_, err := e.AddGroupingPolicy("root", "generic")
	if err != nil {
		t.Error("add group policy failed:", err)
	}
	_, err = e.AddGroupingPolicy("admin", "root")
	if err != nil {
		t.Error("add group policy failed:", err)
	}
	_, err = e.AddPolicy("alice", "data1", "write")
	if err != nil {
		t.Error("add policy failed:", err)
	}
	_, err = e.AddPolicy("bob", "data2", "read")
	if err != nil {
		t.Error("add policy failed:", err)
	}
}

func testInit(t *testing.T) (e *casbin.Enforcer, err error) {
	dbEng, err := xorm.NewEngine("mysql", "root:jiushini@tcp(localhost:3306)/casbin")
	if err != nil {
		return
	}

	cxa, err := NewCasbinXormAdapter(dbEng, &Permission{})
	if err != nil {
		return
	}

	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

	//cxa.AddPolicy("p", "p", []string{"root", "/mgr", "GET"})

	e, err = casbin.NewEnforcer(m, cxa)
	return
}
