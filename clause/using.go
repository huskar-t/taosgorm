package clause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UsingClause struct {
	sTable string
	tags   []interface{}
}

func (i UsingClause) ModifyStatement(stmt *gorm.Statement) {
	c, _ := stmt.Clauses["INSERT"]
	if c.AfterExpression == nil {
		c.AfterExpression = i
	} else if _, ok := c.AfterExpression.(UsingClause); ok {
		c.AfterExpression = i
	} else {
		c.AfterExpression = Expressions{c.AfterExpression, i}
	}
	stmt.Clauses["INSERT"] = c
}

func (i UsingClause) Build(builder clause.Builder) {
	builder.WriteString("using ")
	builder.WriteString(i.sTable)
	builder.WriteString(" tags")
	builder.AddVar(builder, i.tags)
}

func SetUsing(sTable string, tags []interface{}) UsingClause {
	return UsingClause{
		sTable: sTable,
		tags:   tags,
	}
}
