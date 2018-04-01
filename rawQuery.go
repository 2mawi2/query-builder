package main

type QueryParameter struct {
	tag       string
	condition interface{}
}

type RawQuery struct {
	commandText string
	parameters  []QueryParameter
}

type WhereComparator int

const (
	Equals WhereComparator = iota
	In     WhereComparator = iota
)

type WhereCondition struct {
	column     string
	condition  interface{}
	comparator WhereComparator
}

type QueryOptions struct {
	tableName string
	columns   []string
	wheres    []WhereCondition
}
