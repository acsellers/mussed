package mussed

import (
	"reflect"
)

func upscope(d interface{}, i interface{}) map[string]interface{} {
	dot, ok := d.(map[string]interface{})
	if !ok {
		dot = make(map[string]interface{})
	} else {
		if dot == nil {
			dot = make(map[string]interface{})
		}
	}
	var changes *mussedRecord
	subject := i

	switch {
	case mapType(subject):
		changes = upscopeMap(dot, subject)
	case structType(subject):
		changes = upscopeStruct(dot, subject)
	default:
		dot["mussedItem"] = i
		return dot
	}
	if list, ok := dot["mussedScopeList"]; ok {
		if records, ok := list.([]*mussedRecord); ok {
			dot["mussedScopeList"] = append(records, changes)
		} else {
			dot["mussedScopeList"] = []*mussedRecord{changes}
		}
	} else {
		dot["mussedScopeList"] = []*mussedRecord{changes}
	}

	return dot
}

func mapType(i interface{}) bool {
	if i != nil {
		return reflect.TypeOf(i).Kind() == reflect.Map
	}
	return false
}

func structType(i interface{}) bool {
	it := reflect.TypeOf(i)
	if it.Kind() == reflect.Ptr {
		it = it.Elem()
	}
	return it.Kind() == reflect.Struct
}
func upscopeMap(dot map[string]interface{}, subject interface{}) *mussedRecord {
	changes := &mussedRecord{
		replaced: map[string]interface{}{},
	}

	subjectValue := reflect.ValueOf(subject)
	if subjectValue.Type().Kind() != reflect.Map {
		return nil
	}
	for _, key := range subjectValue.MapKeys() {
		keyName := key.String()
		if val, ok := dot[keyName]; ok {
			changes.replaced[keyName] = val
		} else {
			changes.added = append(changes.added, keyName)
		}
		dot[key.String()] = subjectValue.MapIndex(key).Interface()
	}
	return changes
}

func upscopeStruct(dot map[string]interface{}, subject interface{}) *mussedRecord {
	return nil
}

func downscope(dot map[string]interface{}) map[string]interface{} {
	record, ok := dot["mussedScopeList"]
	if !ok {
		return dot
	}
	records, ok := record.([]*mussedRecord)
	if !ok {
		return dot
	}
	if len(records) > 0 {
		subject := records[len(records)-1]

		if subject == nil {
			return dot
		}
		for _, added := range subject.added {
			delete(dot, added)
		}

		for key, previous := range subject.replaced {
			dot[key] = previous
		}
	}

	return dot
}

type mussedRecord struct {
	replaced map[string]interface{}
	added    []string
}
