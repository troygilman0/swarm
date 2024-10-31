package internal

import (
	"math/rand"
	"reflect"
)

type random struct {
	*rand.Rand
	maxSliceLen int
}

func newRandom(seed int64) random {
	return random{
		Rand:        rand.New(rand.NewSource(seed)),
		maxSliceLen: 100,
	}
}

func (r random) Any(typ reflect.Type) any {
	v := reflect.New(typ)
	r.value(v.Elem())
	return v.Elem().Interface()
}

func (r random) value(v reflect.Value) {
	if r.Intn(4) == 0 {
		return
	}

	kind := v.Kind()
	switch kind {
	case reflect.Ptr:
		r.value(v.Elem())

	case reflect.Struct:
		for i := range v.NumField() {
			field := v.Field(i)
			if !field.IsValid() {
				continue
			}
			if !field.CanSet() {
				continue
			}
			r.value(field)
		}

	case reflect.String:
		buff := make([]byte, r.Intn(r.maxSliceLen))
		r.Read(buff)
		v.SetString(string(buff))

	case reflect.Bool:
		v.SetBool(r.Intn(2) == 0)

	case reflect.Array, reflect.Slice:
		if kind == reflect.Slice {
			v.Grow(r.Intn(r.maxSliceLen))
			v.SetLen(v.Cap())
		}
		for i := range v.Len() {
			r.value(v.Index(i))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(r.Int63())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(r.Uint64())

	case reflect.Float32, reflect.Float64:
		v.SetFloat(r.Float64())

	}
}
