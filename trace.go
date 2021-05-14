package structtrace

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

const (
	keyTagField = "key="
	ignoreTag   = "ignore"
)

// Span represents the tracing span.
type Span interface {
	SetTag(key string, value interface{})
}

// StructTrace starts the struct scan recursively processing any nested structs.
func StructTrace(span Span, value interface{}) {
	structTrace("", span, value)
}

func structTrace(baseKey string, span Span, value interface{}) {
	val := reflect.ValueOf(value)

	// Only allow us to trace a struct.
	if val.Kind() != reflect.Struct {
		return
	}

fieldItr:
	for i := 0; i < val.NumField(); i++ {
		fType := val.Type().Field(i)

		name := fType.Name
		tag := fType.Tag.Get("trace")
		items := strings.Split(tag, ",")
		for i := range items {
			if tag == ignoreTag {
				continue fieldItr
			}
			if strings.HasPrefix(items[i], keyTagField) {
				name = strings.TrimPrefix(tag, keyTagField)
			} else {
				name = strcase.ToSnake(fType.Name)
			}
		}

		f := val.Field(i)

		if baseKey != "" {
			name = baseKey + "." + name
		}

		if f.Kind() == reflect.Ptr {
			if f.IsNil() {
				span.SetTag(name, "nil")
				continue
			}
			f = f.Elem()
		}

		switch f.Kind() {
		case reflect.String:
			span.SetTag(name, f.String())
		case reflect.Bool:
			span.SetTag(name, f.Bool())
		case reflect.Float32, reflect.Float64:
			span.SetTag(name, f.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			span.SetTag(name, f.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			span.SetTag(name, f.Uint())
		case reflect.Struct:
			structTrace(name, span, f.Interface())
		}
	}
}
