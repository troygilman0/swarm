package internal

import (
	"math/rand"
	"reflect"
)

type random struct {
	*rand.Rand
}

func newRandom(seed int64) random {
	return random{
		Rand: rand.New(rand.NewSource(seed)),
	}
}

func (r random) newRandomValue(typ reflect.Type) any {
	v := reflect.New(typ)
	r.randomizeValue(v.Elem())
	return v.Elem().Interface()
}

func (r random) randomizeValue(v reflect.Value) {
	if r.Rand.Intn(4) == 0 {
		return
	}

	switch v.Kind() {
	case reflect.Ptr:
		r.randomizeValue(v.Elem())

	case reflect.Struct:
		for i := range v.NumField() {
			field := v.Field(i)
			if !field.IsValid() {
				continue
			}
			if !field.CanSet() {
				continue
			}
			r.randomizeValue(field)
		}

	case reflect.String:
		buff := make([]byte, r.Rand.Intn(100))
		r.Rand.Read(buff)
		v.SetString(string(buff))

	case reflect.Int:
		v.SetInt(r.Rand.Int63())

	}
}
