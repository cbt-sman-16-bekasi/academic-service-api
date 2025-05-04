package helper

import (
	"reflect"
	"strings"
)

func StructToFormMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		formTag := field.Tag.Get("form")
		if formTag == "" || formTag == "-" {
			continue
		}

		// Handle omitempty
		tagParts := strings.Split(formTag, ",")
		formKey := tagParts[0]
		omitEmpty := false
		for _, tagOpt := range tagParts[1:] {
			if tagOpt == "omitempty" {
				omitEmpty = true
			}
		}

		fieldValue := v.Field(i)
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				// Jangan masukkan field jika nil
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		// Skip zero value jika omitempty
		if omitEmpty && isZeroValue(fieldValue) {
			continue
		}

		result[formKey] = fieldValue.Interface()
	}

	return result
}

// isZeroValue memeriksa apakah nilai adalah zero value untuk tipe-nya
func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
