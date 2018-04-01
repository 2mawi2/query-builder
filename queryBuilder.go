package main

import (
	"bytes"
	"fmt"
	"strconv"
)

type IQueryBuilder interface {
	Build(options QueryOptions) Query
}

type QueryBuilder struct {
	b       bytes.Buffer
	options QueryOptions
	query   Query
}

type argError struct {
	errorMessage string
}

func (e *argError) Error() string {
	return e.errorMessage
}

func (q QueryBuilder) Build(options QueryOptions) (Query, error) {
	e := validateOptions(options)
	if e != nil {
		return Query{}, e
	}

	q.options = options

	q.addSelect()
	q.addWheres()

	return Query{
		commandText: q.b.String(),
		parameters:  q.query.parameters,
	}, nil
}

func validateOptions(options QueryOptions) error {
	if options.tableName == " " || options.tableName == "" {
		return &argError{errorMessage: "invalid table tag"}
	}
	return nil
}

func (q *QueryBuilder) addSelect() {
	q.b.WriteString("SELECT ")
	if len(q.options.columns) > 0 {
		for i, column := range q.options.columns {
			q.b.WriteString(fmt.Sprintf("[%v]", column))
			if i != len(q.options.columns)-1 {
				q.b.WriteString(", ")
			}
		}
	} else {
		q.b.WriteString("*")
	}
	q.b.WriteString(fmt.Sprintf(" FROM [%v]", q.options.tableName))
}

func (q *QueryBuilder) addWheres() {
	if len(q.options.wheres) > 0 {
		q.b.WriteString(" WHERE ")
		for i, where := range q.options.wheres {
			tag := fmt.Sprintf("arg%v", strconv.Itoa(i))
			q.query.parameters = append(q.query.parameters, QueryParameter{tag: tag, condition: where.condition})
			q.b.WriteString(fmt.Sprintf("[%v].[%v] = @%v", q.options.tableName, where.column, tag))
			if i != len(q.options.wheres)-1 {
				q.b.WriteString(" AND ")
			}
		}
	}
}
