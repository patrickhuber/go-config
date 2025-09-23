package config

import (
	"reflect"
)

// Merge merges 'from' into 'to' returning a new value without mutating inputs.
// Semantics (retains original intention):
// - Primitive / non-composite conflicts: keep 'to' value
// - Map conflicts: recursively merge; on leaf conflict keep 'to'
// - Slice/array conflicts: concatenate (from... + to...) if same concrete type
// - Nil handling: if one side is nil, return the other (deep copied for composites)
func Merge(from any, to any) (any, error) {
	if from == nil && to == nil {
		return nil, nil
	}
	if from == nil {
		return to, nil
	}
	if to == nil {
		return deepCopy(from), nil
	}

	vFrom := reflect.ValueOf(from)
	vTo := reflect.ValueOf(to)

	switch vFrom.Kind() {
	case reflect.Map:
		return mergeMapReflect(vFrom, vTo)
	case reflect.Slice, reflect.Array:
		return mergeSliceReflect(vFrom, vTo)
	default:
		// primitives / structs / others -> keep 'to'
		return to, nil
	}
}

func mergeMapReflect(vFrom reflect.Value, vTo reflect.Value) (any, error) {
	if !vTo.IsValid() || (vTo.Kind() == reflect.Interface && vTo.IsNil()) {
		return deepCopyValue(vFrom).Interface(), nil
	}
	if vTo.Kind() != reflect.Map {
		return vTo.Interface(), nil
	}
	if vFrom.Type() != vTo.Type() {
		return vTo.Interface(), nil
	}

	result := reflect.MakeMapWithSize(vTo.Type(), vTo.Len())
	for _, key := range vTo.MapKeys() {
		result.SetMapIndex(key, deepCopyValue(vTo.MapIndex(key)))
	}
	for _, key := range vFrom.MapKeys() {
		fVal := vFrom.MapIndex(key)
		if existing := result.MapIndex(key); existing.IsValid() {
			merged, err := Merge(fVal.Interface(), existing.Interface())
			if err != nil {
				return nil, err
			}
			result.SetMapIndex(key, reflect.ValueOf(merged))
		} else {
			result.SetMapIndex(key, deepCopyValue(fVal))
		}
	}
	return result.Interface(), nil
}

func mergeSliceReflect(vFrom reflect.Value, vTo reflect.Value) (any, error) {
	if !vTo.IsValid() || (vTo.Kind() == reflect.Interface && vTo.IsNil()) {
		return deepCopyValue(vFrom).Interface(), nil
	}
	if !(vTo.Kind() == reflect.Slice || vTo.Kind() == reflect.Array) {
		return vTo.Interface(), nil
	}
	if vFrom.Type() != vTo.Type() {
		return vTo.Interface(), nil
	}
	total := vFrom.Len() + vTo.Len()
	out := reflect.MakeSlice(vTo.Type(), 0, total)
	for i := 0; i < vFrom.Len(); i++ {
		out = reflect.Append(out, deepCopyValue(vFrom.Index(i)))
	}
	for i := 0; i < vTo.Len(); i++ {
		out = reflect.Append(out, vTo.Index(i))
	}
	return out.Interface(), nil
}

func deepCopy(value any) any {
	if value == nil {
		return nil
	}
	return deepCopyValue(reflect.ValueOf(value)).Interface()
}

func deepCopyValue(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}
	if v.Kind() == reflect.Interface && !v.IsNil() {
		elem := v.Elem()
		copied := deepCopyValue(elem)
		if copied.Type() != elem.Type() {
			return copied.Convert(elem.Type())
		}
		return copied
	}
	switch v.Kind() {
	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		out := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			out.Index(i).Set(deepCopyValue(v.Index(i)))
		}
		return out
	case reflect.Array:
		out := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			out.Index(i).Set(deepCopyValue(v.Index(i)))
		}
		return out
	case reflect.Map:
		if v.IsNil() {
			return v
		}
		out := reflect.MakeMapWithSize(v.Type(), v.Len())
		for _, key := range v.MapKeys() {
			out.SetMapIndex(key, deepCopyValue(v.MapIndex(key)))
		}
		return out
	case reflect.Ptr:
		if v.IsNil() {
			return v
		}
		out := reflect.New(v.Elem().Type())
		out.Elem().Set(deepCopyValue(v.Elem()))
		return out
	case reflect.Struct:
		out := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			if out.Field(i).CanSet() {
				out.Field(i).Set(deepCopyValue(v.Field(i)))
			}
		}
		return out
	default:
		return v
	}
}
