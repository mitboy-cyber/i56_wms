// Package validate provides struct-tag-based validation for API requests.
package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Struct validates a struct pointer using "validate" tags.
// Returns nil if all validations pass, or a map of field→error.
func Struct(v interface{}) map[string]string {
	errs := map[string]string{}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		errs["_"] = "expected a struct"
		return errs
	}
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		jsonName := field.Tag.Get("json")
		if jsonName == "" {
			jsonName = field.Name
		} else if idx := strings.Index(jsonName, ","); idx > 0 {
			jsonName = jsonName[:idx]
		}
		fieldVal := val.Field(i)
		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if msg := checkRule(rule, fieldVal, jsonName); msg != "" {
				errs[jsonName] = msg
				break // one error per field
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func checkRule(rule string, val reflect.Value, name string) string {
	// Handle numeric types
	switch val.Kind() {
	case reflect.String:
		s := val.String()
		switch {
		case rule == "required" && s == "":
			return fmt.Sprintf("%s 为必填项", name)
		case strings.HasPrefix(rule, "min="):
			n, _ := strconv.Atoi(rule[4:])
			if len(s) < n && s != "" {
				return fmt.Sprintf("%s 长度不能少于 %d 个字符", name, n)
			}
		case strings.HasPrefix(rule, "max="):
			n, _ := strconv.Atoi(rule[4:])
			if len(s) > n {
				return fmt.Sprintf("%s 长度不能超过 %d 个字符", name, n)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := val.Int()
		switch {
		case rule == "required" && v == 0:
			return fmt.Sprintf("%s 为必填项", name)
		case strings.HasPrefix(rule, "gt="):
			n, _ := strconv.ParseInt(rule[3:], 10, 64)
			if v <= n {
				return fmt.Sprintf("%s 必须大于 %d", name, n)
			}
		case strings.HasPrefix(rule, "gte="):
			n, _ := strconv.ParseInt(rule[4:], 10, 64)
			if v < n {
				return fmt.Sprintf("%s 必须 >= %d", name, n)
			}
		case strings.HasPrefix(rule, "lte="):
			n, _ := strconv.ParseInt(rule[4:], 10, 64)
			if v > n {
				return fmt.Sprintf("%s 必须 <= %d", name, n)
			}
		}
	case reflect.Float64, reflect.Float32:
		v := val.Float()
		switch {
		case strings.HasPrefix(rule, "gt="):
			n, _ := strconv.ParseFloat(rule[3:], 64)
			if v <= n {
				return fmt.Sprintf("%s 必须大于 %.1f", name, n)
			}
		case strings.HasPrefix(rule, "gte="):
			n, _ := strconv.ParseFloat(rule[4:], 64)
			if v < n {
				return fmt.Sprintf("%s 必须 >= %.1f", name, n)
			}
		}
	}
	return ""
}
