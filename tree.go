package options

import (
	"strings"
	"errors"
)

type OptionsTree map[string]interface{}

func NewOptionsTree(data ...map[string]interface{}) OptionsTree {
	oc := make(OptionsTree)
	for _, d := range data {
		for k, v := range d {
			oc.Set(k, v)
		}
	}
	return oc
}

func (oc OptionsTree) Set(key string, value interface{}) OptionsTree {
	if key == "" {
		panic(errors.New("OptionsTree.Set: Key is empty"))
	}
	parts := strings.Split(key, ".")
	end, parts := parts[len(parts)-1], parts[:len(parts)-1]
	o := oc
	for _, p := range parts {
		parent, ok := (o)[p]
		if !ok {
			parent = make(OptionsTree)
			(o)[p] = parent
		}
		o = parent.(OptionsTree)
	}
	(o)[end] = value
	return oc
}

func (oc OptionsTree) Merge(key string, values ...map[string]interface{}) OptionsTree {
	if key == "" {
		panic(errors.New("OptionsTree.Merge: Key is empty"))
	}
	for _, value := range values {
		for k, v := range value {
			if vmap, ok := v.(map[string]interface{}); ok {
				oc.GetMany(k, true).Merge(k, vmap)
			} else {
				(oc)[k] = v
			}
		}
	}
	return oc
}

func (oc OptionsTree) SetMany(key string, values ...map[string]interface{}) OptionsTree {
	oc.Set(key, nil)
	oc.GetMany(key, true).Merge(key, values...)
	return oc
}

func (oc OptionsTree) GetMany(key string, createIfNil ...bool) (value OptionsTree) {
	v, _ := oc.Get(key)

	if v == nil && len(createIfNil) == 1 && createIfNil[0] {
		value = make(OptionsTree)
		oc.Set(key, value)
		return value
	}

	return v.(OptionsTree)
}

func (oc OptionsTree) Get(key string) (value interface{}, ok bool) {
	if key == "" {
		panic(errors.New("OptionsTree.Get: Key is empty"))
	}
	parts := strings.Split(key, ".")
	end, parts := parts[len(parts)-1], parts[:len(parts)-1]
	o := oc
	for _, p := range parts {
		parent, ok := (o)[p]
		if !ok {
			return nil, false
		}
		o = parent.(OptionsTree)
	}
	value, ok = (o)[end]
	return
}

func (oc OptionsTree) GetBool(key string, defaul ... bool) bool {
	v, _ := oc.Get(key)
	if v != nil {
		return v.(bool)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return false
}

func (oc OptionsTree) GetString(key string, defaul ...string) string {
	v, _ := oc.Get(key)
	if v != nil {
		return v.(string)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return ""
}

func (oc OptionsTree) GetInt(key string, defaul ... int) int {
	v, _ := oc.Get(key)
	if v != nil {
		return v.(int)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return 0
}

func (oc OptionsTree) GetSlice(key string, defaul ...[]interface{}) (r []interface{}) {
	v, _ := oc.Get(key)
	if v != nil {
		return v.([]interface{})
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return
}

func (oc OptionsTree) GetStrings(key string, defaul ...[]string) (r []string) {
	v, _ := oc.Get(key)
	if v != nil {
		return v.([]string)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return
}

func (oc OptionsTree) GetInterface(key string, defaul ...interface{}) interface{} {
	if v, ok := oc.Get(key); ok {
		return v
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return nil
}

func (oc OptionsTree) On(key string, f func(ok bool, value interface{}) interface{}) interface{} {
	v, ok := oc.Get(key)
	return f(ok, v)
}
