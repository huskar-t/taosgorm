package clause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type IntervalClause struct {
	value uint64
	unit  IntervalUnitType
}

type IntervalUnitType string

//u(微秒)、a(毫秒)、s(秒)、m(分)、h(小时)、d(天)、w(周) n(自然月) 和 y(自然年)
const (
	Microsecond IntervalUnitType = "u"
	Millisecond IntervalUnitType = "a"
	Second      IntervalUnitType = "s"
	Minute      IntervalUnitType = "m"
	Hour        IntervalUnitType = "h"
	Day         IntervalUnitType = "d"
	Week        IntervalUnitType = "w"
	Month       IntervalUnitType = "n"
	Year        IntervalUnitType = "y"
)

func (i IntervalClause) ModifyStatement(stmt *gorm.Statement) {
	c, _ := stmt.Clauses["WHERE"]
	if c.AfterExpression == nil {
		c.AfterExpression = i
	} else if _, ok := c.AfterExpression.(IntervalClause); ok {
		c.AfterExpression = i
	} else {
		c.AfterExpression = Expressions{c.AfterExpression, i}
	}
	stmt.Clauses["WHERE"] = c
}

func (i IntervalClause) Build(builder clause.Builder) {
	if i.unit == "" {
		i.unit = Second
	}
	if i.value == 0 {
		return
	}
	builder.WriteString("interval(")
	builder.WriteString(strconv.FormatUint(i.value, 10))
	builder.WriteString(string(i.unit))
	builder.WriteByte(')')
}

func SetInterval(value uint64, unit IntervalUnitType) IntervalClause {
	return IntervalClause{
		value: value,
		unit:  unit,
	}
}
