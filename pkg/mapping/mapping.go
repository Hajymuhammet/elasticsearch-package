package mapping

import (
	"reflect"
	"strings"
)

func BuildMappingFromStruct(sample any) map[string]any {
	properties := map[string]any{}
	v := reflect.ValueOf(sample)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return map[string]any{}
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name := f.Tag.Get("json")
		if name == "" {
			name = f.Name
		}
		name = strings.Split(name, ",")[0]
		esTag := f.Tag.Get("es")
		prop := map[string]any{}
		if esTag != "" {
			parts := strings.Split(esTag, ",")
			for _, p := range parts {
				kv := strings.SplitN(p, "=", 2)
				if len(kv) == 2 {
					prop[kv[0]] = kv[1]
				} else {
					prop[p] = true
				}
			}
		}
		if _, ok := prop["type"]; !ok {
			switch f.Type.Kind() {
			case reflect.String:
				prop["type"] = "text"
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				prop["type"] = "long"
			case reflect.Float32, reflect.Float64:
				prop["type"] = "double"
			case reflect.Bool:
				prop["type"] = "boolean"
			default:
				prop["type"] = "object"
			}
		}
		properties[name] = prop
	}
	return map[string]any{"mappings": map[string]any{"properties": properties}}
}
