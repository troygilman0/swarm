package internal

import (
	"math/rand"
	"reflect"
)

func randomize(a any, r *rand.Rand) any {
	v := reflect.New(reflect.TypeOf(a))
	randomizeValue(v.Elem(), r)
	return v.Interface()
}

func randomizeValue(v reflect.Value, r *rand.Rand) {
	switch v.Kind() {
	case reflect.Ptr:
		randomizePointer(v, r)
	case reflect.Struct:
		randomizeStruct(v, r)
	case reflect.String:
		randomizeString(v, r)
	}
}

func randomizePointer(v reflect.Value, r *rand.Rand) {
	if r.Intn(2) == 0 {
		randomizeValue(v.Elem(), r)
	}
}

func randomizeStruct(v reflect.Value, r *rand.Rand) {
	for i := range v.NumField() {
		field := v.Field(i)
		if !field.IsValid() {
			continue
		}
		if !field.CanSet() {
			continue
		}
		randomizeValue(field, r)
	}
}

func randomizeString(v reflect.Value, r *rand.Rand) {
	buff := make([]byte, r.Intn(100))
	r.Read(buff)
	v.SetString(string(buff))
}
