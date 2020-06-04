package maplization

import (
	"log"
	"testing"
	"time"
)

// Test struct
type Test struct {
	ID           *string                 `bson:"id,omitempty"`
	String       string                  `bson:"string,omitzero"`
	Arr          []interface{}           `bson:"arr,omitzero"`
	ArrZero      []interface{}           `bson:"arr_zero,omitzero"`
	ArrPtr       *[]interface{}          `bson:"arr_ptr,omitempty"`
	ArrPtrNil    *[]interface{}          `bson:"arr_ptr_nil,omitempty"`
	Map          map[string]interface{}  `bson:"map,omitzero"`
	MapZero      map[string]interface{}  `bson:"map_zero,omitzero"`
	MapPtr       *map[string]interface{} `bson:"map_ptr,omitempty"`
	MapPtrNil    *map[string]interface{} `bson:"map_ptr_nil,omitempty"`
	I32          int32                   `bson:"i32,omitzero"`
	I32Ptr       *int32                  `bson:"i32_ptr,omitempty"`
	F32          float32                 `bson:"f32,omitzero"`
	F32Ptr       *float32                `bson:"f32_ptr,omitempty"`
	Time         time.Time               `bson:"time,omitzero" formatter:"now"`
	TimePtr      *time.Time              `bson:"time_ptr,omitempty" formatter:"now"`
	Struct       *Test                   `bson:"struct,omitempty"`
	StructNil    *Test                   `bson:"struct_nil,omitempty"`
	OmitZero     string                  `bson:"omit_zero"`
	OmitEmpty    *string                 `bson:"omit_empty"`
	InterfacePtr *interface{}            `bson:"interface_ptr,omitempty"`
}

// mapliztion test
func TestMapliztion(test *testing.T) {
	formatter := map[string]Formatter{
		"now": func(i interface{}) (interface{}, error) {
			return time.Now().Local(), nil
		},
	}

	mapper := NewMapper(formatter)

	var str string = "5e73018c324cecdcfda7bcac"
	var id *string = &str
	var arrZero []interface{}
	var arrPtrNil *[]interface{}
	var mapZero map[string]interface{}
	var mapPtrNil *map[string]interface{}
	var i32 int32 = 32
	var i32Ptr *int32 = &i32
	var f32 float32 = 32.0
	var f32Ptr *float32 = &f32
	var t time.Time = time.Now()
	var timePtr *time.Time = &t
	var structNil *Test
	var omitZero string
	var omitEmpty *string
	var it interface{}
	var interfacePtr *interface{} = &it

	subStr := &Test{
		String: "subStr",
	}

	arr := []interface{}{
		str,
		id,
		arrZero,
		arrPtrNil,
		mapZero,
		mapPtrNil,
		i32,
		i32Ptr,
		f32,
		f32Ptr,
		t,
		timePtr,
		subStr,
		structNil,
		omitZero,
		omitEmpty,
	}
	var arrPtr *[]interface{} = &arr
	m := map[string]interface{}{
		"str,":       str,
		"id,":        id,
		"arrZero,":   arrZero,
		"arrPtrNil,": arrPtrNil,
		"mapZero,":   mapZero,
		"mapPtrNil,": mapPtrNil,
		"i32,":       i32,
		"i32Ptr,":    i32Ptr,
		"f32,":       f32,
		"f32Ptr,":    f32Ptr,
		"t,":         t,
		"timePtr,":   timePtr,
		"subStr":     subStr,
		"structNil,": structNil,
		"omitZero,":  omitZero,
		"omitEmpty,": omitEmpty,
	}
	var mapPtr *map[string]interface{} = &m

	tt := &Test{
		ID:        id,
		String:    str,
		Arr:       arr,
		ArrZero:   arrZero,
		ArrPtr:    arrPtr,
		ArrPtrNil: arrPtrNil,
		Map:       m,
		MapZero:   mapZero,
		MapPtr:    mapPtr,
		MapPtrNil: mapPtrNil,
		I32:       i32,
		I32Ptr:    i32Ptr,
		F32:       f32,
		F32Ptr:    f32Ptr,
		Time:      t,
		TimePtr:   timePtr,
		Struct:    subStr,
		StructNil: structNil,
		// OmitZero:  omitZero,
		// OmitEmpty: omitEmpty,
		InterfacePtr: interfacePtr,
	}

	mm, err := mapper.Conver2Map(tt)
	mm["set"] = "set"

	if err != nil {
		test.Log(err)
	}

	log.Printf("%+v\n", mm)
	test.Log(mm)

}
