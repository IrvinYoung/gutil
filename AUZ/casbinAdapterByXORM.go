package auz

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"reflect"
	"xorm.io/xorm"
)

type CasbinXormAdapter struct {
	dbEng  *xorm.Engine
	policy interface{}
}

type CasbinPolicy interface {
	ToPolicy() string                                         //struct -> string (join by comma)
	FromPolicy(sec, ptype string, rule []string) CasbinPolicy //sec,type,rules -> struct
	CompareFields() CasbinPolicy
}

// LoadPolicy loads all policy rules from the storage.
func (cxa *CasbinXormAdapter) LoadPolicy(model model.Model) (err error) {
	// Create a slice to begin with
	myType := reflect.TypeOf(cxa.policy)
	slice := reflect.MakeSlice(reflect.SliceOf(myType), 0, 0)

	// Create a pointer to a slice value and set it to the slice
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)

	if err = cxa.getSession().Find(x.Interface()); err != nil {
		return
	}
	e := x.Elem()
	for i := 0; i < e.Len(); i++ {
		p := e.Index(i).Interface()
		persist.LoadPolicyLine(p.(CasbinPolicy).ToPolicy(), model)
	}
	return
}

// SavePolicy saves all policy rules to the storage.
func (cxa *CasbinXormAdapter) SavePolicy(model model.Model) (err error) {
	//policy_definition
	for ptype, ast := range model["p"] {
		for _, policy := range ast.Policy {
			if err = cxa.mergePolicy("p", ptype, policy); err != nil {
				return
			}
		}
	}
	//role_definition
	for ptype, ast := range model["g"] {
		for _, policy := range ast.Policy {
			if err = cxa.mergePolicy("g", ptype, policy); err != nil {
				return
			}
		}
	}
	return
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (cxa *CasbinXormAdapter) AddPolicy(sec string, ptype string, rule []string) (err error) {
	ca := cxa.policy.(CasbinPolicy)
	n, err := cxa.getSession().InsertOne(ca.FromPolicy(sec, ptype, rule))
	if err != nil {
		return
	}
	if n != 1 {
		err = fmt.Errorf("add policy failed, record affect count is %d", n)
	}
	return
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (cxa *CasbinXormAdapter) RemovePolicy(sec string, ptype string, rule []string) (err error) {
	ca := cxa.policy.(CasbinPolicy)
	p := ca.FromPolicy(sec, ptype, rule).CompareFields()
	has, err := cxa.getSession().Exist(p)
	if err != nil {
		return
	}
	if !has {
		return
	}
	n, err := cxa.getSession().Delete(p)
	if err != nil {
		return
	}
	if n != 1 {
		err = fmt.Errorf("remove policy failed, record affect count is %d", n)
	}
	return
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (cxa *CasbinXormAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (err error) {
	count := fieldIndex + len(fieldValues)
	rule := make([]string, count)
	for i := fieldIndex; i < count; i++ {
		rule[i] = fieldValues[i-fieldIndex]
	}
	err = cxa.RemovePolicy(sec, ptype, rule)
	return
}

// NewCasbinXormAdapter create a new instance
func NewCasbinXormAdapter(dbEng *xorm.Engine, policy interface{}) (cxa *CasbinXormAdapter, err error) {
	if dbEng == nil || policy == nil {
		err = errors.New("params are invalid")
		return
	}

	cxa = &CasbinXormAdapter{
		dbEng:  dbEng,
		policy: policy,
	}

	err = cxa.dbEng.Sync2(policy)
	return
}

func (cxa *CasbinXormAdapter) getSession() *xorm.Session {
	return cxa.dbEng.Table(cxa.policy)
}

func (cxa *CasbinXormAdapter) mergePolicy(sec string, ptype string, rule []string) (err error) {
	ca := cxa.policy.(CasbinPolicy)
	p := ca.FromPolicy(sec, ptype, rule).CompareFields()
	has, err := cxa.getSession().Exist(p)
	if err != nil {
		return
	}
	if has {
		return
	}
	n, err := cxa.getSession().InsertOne(ca.FromPolicy(sec, ptype, rule))
	if err != nil {
		return
	}
	if n != 1 {
		err = fmt.Errorf("merge policy failed, record affect count is %d", n)
	}
	return
}
