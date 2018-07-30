package options

import (
	"errors"

	"github.com/moisespsena/go-error-wrap"
)

var EmptyKey = errors.New("Key is empty")

type Options map[string]interface{}

func NewOptions(data ...map[string]interface{}) Options {
	oc := make(Options)
	for _, d := range data {
		for k, v := range d {
			oc.Set(k, v)
		}
	}
	return oc
}

func (oc Options) Set(key string, value interface{}) Options {
	if key == "" {
		panic(errwrap.Wrap(EmptyKey, "Options.Set"))
	}
	(oc)[key] = value
	return oc
}

func (oc Options) Merge(key string, values ...map[string]interface{}) Options {
	if key == "" {
		panic(errwrap.Wrap(EmptyKey, "Options.Merge"))
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

func (oc Options) SetMany(key string, values ...map[string]interface{}) Options {
	oc.Set(key, nil)
	oc.GetMany(key, true).Merge(key, values...)
	return oc
}

func (oc Options) HasMany(key ...string) (ok bool) {
	for i, k := range key {
		if k == "" {
			panic(errwrap.Wrap(EmptyKey, "Options.HasMany Key[%v]", i))
		}
		_, ok = (oc)[k]
		if !ok {
			return false
		}
	}
	return true
}

func (oc Options) Has(key string) (ok bool) {
	if key == "" {
		panic(errwrap.Wrap(EmptyKey, "Options.Has"))
	}
	_, ok = (oc)[key]
	return ok
}

func (oc Options) GetMany(key string, createIfNil ...bool) (value Options) {
	v, _ := oc.Get(key)

	if v == nil && len(createIfNil) == 1 && createIfNil[0] {
		value = make(Options)
		oc.Set(key, value)
		return value
	}

	return v.(Options)
}

func (oc Options) Get(key string) (value interface{}, ok bool) {
	if key == "" {
		panic(errwrap.Wrap(EmptyKey, "Options.Get"))
	}
	value, ok = (oc)[key]
	return
}

func (oc Options) GetBool(key string, defaul ...bool) bool {
	v, _ := oc.Get(key)
	if v != nil {
		return v.(bool)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return false
}

func (oc Options) GetString(key string, defaul ...string) string {
	v, _ := oc.Get(key)
	if v != nil {
		return v.(string)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return ""
}

func (oc Options) GetInt(key string, defaul ...int) int {
	v, _ := oc.Get(key)
	if v != nil {
		return v.(int)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return 0
}

func (oc Options) GetSlice(key string, defaul ...[]interface{}) (r []interface{}) {
	v, _ := oc.Get(key)
	if v != nil {
		return v.([]interface{})
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return
}

func (oc Options) GetStrings(key string, defaul ...[]string) (r []string) {
	v, _ := oc.Get(key)
	if v != nil {
		return v.([]string)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return
}

func (oc Options) GetInterface(key string, defaul ...interface{}) interface{} {
	if v, ok := oc.Get(key); ok {
		return v
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return nil
}

func (oc Options) On(key string, f func(ok bool, value interface{}) interface{}) interface{} {
	v, ok := oc.Get(key)
	return f(ok, v)
}
