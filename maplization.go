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

// Conver2Map map[string]interface{}'s sugar
func (m *Maplization) Conver2Map(i interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(i)
	o, err := m.dispather(v)
	return o.(map[string]interface{}), err
}

// Conver conver 2 map[string]interface{}
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
	case reflect.Invalid:
		return nil, nil
	default:
		return m.safeInterface(v), nil
	}
}

func (m *Maplization) structHandler(v reflect.Value) (map[string]interface{}, error) {
	t := v.Type()
	o := make(map[string]interface{})
	var err error

	for i := 0; i < t.NumField(); i++ {
		var omitempty, skip, omitzero bool

		tagsName := t.Field(i).Name
		cur := v.Field(i)

		if tags, ok := t.Field(i).Tag.Lookup("bson"); ok {
			tagsArr := strings.Split(tags, ",")
			tagsName = tagsArr[0]
			for _, v := range tagsArr {
				switch v {
				case "omitempty":
					omitempty = true
				case "omitzero":
					omitzero = true
				case "-":
					skip = true
				}
			}
		}

		if skip || (cur.Kind() == reflect.Ptr && cur.IsNil() && omitempty) || (cur.IsZero() && omitzero) {
			continue
		}

		ele, err := m.Conver(m.safeInterface(cur))

		if err != nil {
			return nil, err
		}

		if formatter, ok := t.Field(i).Tag.Lookup("formatter"); ok {
			o[tagsName], err = m.Formatters[formatter](ele)
		} else {
			o[tagsName], err = m.Conver(ele)
		}

	}

	return o, err
}

func (m *Maplization) mapHandler(v reflect.Value) (map[string]interface{}, error) {
	var err error
	o := make(map[string]interface{})

	for _, idx := range v.MapKeys() {
		o[idx.String()], err = m.Conver(m.safeInterface(v.MapIndex(idx)))
	}
	return o, err
}

func (m *Maplization) sliceHandler(v reflect.Value) ([]interface{}, error) {
	var err error
	l := v.Len()
	o := make([]interface{}, l)

	for i := 0; i < l; i++ {
		o[i], err = m.Conver(m.safeInterface(v.Index(i)))
	}
	return o, err
}

func (m *Maplization) ptrHandler(v reflect.Value) (interface{}, error) {
	return m.Conver(m.safeInterface(reflect.Indirect(v)))
}

func (m *Maplization) safeInterface(v reflect.Value) interface{} {
	if v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return nil
}
