package confluence

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Supported options
const (
	_OPTION_UNWRAP   = "unwrap"
	_OPTION_RESPECT  = "respect"
	_OPTION_REVERSE  = "reverse"
	_OPTION_TIMEDATE = "timedate"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// paramsToQuery convert params to query string
func paramsToQuery(params interface{}) string {
	var result string

	t := reflect.TypeOf(params)
	v := reflect.ValueOf(params)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("query")

		switch value.Type().String() {
		case "string":
			result += formatString(tag, value)

		case "int", "int64":
			result += formatInt(tag, value)

		case "bool":
			result += formatBool(tag, value)

		case "time.Time":
			result += formatTime(tag, value)

		case "[]string":
			result += formatSlice(tag, value)
		}
	}

	if result == "" {
		return ""
	}

	return result[:len(result)-1]
}

// formatString returns string representation of string for query string
func formatString(tag string, value reflect.Value) string {
	if value.String() != "" {
		return tag + "=" + esc(value.String()) + "&"
	} else {
		if hasTagOption(tag, _OPTION_RESPECT) {
			return getTagName(tag) + "=&"
		}
	}

	return ""
}

// formatInt returns string representation of int/int64 for query string
func formatInt(tag string, value reflect.Value) string {
	if value.Int() != 0 {
		return tag + "=" + fmt.Sprintf("%d", value.Int()) + "&"
	} else {
		if hasTagOption(tag, _OPTION_RESPECT) {
			return getTagName(tag) + "=0&"
		}
	}

	return ""
}

// formatBool returns string representation of boolean for query string
func formatBool(tag string, value reflect.Value) string {
	b := value.Bool()

	if hasTagOption(tag, _OPTION_REVERSE) && b {
		return getTagName(tag) + "=false&"
	} else {
		if b {
			return getTagName(tag) + "=true&"
		} else {
			if hasTagOption(tag, _OPTION_RESPECT) {
				return getTagName(tag) + "=false&"
			}
		}
	}

	return ""
}

// formatTime returns string representation of time and date for query string
func formatTime(tag string, value reflect.Value) string {
	d := value.Interface().(time.Time)

	if !d.IsZero() {
		if hasTagOption(tag, _OPTION_TIMEDATE) {
			return getTagName(tag) + "=" + d.Format("2006-01-02T15:04:05Z") + "&"
		} else {
			return tag + "=" + d.Format("2006-01-02") + "&"
		}
	}

	return ""
}

// formatSlice returns string representation of slice for query string
func formatSlice(tag string, value reflect.Value) string {
	if value.Len() == 0 {
		return ""
	}

	var result string

	name := getTagName(tag)
	unwrap := hasTagOption(tag, _OPTION_UNWRAP)

	if !unwrap {
		result += name + "="
	}

	for i := 0; i < value.Len(); i++ {
		v := value.Index(i)

		if unwrap {
			result += name + "=" + esc(v.String()) + "&"
		} else {
			result += esc(v.String()) + ","
		}
	}

	return result[:len(result)-1] + "&"
}

// getTagOption extract option from tag
func hasTagOption(tag, option string) bool {
	if !strings.Contains(tag, ",") {
		return false
	}

	return tag[strings.Index(tag, ",")+1:] == option
}

// getTagName return tag name
func getTagName(tag string) string {
	if !strings.Contains(tag, ",") {
		return tag
	}

	return tag[:strings.Index(tag, ",")]
}

// esc escapes the string so it can be safely placed inside a URL query
func esc(s string) string {
	return url.QueryEscape(s)
}
