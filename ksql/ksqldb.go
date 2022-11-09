package ksql

import (
	"fmt"
	"github.com/rmoff/ksqldb-go"
	"gitlab.com/dipper-iot/shared/logger"
	"reflect"
	"strings"
	"time"
)

type FieldData struct {
	Name  string
	Value interface{}
	Type  reflect.Type
}

/*
type USBLog struct {
	schema      struct{}  `name:"USB_LOG_STREAM" topic:"data.v1.usb_log" format:"JSON" type:"STREAM" partitions:"1"`
	Id          string    `json:"id" type:"VARCHAR KEY"`
	ClientId    string    `json:"client_id" type:"VARCHAR"`
	UsbId       string    `json:"usb_id" type:"VARCHAR"`
	Message     string    `json:"message" type:"VARCHAR"`
	Description string    `json:"description"  type:"VARCHAR"`
	Code        int32     `json:"code" type:"INT"`
	OnTime      time.Time `json:"on_time" type:"timestamp"`
	ServerTime  time.Time `json:"server_time,omitempty" type:"timestamp"`
	ServerId    string    `json:"server_id" type:"VARCHAR"`
}

type Application struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	InstallDate string `json:"install_date"`
	Location    string `json:"location"`
	Publisher   string `json:"publisher"`
}

type ApplicationLog struct {
	schema        struct{}       `name:"APPLICATION_LOG_STREAM" topic:"data.v1.application_log" format:"JSON" type:"STREAM" partitions:"1"`
	Id            string         `json:"id"  type:"VARCHAR KEY"`
	Applications  []*Application `json:"applications" type:"ARRAY<STRUCT<name VARCHAR,version VARCHAR, install_date VARCHAR,location VARCHAR, publisher VARCHAR>>"`
	ServerId      string         `json:"server_id" type:"VARCHAR"`
	ClientGroupId string         `json:"client_group_id" type:"VARCHAR"`
	OnTime        time.Time      `json:"on_time" type:"timestamp"`
	ServerTime    time.Time      `json:"server_time" type:"timestamp"`
}

*/

type FieldSchema struct {
	Name     string
	MetaData map[string]string
}

func ScanSchema(s interface{}) ([]FieldSchema, error) {
	list := make([]FieldSchema, 0)

	r := reflect.TypeOf(s)
	if r.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not is pointer")
	}
	dataType := r.Elem()

	for i := 0; i < dataType.NumField(); i++ {
		f := dataType.Field(i)

		metaData := make(map[string]string)
		name, success := f.Tag.Lookup("json")
		if !success {
			name = f.Name
		} else {
			name = strings.Split(name, ",")[0]
		}
		if f.Name == "schema" {
			name = "schema"

			val, success := f.Tag.Lookup("name")
			if success {
				metaData["name"] = val
			}

			val, success = f.Tag.Lookup("key")
			if success {
				metaData["key"] = val
			}

			val, success = f.Tag.Lookup("topic")
			if success {
				metaData["topic"] = val
			}

			val, success = f.Tag.Lookup("format")
			if success {
				metaData["format"] = val
			}

			val, success = f.Tag.Lookup("partitions")
			if success {
				metaData["partitions"] = val
			}

		}

		val, success := f.Tag.Lookup("type")
		if success {
			metaData["type"] = val
		}

		item := FieldSchema{
			Name:     name,
			MetaData: metaData,
		}
		list = append(list, item)
	}

	return list, nil
}

func ScanField(s interface{}) ([]FieldData, error) {
	list := make([]FieldData, 0)

	r := reflect.TypeOf(s)
	if r.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not is pointer")
	}
	dataType := r.Elem()
	val := reflect.ValueOf(s).Elem()

	for i := 0; i < dataType.NumField(); i++ {
		f := dataType.Field(i)
		v := val.Field(i)

		name, success := f.Tag.Lookup("json")

		if !success {
			name = f.Name
		} else {
			name = strings.Split(name, ",")[0]
		}
		var value interface{}
		if f.Name != "schema" {
			value = v.Interface()
		} else {
			name = "schema"
			value = f.Tag.Get("name")
		}
		item := FieldData{
			Name:  name,
			Type:  v.Type(),
			Value: value,
		}
		list = append(list, item)
	}

	return list, nil
}

func CreateSchema(s interface{}) (string, error) {

	list, err := ScanSchema(s)
	if err != nil {
		return "", err
	}

	schemaItem := make([]string, 0)
	optionItem := make([]string, 0)
	typeCreate := ""
	nameCreate := ""

	for _, data := range list {
		if data.Name == "schema" {
			typeCreate = data.MetaData["type"]
			nameCreate = data.MetaData["name"]

			key, success := data.MetaData["key"]
			if success {
				optionItem = append(optionItem, fmt.Sprintf("KEY='%s'", key))
			}

			topic := data.MetaData["topic"]
			optionItem = append(optionItem, fmt.Sprintf("KAFKA_TOPIC='%s'", topic))

			format := data.MetaData["format"]
			optionItem = append(optionItem, fmt.Sprintf("VALUE_FORMAT='%s'", format))

			partitions := data.MetaData["partitions"]
			optionItem = append(optionItem, fmt.Sprintf("PARTITIONS='%s'", partitions))
		} else {
			schemaItem = append(schemaItem, fmt.Sprintf("%s %s", data.Name, data.MetaData["type"]))
		}
	}

	return fmt.Sprintf("CREATE OR REPLACE %s IF NOT EXISTS %s (\n\t%s\n) WITH (%s);", typeCreate, nameCreate, strings.Join(schemaItem, ",\n\t"), strings.Join(optionItem, ",")), nil
}

