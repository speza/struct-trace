package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
)

type Span interface {
	SetTag(key string, value interface{})
}

func StructTrace(span Span, value interface{}) error {
	return structTrace("", span, value)
}

func structTrace(key string, span Span, value interface{}) error {
	val := reflect.ValueOf(value)

	if val.Kind() != reflect.Struct {
		return errors.New("value must be a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		fType := val.Type().Field(i)

		ignore := fType.Tag.Get("trace_ignore")
		if ignore == "true" {
			continue
		}

		name := fType.Name
		tag := fType.Tag.Get("trace")
		if tag != "" {
			name = tag
		} else {
			name = strcase.ToSnake(fType.Name)
		}

		f := val.Field(i)

		if key != "" {
			name = key + "." + name
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
			if err := structTrace(name, span, f.Interface()); err != nil {
				return fmt.Errorf("failed to parse nested struct %s", name)
			}
		}
	}

	return nil
}
