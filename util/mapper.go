package util

import (
	"encoding/json"
	"errors"
	"gitlab.com/dipper-iot/shared/convert"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	ZeroValue    reflect.Value
	fieldNameMap sync.Map
	timeType     = reflect.TypeOf(time.Now())
)

func Mapper(objectFrom interface{}, objectTo interface{}) error {
	buffer, err := json.Marshal(objectFrom)
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, objectTo)
}

func CheckExistsField(elem reflect.Value, fieldName string) (realFieldName string, exists bool) {
	realRaw, isOk := fieldNameMap.Load(elem)
	if !isOk {
		return "", false
	}
	mapName := realRaw.(map[string]string)
	name, exits := mapName[fieldName]
	if !exits {
		return "", false
	}
	return name, true
}

func setFieldValue(realFieldName string, fieldValue reflect.Value, fieldKind reflect.Kind, value interface{}) error {

	switch fieldKind {
	case reflect.Bool:
		if value == nil {
			fieldValue.SetBool(false)
		} else if v, ok := value.(bool); ok {
			fieldValue.SetBool(v)
		} else {
			v, _ := convert.Convert(convert.ToString(value)).Bool()
			fieldValue.SetBool(v)
		}

	case reflect.String:
		if value == nil {
			fieldValue.SetString("")
		} else {
			fieldValue.SetString(convert.ToString(value))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == nil {
			fieldValue.SetInt(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetInt(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetInt(int64(val.Uint()))
			default:
				v, _ := convert.Convert(convert.ToString(value)).Int64()
				fieldValue.SetInt(v)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == nil {
			fieldValue.SetUint(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetUint(uint64(val.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetUint(val.Uint())
			default:
				v, _ := convert.Convert(convert.ToString(value)).Uint64()
				fieldValue.SetUint(v)
			}
		}
	case reflect.Float64, reflect.Float32:
		if value == nil {
			fieldValue.SetFloat(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Float64:
				fieldValue.SetFloat(val.Float())
			default:
				v, _ := convert.Convert(convert.ToString(value)).Float64()
				fieldValue.SetFloat(v)
			}
		}
	case reflect.Struct:
		// struct
	default:
		if reflect.ValueOf(value).Type() == fieldValue.Type() {
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	return nil
}

func RegisterMapper(val reflect.Value) {
	_, exits := fieldNameMap.Load(val)
	if exits {
		return
	}
	mapName := make(map[string]string)

	for i := 0; i < val.Type().NumField(); i++ {
		nameJson := val.Type().Field(i).Tag.Get("json")
		nameJson = strings.ReplaceAll(nameJson, "omitempty", "")
		nameJson = strings.ReplaceAll(nameJson, ",", "")
		name := val.Type().Field(i).Name
		if nameJson == "id" {
			mapName["_id"] = name
		}
		mapName[nameJson] = name
	}

	fieldNameMap.Store(val, mapName)
}

func MapperMap(fromMap map[string]interface{}, toObj interface{}) error {

	toElem := reflect.ValueOf(toObj).Elem()
	if toElem == ZeroValue {
		return errors.New("to obj is not legal value")
	}
	//check register flag
	RegisterMapper(toElem)
	//if not register, register it
	for k, v := range fromMap {
		fieldName := k
		//check field is exists
		realFieldName, exists := CheckExistsField(toElem, fieldName)
		if !exists {
			continue
		}
		fieldInfo, exists := toElem.Type().FieldByName(realFieldName)
		if !exists {
			continue
		}
		fieldKind := fieldInfo.Type.Kind()
		fieldValue := toElem.FieldByName(realFieldName)
		setFieldValue(realFieldName, fieldValue, fieldKind, v)
	}
	return nil
}