func QueryInsert(s interface{}) (string, error) {

	list, err := ScanField(s)
	if err != nil {
		return "", err
	}

	header := make([]string, 0)
	values := make([]string, 0)
	table := ""

	for _, data := range list {
		if data.Name == "schema" {
			table = data.Value.(string)
		} else {
			header = append(header, data.Name)
			switch data.Type.Kind() {
			case reflect.String:
				{
					values = append(values, "'"+data.Value.(string)+"'")
				}
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int16, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint64, reflect.Uint16, reflect.Uint32:
				{
					values = append(values, fmt.Sprintf("%d", data.Value))
				}
			case reflect.Float32, reflect.Float64:
				{
					values = append(values, fmt.Sprintf("%f", data.Value))
				}
			case reflect.Map, reflect.Struct, reflect.Interface:
				{
					switch data.Value.(type) {
					case time.Time:
						dataTime := data.Value.(time.Time)
						values = append(values, fmt.Sprintf("'%s'", dataTime.UTC().Format("2006-01-02T15:04:05")))
					case *time.Time:
						dataTime := data.Value.(*time.Time)
						values = append(values, fmt.Sprintf("'%s'", dataTime.UTC().Format("2006-01-02T15:04:05")))

					default:
						rs, err := convertStructKSQL(data.Value)
						if err != nil {
							return "", err
						}
						values = append(values, rs)
					}

				}
			case reflect.Array, reflect.Slice:
				{
					y := make([]interface{}, 0)
					s := reflect.ValueOf(data.Value)
					for i := 0; i < s.Len(); i++ {
						y = append(y, s.Index(i).Interface())
					}

					rs, err := convertArrayKSQL(y)
					if err != nil {
						return "", err
					}
					values = append(values, rs)
				}
			default:
				logger.Infof(data.Type.String())
				values = append(values, fmt.Sprintf("%v", data.Value))
			}

		}
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", table, strings.Join(header, ","), strings.Join(values, ","))
	query = strings.ReplaceAll(query, "\\", "\\\\")
	return query, nil
}

func convertArrayKSQL(s []interface{}) (string, error) {
	values := make([]string, 0)
	if len(s) == 0 {
		return "", nil
	}
	for _, item := range s {
		switch item.(type) {
		case string:
			values = append(values, fmt.Sprintf(`'%s'`, item))
		case int, int8, int32, int64, uint, uint8, uint16, uint32, uint64:
			values = append(values, fmt.Sprintf(`%d`, item))
		case float64, float32:
			values = append(values, fmt.Sprintf(`%f`, item))
		default:
			rs, err := convertStructKSQL(item)
			if err != nil {
				return "", err
			}
			values = append(values, fmt.Sprintf(`%s`, rs))
		}
	}

	return fmt.Sprintf("ARRAY[%s]", strings.Join(values, ",")), nil
}

func convertStructKSQL(s interface{}) (string, error) {
	values := make([]string, 0)
	list, err := ScanField(s)
	if err != nil {
		return "", err
	}
	for _, data := range list {
		switch data.Type.Kind() {
		case reflect.String:
			{
				values = append(values, fmt.Sprintf("%s:='%s'", data.Name, data.Value))
			}
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int16, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint64, reflect.Uint16, reflect.Uint32:
			{
				values = append(values, fmt.Sprintf("%s:=%d", data.Name, data.Value))
			}
		case reflect.Float32, reflect.Float64:
			{
				values = append(values, fmt.Sprintf("%s:=%f", data.Name, data.Value))
			}
		case reflect.Struct, reflect.Interface:
			{
				switch data.Value.(type) {
				case time.Time:
					dataTime := data.Value.(time.Time)
					values = append(values, fmt.Sprintf("'%s'", dataTime.UTC().Format(time.RFC3339)))
				case *time.Time:
					dataTime := data.Value.(*time.Time)
					values = append(values, fmt.Sprintf("'%s'", dataTime.UTC().Format(time.RFC3339)))
				default:
					rs, err := convertStructKSQL(data.Value)
					if err != nil {
						return "", err
					}
					values = append(values, fmt.Sprintf("%s:=%s", data.Name, rs))
				}

			}
		case reflect.Array, reflect.Slice:
			{
				rs, err := convertArrayKSQL(data.Value.([]interface{}))
				if err != nil {
					return "", err
				}
				values = append(values, fmt.Sprintf(`%s:=%s`, data.Name, rs))
			}
		default:
			logger.Infof(data.Type.String())
			values = append(values, fmt.Sprintf("%s:=%v", data.Name, data.Value))
		}
	}

	return fmt.Sprintf("STRUCT(%s)", strings.Join(values, ",")), nil
}

func convertCharSpecial(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}

func RegistrySchema(client *ksqldb.Client, data ...interface{}) error {
	for _, item := range data {
		query, err := CreateSchema(item)
		if err != nil {
			logger.Error(err)
			return err
		}
		err = client.Execute(query)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}
