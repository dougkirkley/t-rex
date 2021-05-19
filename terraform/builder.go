package terraform

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
)

func dereferenceFieldIfPointer(field reflect.StructField) {}

// BuildArgs builds args
func BuildArgs(x interface{}) []string {

	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	flags := make([]string, 0)

	positionalArgs := make(map[int]string)

	for i := 0; i < t.NumField(); i++ {

		fieldType := t.Field(i)
		fieldValue := v.Field(i)
		fieldKind := fieldType.Type.Kind()

		flag := fieldType.Tag.Get("flag")
		pos := fieldType.Tag.Get("pos")
		if flag == "" && pos == "" {
			continue
		}

		rootKind := fieldType.Type.Kind()
		rootValue := fieldValue

		if fieldKind == reflect.Ptr {
			rootKind = fieldType.Type.Elem().Kind()
			rootValue = fieldValue.Elem()
		}

		if rootValue.Kind() == reflect.Invalid {
			continue
		}

		if pos != "" {
			index, err := strconv.Atoi(pos)
			if err != nil {
				log.Fatal(err)
			}
			positionalArgs[index] = rootValue.String()
			continue
		}

		switch rootKind {
		case reflect.String:
			flags = append(flags, fmt.Sprintf("%s=%v", flag, rootValue.String()))

		case reflect.Bool:
			if rootValue.Bool() {
				flags = append(flags, fmt.Sprintf("%s", flag))
			} else {
				flags = append(flags, fmt.Sprintf("%s=%v", flag, rootValue.Bool()))
			}
		case reflect.Map:
			for _, keyValue := range rootValue.MapKeys() {
				valueValue := rootValue.MapIndex(keyValue)
				flags = append(flags, fmt.Sprintf("%s '%s=%s'", flag, keyValue.String(), valueValue.String()))
			}
		case reflect.Slice:
			for j := 0; j < rootValue.Len(); j++ {

				flags = append(flags, fmt.Sprintf("%s='%s'", flag, rootValue.Index(j).String()))
			}
		}

	}

	for _, key := range sortedMapKeys(positionalArgs) {
		flags = append(flags, positionalArgs[key])
	}

	return flags
}

func sortedMapKeys(m map[int]string) []int {
	keys := make([]int, 0)
	for key := range m {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys
}
