package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type QueryBuilder struct {
	b          bytes.Buffer
	options    QueryOptions
	query      RawQuery
	tagCounter int
}

type argError struct {
	errorMessage string
}

func (e *argError) Error() string {
	return e.errorMessage
}

func (q QueryBuilder) Build(options QueryOptions) (RawQuery, error) {
	e := validateOptions(options)
	if e != nil {
		return RawQuery{}, e
	}

	q.options = options
	q.tagCounter = 0

	q.addSelect()
	q.addWheres()

	return RawQuery{
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

func castConditions(where WhereCondition) []interface{} {
	if where.comparator == In {
		return InterfaceSlice(where.condition)
	} else {
		return []interface{}{where.condition}
	}
}

func (q *QueryBuilder) addWheres() {

	if len(q.options.wheres) <= 0 {
		return
	}
	q.b.WriteString(" WHERE ")
	for i, where := range q.options.wheres {
		q.appendWhereCondition(where, i)
	}
}

func (q *QueryBuilder) appendWhereCondition(where WhereCondition, i int) {
	q.b.WriteString(fmt.Sprintf("[%v].[%v] %v ", q.options.tableName, where.column, getComparator(where)))
	conditions := castConditions(where)
	params := q.getParams(conditions)
	q.appendWhereConditions(params, where)
	q.addParamsToQuery(conditions, params)
	if i != len(q.options.wheres)-1 {
		q.b.WriteString(" AND ")
	}
}

func (q *QueryBuilder) addParamsToQuery(conditions []interface{}, params []string) {
	for i := range conditions {
		param := QueryParameter{tag: params[i], condition: conditions[i]}
		q.query.parameters = append(q.query.parameters, param)
	}
}

func (q *QueryBuilder) getParams(conditions []interface{}) []string {
	var params []string
	for range conditions {
		arg := "arg" + strconv.Itoa(q.tagCounter)
		q.tagCounter++
		params = append(params, arg)
	}
	return params
}

func (q *QueryBuilder) appendWhereConditions(params []string, where WhereCondition) {
	var paramsWithAt []string
	for i := range params {
		paramsWithAt = append(paramsWithAt, "@"+params[i])
	}
	whereConditions := strings.Join(paramsWithAt, ", ")
	if where.comparator == In {
		whereConditions = "(" + whereConditions + ")"
	}
	q.b.WriteString(whereConditions)
}

func getComparator(where WhereCondition) string {
	if where.comparator == In {
		return "IN"
	} else if where.comparator == Equals {
		return "="
	} else {
		panic("No such comparator type implemented")
	}
}
