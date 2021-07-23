package clause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type FillClause struct {
	value    float64
	fillType FillType
}
type FillType string

const (
	FillNone   FillType = "NONE"
	FillValue  FillType = "VALUE"
	FillPrev   FillType = "PREV"
	FillNull   FillType = "NULL"
	FillLinear FillType = "LINEAR"
	FillNext   FillType = "NEXT"
)

func (f FillClause) ModifyStatement(stmt *gorm.Statement) {
	c, _ := stmt.Clauses["WHERE"]
	if c.AfterExpression == nil {
		c.AfterExpression = f
	} else if _, ok := c.AfterExpression.(FillClause); ok {
		c.AfterExpression = f
	} else {
		c.AfterExpression = Expressions{c.AfterExpression, f}
	}
	stmt.Clauses["WHERE"] = c
}

func (f FillClause) Build(builder clause.Builder) {
	builder.WriteString("fill(")
	builder.WriteString(string(f.fillType))
	if f.fillType == FillValue {
		builder.WriteByte(',')
		builder.WriteString(strconv.FormatFloat(f.value, 'g', -1, 64))
	}
	builder.WriteByte(')')
}

func SetFill(fillType FillType) FillClause {
	return FillClause{
		fillType: fillType,
	}
}

func (f *FillClause) SetValue(value float64) FillClause {
	f.value = value
	return *f
}
