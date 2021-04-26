package utils

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// Get field form
func GetForm(field reflect.StructField) string {
	tag := field.Tag.Get("form")
	if tag == "" {
		tag = field.Tag.Get("json")
	}
	desc := strings.Trim(tag, " ")
	return desc
}

// Get field description
func GetDescription(field reflect.StructField) string {
	tag := field.Tag.Get("description")
	desc := strings.Trim(tag, " ")
	return desc
}

// Get require
func GetRequired(field reflect.StructField) bool {
	tag := field.Tag.Get("binding")
	return strings.Contains(tag, "required")
}

// Traverse the structure for processing
func RangeStruct(s interface{}, f func(reflect.Value, reflect.StructField) bool) error {
	if s == nil {
		return errors.New("the params is nil")
	}
	val := reflect.ValueOf(s)

	if !val.IsValid() {
		return errors.New("the params is error")
	}
	if val.Kind() == reflect.Slice {
		return nil
	}
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice {
		val = reflect.New(val.Type().Elem())
	}

	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	typ := val.Type()
	cnt := val.NumField()
	for i := 0; i < cnt; i++ {
		fd := val.Field(i)
		ty := typ.Field(i)
		if !f(fd, ty) {
			break
		}
	}
	return nil
}
