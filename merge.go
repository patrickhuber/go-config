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

	// if from is nil, return to
	if fromMap == nil {
		return to, nil
	}

	// check if to is a map, if not return to
	toMap, ok := to.(map[string]any)
	if !ok {
		return to, nil
	}

	// deepCopy the toMap as it is the target
	merged := make(map[string]any, len(toMap))
	for keyTo, valueTo := range toMap {
		merged[keyTo] = deepCopy(valueTo)
	}

	// merge keys in the from map with keys in the to map
	for keyFrom, valueFrom := range fromMap {
		// if the key exists in both maps, merge the values
		if vTo, ok := merged[keyFrom]; ok {
			var err error
			valueFrom, err = Merge(valueFrom, vTo)
			if err != nil {
				return nil, err
			}
		}
		// if the key does not exist in the toMap, copy the value from the fromMap
		merged[keyFrom] = valueFrom
	}

	return merged, nil
}

func deepCopy(value any) any {
	if value == nil {
		return nil
	}
	switch concrete := value.(type) {
	case []any:
		sliceClone := make([]any, len(concrete))
		for i, v := range concrete {
			sliceClone[i] = deepCopy(v)
		}
		return sliceClone
	case map[string]any:
		mapClone := make(map[string]any, len(concrete))
		for k, v := range concrete {
			mapClone[k] = deepCopy(v)
		}
		return mapClone
	case string:
		return concrete
	case float64:
		return concrete
	case bool:
		return concrete
	default:
		return concrete
	}
}
