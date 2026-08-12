package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Format(v any, formatType string) (string, error) {
	switch strings.ToLower(formatType) {
	case "json":
		return FormatJSON(v)
	case "table", "":
		return FormatTable(v)
	default:
		return "", fmt.Errorf("unsupported output format: %s", formatType)
	}
}

func FormatJSON(v any) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FormatTable(v any) (string, error) {
	val := reflect.ValueOf(v)
	var sliceVal reflect.Value
	var elemType reflect.Type

	switch val.Kind() {
	case reflect.Struct:
		sliceVal = reflect.MakeSlice(reflect.SliceOf(val.Type()), 1, 1)
		sliceVal.Index(0).Set(val)
		elemType = val.Type()
	case reflect.Ptr:
		if val.Elem().Kind() == reflect.Struct {
			sliceVal = reflect.MakeSlice(reflect.SliceOf(val.Type()), 1, 1)
			sliceVal.Index(0).Set(val)
			elemType = val.Elem().Type()
		} else {
			return "", fmt.Errorf("provided argument is not a struct")
		}
	case reflect.Slice:
		if val.Len() == 0 {
			return "No resources found.\n", nil
		}
		elemType = val.Index(0).Type()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() != reflect.Struct {
			return "", fmt.Errorf("slice elements must be structs or pointers to structs")
		}
		sliceVal = val
	default:
		return "", fmt.Errorf("provided argument is not a struct or slice")
	}

	t := table.NewWriter()

	numFields := elemType.NumField()

	headers := make(table.Row, numFields)
	for i := range numFields {
		headers[i] = strings.ToUpper(elemType.Field(i).Name)
	}
	t.AppendHeader(headers)

	for r := range sliceVal.Len() {
		row := make(table.Row, numFields)
		item := reflect.Indirect(sliceVal.Index(r))

		for c := range numFields {
			row[c] = item.Field(c).Interface()
		}
		t.AppendRow(row)
	}

	return t.Render() + "\n", nil
}
