package maplization

import (
	"reflect"
	"strings"
)

// Formatter is custom format func
type Formatter func(interface{}) (interface{}, error)

// Maplization can conver struct 2 map
type Maplization struct {
	Formatters map[string]Formatter
}

// NewMapper init func
func NewMapper(m map[string]Formatter) *Maplization {
	return &Maplization{
		Formatters: m,
	}
}

// Conver 2 map
func (m *Maplization) Conver(i interface{}) (interface{}, error) {
	v := reflect.ValueOf(i)
	return m.dispather(v)
}

func (m *Maplization) dispather(v reflect.Value) (interface{}, error) {
	switch v.Kind() {
	case reflect.Struct:
		return m.structHandler(v)
	case reflect.Map:
		return m.mapHandler(v)
	case reflect.Slice:
		return m.sliceHandler(v)
	case reflect.Ptr:
		return m.ptrHandler(v)
	default:
		if v.CanInterface() {
			return v.Interface(), nil
		}
		return v, nil
	}
}

func (m *Maplization) structHandler(v reflect.Value) (o map[string]interface{}, err error) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		var canEmpty, skip bool

		tagsName := t.Field(i).Name
		curElem := v.Field(i)

		if tags, ok := t.Field(i).Tag.Lookup("bson"); ok {
			tagsArr := strings.Split(tags, ",")
			tagsName = tagsArr[0]
			for _, v := range tagsArr {
				switch v {
				case "omitempty":
					canEmpty = true
				case "-":
					skip = true
				}
			}
		}

		if skip || (curElem.IsNil() && canEmpty) {
			continue
		}

		if formatter, ok := t.Field(i).Tag.Lookup("formatter"); ok {
			o[tagsName], err = m.Formatters[formatter](curElem.Interface())
		} else {
			o[tagsName], err = m.dispather(curElem)
		}

	}

	return o, err
}

func (m *Maplization) mapHandler(v reflect.Value) (o map[string]interface{}, err error) {
	for _, idx := range v.MapKeys() {
		o[idx.Interface().(string)], err = m.dispather(v.MapIndex(idx))
	}
	return o, err
}

func (m *Maplization) sliceHandler(v reflect.Value) (o []interface{}, err error) {
	for i := 0; i < v.Len(); i++ {
		o[i], err = m.dispather(v.Index(i))
	}
	return o, err
}

func (m *Maplization) ptrHandler(v reflect.Value) (interface{}, error) {
	if v.IsNil() {
		return nil, nil
	}
	return m.dispather(v.Elem())
}
