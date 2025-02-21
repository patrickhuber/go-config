package config

import (
	"fmt"
	"maps"
	"strconv"
)

type ChangeType interface {
	changeType()
}

type changeType string

func (changeType) changeType() {}

const (
	Create changeType = "create"
	Update changeType = "update"
	Delete changeType = "delete"
)

type Change struct {
	ChangeType ChangeType
	Path       []string
	From       any
	To         any
}

func Diff(from, to any) ([]Change, error) {
	return diff(from, to, []string{})
}

func diff(from, to any, path []string) ([]Change, error) {
	if from == nil && to == nil {
		return nil, nil
	}
	if from == nil && to != nil {
		return []Change{
			{
				ChangeType: Create,
				Path:       path,
				From:       from,
				To:         to,
			},
		}, nil
	}
	if from != nil && to == nil {
		return []Change{
			{
				ChangeType: Delete,
				Path:       path,
				From:       from,
				To:         to,
			},
		}, nil
	}
	switch f := from.(type) {
	case string:
		return diffOf(f, to, path), nil
	case float64:
		return diffOf(f, to, path), nil
	case bool:
		return diffOf(f, to, path), nil
	case map[string]any:
		return mapDiff(f, to, path)
	case []any:
		return sliceDiff(f, to, path)
	}
	return nil, fmt.Errorf("unrecognized source type %T. expected string, float64, bool, map[string]any or []any", from)
}

func diffOf[T comparable](from T, to any, path []string) []Change {
	t, isTypeT := to.(T)
	if isTypeT && t == from {
		return nil
	}
	return []Change{
		{
			ChangeType: Update,
			Path:       path,
			From:       from,
			To:         to,
		},
	}
}

func mapDiff(fromMap map[string]any, to any, path []string) ([]Change, error) {
	toMap, toIsMap := to.(map[string]any)
	if !toIsMap {
		return []Change{
			{
				ChangeType: Update,
				Path:       path,
				From:       fromMap,
				To:         to,
			},
		}, nil
	}
	var changes []Change

	// process deletions and updates
	for fromKey := range maps.Keys(fromMap) {
		value, exists := toMap[fromKey]
		if !exists {
			value = nil
		}
		recurseChanges, err := diff(fromMap[fromKey], value, append(path, fromKey))
		if err != nil {
			return nil, err
		}
		changes = append(changes, recurseChanges...)
	}
	// process additions
	for toKey := range maps.Keys(toMap) {
		_, exists := fromMap[toKey]
		if exists {
			continue
		}
		recurseChanges, err := diff(nil, toMap[toKey], append(path, toKey))
		if err != nil {
			return nil, err
		}
		changes = append(changes, recurseChanges...)
	}
	return changes, nil
}

func sliceDiff(fromSlice []any, to any, path []string) ([]Change, error) {
	toSlice, toIsSlice := to.([]any)
	if !toIsSlice {
		return []Change{
			{
				ChangeType: Update,
				Path:       path,
				From:       fromSlice,
				To:         toSlice,
			},
		}, nil
	}
	var changes []Change
	f, t := 0, 0
	for {
		if t >= len(toSlice) && f >= len(fromSlice) {
			break
		} else if t >= len(toSlice) {
			for ; f < len(fromSlice); f++ {
				change := Change{
					ChangeType: Delete,
					Path:       append(path, strconv.Itoa(f)),
					From:       fromSlice[f],
					To:         nil}
				changes = append(changes, change)
			}
			break
		} else if f == len(fromSlice) {
			for ; t < len(toSlice); t++ {
				change := Change{
					ChangeType: Create,
					Path:       append(path, strconv.Itoa(f)),
					From:       nil,
					To:         toSlice[t]}
				changes = append(changes, change)
			}
			break
		} else {
			recurseChanges, err := diff(fromSlice[f], toSlice[t], append(path, strconv.Itoa(f)))
			if err != nil {
				return nil, err
			}
			changes = append(changes, recurseChanges...)
			f++
			t++
		}
	}

	return changes, nil
}