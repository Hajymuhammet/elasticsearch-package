package mapping

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// BuildMappingFromStruct generates Elasticsearch index mapping from Go struct with JSON and ES tags.
// It also adds default analyzers to handle multilingual search.
func BuildMappingFromStruct[T any](sample T) (map[string]any, error) {
	typ := reflect.TypeOf(sample)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("BuildMappingFromStruct: expected struct, got %s", typ.Kind())
	}

	properties := make(map[string]any)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		jsonTag = strings.Split(jsonTag, ",")[0] // <- düzediş
		esTag := field.Tag.Get("es")

		esField := make(map[string]any)

		if esTag != "" {
			tags := strings.Split(esTag, ",")
			for _, tag := range tags {
				parts := strings.SplitN(tag, "=", 2)
				key := strings.TrimSpace(parts[0])
				if len(parts) == 2 {
					esField[key] = strings.TrimSpace(parts[1])
				} else {
					esField[key] = true
				}
			}
		}

		if _, ok := esField["type"]; !ok {
			switch field.Type.Kind() {
			case reflect.String:
				esField["type"] = "text"
				esField["analyzer"] = "universal_analyzer"
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				esField["type"] = "long"
			case reflect.Float32, reflect.Float64:
				esField["type"] = "float"
			case reflect.Bool:
				esField["type"] = "boolean"
			case reflect.Struct:
				if field.Type == reflect.TypeOf(time.Time{}) {
					esField["type"] = "date"
				} else {
					esField["type"] = "object"
				}
			}
		}

		properties[jsonTag] = esField
	}

	mapping := map[string]any{
		"settings": map[string]any{
			"analysis": map[string]any{
				"analyzer": map[string]any{
					"universal_analyzer": map[string]any{
						"tokenizer": "icu_tokenizer",
						"filter":    []string{"icu_folding", "lowercase"},
					},
				},
			},
		},
		"mappings": map[string]any{
			"properties": properties,
		},
	}

	return mapping, nil
}
