package contract_test

import (
	"reflect"
	"testing"
)

func TestReflectMap(t *testing.T) {
	v := map[string]interface{}{}
	vt := reflect.TypeOf(v)
	t.Log(vt)
	var newV reflect.Value
	if vt.Kind() != reflect.Map {
		newV = reflect.New(vt)
	} else {
		newV = reflect.MakeMap(vt)
	}
	t.Log(newV)
	if newV.Kind() == reflect.Map {
		newV.SetMapIndex(reflect.ValueOf("a"), reflect.ValueOf("b"))
	}
	// get the poiner of newV
	other := newV.Interface()
	t.Log(&other)
}
