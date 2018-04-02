package main

import "reflect"

func validateOptions(options QueryOptions) error {
	if isEmpty(options.tableName) {
		return &argError{errorMessage: "invalid table tag"}
	}

	for _, w := range options.wheres {
		if isInvalid(w) {
			return &argError{errorMessage: "condition must be of type slice when using In comparator"}
		}
	}

	return nil
}

func isInvalid(w Where) bool {
	return w.comparator == In && !(isOfTypeSlice(w) && isOfTypeString(w))
}

func isOfTypeString(w Where) bool {
	return reflect.ValueOf(w.condition).Kind() != reflect.String
}

func isOfTypeSlice(w Where) bool {
	return reflect.ValueOf(w.condition).Kind() == reflect.Slice
}

func isEmpty(val string) bool {
	return val == " " || val == ""
}
