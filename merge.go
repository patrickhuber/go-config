package config

import (
	"maps"
)

func Merge(from any, to any) (any, error) {
	if from == nil && to == nil {
		return nil, nil
	}
	if from == nil {
		return to, nil
	}
	if to == nil {
		return from, nil
	}
	switch concrete := from.(type) {
	case []any:
		return to, nil
	case map[string]any:
		return mergeMap(concrete, to)
	case string:
		return to, nil
	case float64:
		return to, nil
	case bool:
		return to, nil
	}
	return from, nil
}

func mergeMap(fromMap map[string]any, to any) (any, error) {
	// if to is nil, return fromMap cloned
	if to == nil {
		return maps.Clone(fromMap), nil
	}

	// if from is nil, return toMap cloned
	if fromMap == nil {
		return to, nil
	}

	// check if to is a map, if not return to
	toMap, ok := to.(map[string]any)
	if !ok {
		return to, nil
	}

	// clone the from map to avoid mutating
	m := maps.Clone(fromMap)

	// merge keys in the from map with keys in the to map
	for k, vFrom := range m {
		// if the key exists in both maps, merge the values
		if vTo, ok := toMap[k]; ok {
			var err error
			vFrom, err = Merge(vFrom, vTo)
			if err != nil {
				return nil, err
			}
		}
		// if the key does not exist in the toMap, copy the value from the fromMap
		m[k] = vFrom
	}

	// go in the reverse direction
	for k, vTo := range toMap {
		if vFrom, ok := m[k]; ok {
			var err error
			vTo, err = Merge(vTo, vFrom)
			if err != nil {
				return nil, err
			}
		}
		m[k] = vTo
	}

	return m, nil
}
