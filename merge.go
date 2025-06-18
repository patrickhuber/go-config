package config

import (
	"maps"
)

func Merge(from any, to any) (any, error) {
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

	// check if to is a map, if not return to
	toMap, ok := to.(map[string]any)
	if !ok {
		return to, nil
	}

	// clone the from map to avoid mutating
	m := maps.Clone(fromMap)

	// overwrite keys from fromMap with keys from toMap
	maps.Copy(m, toMap)

	return m, nil
}
