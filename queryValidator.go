package main

import "reflect"

func validateOptions(options QueryOptions) error {
	if options.tableName == " " || options.tableName == "" {
		return &argError{errorMessage: "invalid table tag"}
	}
	for _, w := range options.wheres {
		if w.comparator == In {
			if !(reflect.ValueOf(w.condition).Kind() == reflect.Slice &&
				reflect.ValueOf(w.condition).Kind() != reflect.String) {
				return &argError{errorMessage: "condition must be of type slice when using In comparator"}
			}
		}
	}
	return nil
}
