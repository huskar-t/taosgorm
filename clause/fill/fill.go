package fill

import (
	"gorm.io/gorm/clause"
	"strconv"
)

type Clause struct {
	value    float64
	fillType Type
}
type Type string

const (
	FillNone   Type = "NONE"
	FillValue  Type = "VALUE"
	FillPrev   Type = "PREV"
	FillNull   Type = "NULL"
	FillLinear Type = "LINEAR"
	FillNext   Type = "NEXT"
)

func (f Clause) Build(builder clause.Builder) {
	builder.WriteString("(")
	builder.WriteString(string(f.fillType))
	if f.fillType == FillValue {
		builder.WriteByte(',')
		builder.WriteString(strconv.FormatFloat(f.value, 'g', -1, 64))
	}
	builder.WriteByte(')')
}

//[FILL(fill_mod_and_val)]

func (f Clause) Name() string {
	return "FILL"
}

func (f Clause) MergeClause(c *clause.Clause) {
	c.Expression = f
}

func SetFill(fillType Type) Clause {
	return Clause{
		fillType: fillType,
	}
}

func (f Clause) SetValue(value float64) Clause {
	f.value = value
	return f
}
