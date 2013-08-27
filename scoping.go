package mussed

import (
	"reflect"
)

func upscope(dot map[string]interface{}, subject interface{}) {
	var changes *mussedRecord
	switch {
	case mapType(subject):
		changes = upscopeMap(dot, subject)
	case structType(subject):
		changes = upscopeStruct(dot, subject)
	}

}
func upscopeMap(dot map[string]interface{}, subject interface{}) *mussedRecord {
	changes := &mussedRecord{
		replaced: map[string]interface{}{},
	}

	st := reflect.TypeOf(subject)
	if st.Kind() != reflect.Map {
		return
	}
	subjectValue := reflect.ValueOf(subject)
	if subjectValue.Type().Kind() != reflect.Map {
		fmt.Println("not a map")
		return
	}
	for _, key := range rv.MapKeys() {
		keyName := key.String()
		if val, ok := dot[keyName]; ok {
			changes.replaced[keyName] = val
		} else {
			changes.added = append(changes.added, keyName)
		}
		dot[key.String()] = subjectValue.MapIndex(key).Interface()
	}
}

func upscopeStruct(dot map[string]interface{}, subject interface{}) *mussedRecord {

}

func downscope(dot map[string]interface{}) {
	record, ok := dot["mussedScopeList"]
	if !ok {
		return
	}
	records, ok := record.([]*mussedRecord)
	if !ok {
		return
	}
	if len(records) > 0 {
		subject := records[len(records)-1]

		for _, added := range subject.added {
			delete(dot, added)
		}

		for key, previous := range subject.replaced {
			dot[key] = previous
		}
	}
}

type mussedRecord struct {
	replaced map[string]interface{}
	added    []string
}
