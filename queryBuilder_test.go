package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func NewQueryBuilder() QueryBuilder {
	builder := QueryBuilder{}
	return builder
}

func TestQueryBuilder_Build_ShouldBuildSelectCommand(t *testing.T) {
	result, _ := NewQueryBuilder().Build(QueryOptions{
		tableName: "table1",
	})

	assert.Equal(t, "SELECT * FROM [table1]", result.commandText)
}

func TestQueryBuilder_Build_ShouldBuildMultipleSelectCommands(t *testing.T) {
	result, _ := NewQueryBuilder().Build(QueryOptions{
		columns:   []string{"column1", "column2"},
		tableName: "table1",
	})

	assert.Equal(t, "SELECT [column1], [column2] FROM [table1]", result.commandText)
}

func TestQueryBuilder_Build_ShouldValidateTableName(t *testing.T) {
	for _, j := range []string{" ", ""} {
		_, e := NewQueryBuilder().Build(QueryOptions{
			tableName: j,
		})
		assert.EqualError(t, e, "invalid table tag")
	}
}

func TestQueryBuilder_Build_ShouldBuildWhereStatement(t *testing.T) {
	result, _ := NewQueryBuilder().Build(QueryOptions{
		tableName: "table1",
		wheres: []Where{
			{column: "column1", comparator: Equals, condition: 2},
		},
	})

	assert.Equal(t, "SELECT * FROM [table1] WHERE [table1].[column1] = @arg0", result.commandText)
	assert.Equal(t, result.parameters[0].tag, "arg0")
	assert.Equal(t, result.parameters[0].condition, 2)
}

func TestQueryBuilder_Build_ShouldBuildWhereStatement_Multiple(t *testing.T) {
	result, _ := NewQueryBuilder().Build(QueryOptions{
		tableName: "table1",
		wheres: []Where{
			{column: "column1", comparator: Equals, condition: 2},
			{column: "column2", comparator: Equals, condition: 1},
		},
	})

	assert.Equal(t, "SELECT * FROM [table1] WHERE [table1].[column1] = @arg0 AND [table1].[column2] = @arg1", result.commandText)
	assert.Equal(t, result.parameters[0].tag, "arg0")
	assert.Equal(t, result.parameters[0].condition, 2)
	assert.Equal(t, result.parameters[1].tag, "arg1")
	assert.Equal(t, result.parameters[1].condition, 1)
}

func TestQueryBuilder_Build_ShouldBuildWhereStatement_InComparator(t *testing.T) {

}
