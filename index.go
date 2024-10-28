package gu

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

type Index interface {
	/*
		/获取index名
	*/
	GetIndexName() string

	/*
		/获取struct,实现这个函数需要返回值而不是指针
	*/
	GetStruct() any
}

type ESIndex struct {
	Index Index
}

func CreateIndex(index ...Index) error {
	var success = make([]string, 0, len(index))
	for i := range index {
		name := index[i].GetIndexName()
		if name == "" {
			logrus.Errorf("index name is empty")
			continue
		}

		mapping := GenerateIndexMapping(index[i].GetStruct())

		if mapping == nil {
			logrus.Errorf("generate mapping failed,index:%v", name)
			continue
		}

		if _, err := EsClient.CreateIndex(name).
			BodyJson(map[string]any{
				"mappings": map[string]any{
					"properties": mapping,
				},
			}).Do(Timeout()); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}

		success = append(success, name)
	}

	logrus.Infof("create index success:%v", success)

	return nil
}

func GenerateIndexMapping(data interface{}, granularity ...bool) map[string]any {
	typ := reflect.TypeOf(data)
	if typ.Kind() != reflect.Struct {
		return nil
	}

	properties := make(map[string]any)

	value := reflect.ValueOf(data)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}

		//time.Time类型 自用
		if tag == "create" || tag == "modified" {
			properties[tag] = map[string]interface{}{
				"type":   "date",
				"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd",
			}
			continue
		}

		if field.Type == reflect.TypeOf(time.Time{}) {
			//time类型
			properties[tag] = map[string]interface{}{
				"type":   "date",
				"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||yyyy||yyyy-M||yyyy-M-d||yyyy.M||yyyy.M.d||yyyy-MM-dd'T'HH:mm:ss||yyyy-MM-dd'T'HH:mm:ss.SSSSSSS",
			}
			continue
		}

		switch field.Type.Kind() {

		case reflect.String:
			if len(granularity) > 0 && granularity[0] {
				properties[tag] = map[string]interface{}{
					"type": "keyword",
					"fields": map[string]any{
						"ik_max": map[string]any{
							"type":            "text",
							"analyzer":        "ik_max_word",
							"search_analyzer": "ik_max_word",
						},
						"ik_smart": map[string]any{
							"type":     "text",
							"analyzer": "ik_smart_word",
						},
					},
				}
			} else {
				properties[tag] = map[string]interface{}{
					"type":     "text",
					"analyzer": "ik_max_word",
					"fields": map[string]any{
						"keyword": map[string]any{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			properties[tag] = map[string]interface{}{
				"type": "long",
			}
		case reflect.Float32, reflect.Float64:
			properties[tag] = map[string]interface{}{
				"type": "double",
			}
		case reflect.Struct:
			if field.Anonymous {
				nestedMapping := GenerateIndexMapping(value.Field(i).Interface())
				for key, value := range nestedMapping {
					properties[key] = value
				}
				continue
			}
		case reflect.Slice:
			switch field.Type.Elem().Kind() {
			case reflect.String:
				properties[tag] = map[string]interface{}{
					"type":     "text",
					"analyzer": "ik_max_word",
					"fields": map[string]any{
						"keyword": map[string]any{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				properties[tag] = map[string]interface{}{
					"type": "long",
				}
			case reflect.Float32, reflect.Float64:
				properties[tag] = map[string]interface{}{
					"type": "double",
				}
			case reflect.Struct:
				elemType := reflect.Indirect(reflect.New(field.Type.Elem())).Interface()
				nestedMapping := GenerateIndexMapping(elemType)
				properties[tag] = map[string]any{
					"type":       "nested",
					"properties": nestedMapping,
				}
			}
		default:
			continue
		}
	}

	return properties
}
