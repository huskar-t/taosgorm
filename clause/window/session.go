package window

import (
	"gorm.io/gorm/clause"
	"strconv"
)

const (
	SESSION = iota + 1
	STATE
	INTERVAL
)

//[SESSION(ts_col, tol_val)]
//[STATE_WINDOW(col)]
//[INTERVAL(interval_val [, interval_offset]) [SLIDING sliding_val]]

type Clause struct {
	windowType  int
	tsColumn    string
	stateColumn string
	duration    *Duration
	offset      *Duration
	sliding     *Duration
}

func SetSessionWindow(tsColumn string, duration Duration) Clause {
	return Clause{windowType: SESSION, tsColumn: tsColumn, duration: &duration}
}

func SetStateWindow(column string) Clause {
	return Clause{windowType: STATE, stateColumn: column}
}

func SetInterval(duration Duration) Clause {
	return Clause{windowType: INTERVAL, duration: &duration}
}

func (sc Clause) SetOffset(offset Duration) Clause {
	if sc.windowType == INTERVAL {
		sc.offset = &offset
	}
	return sc
}
func (sc Clause) SetSliding(sliding Duration) Clause {
	if sc.windowType == INTERVAL {
		sc.sliding = &sliding
	}
	return sc
}

func (sc Clause) Build(builder clause.Builder) {
	switch sc.windowType {
	case SESSION:
		builder.WriteString("SESSION(")
		builder.WriteString(sc.tsColumn)
		builder.WriteByte(',')
		builder.WriteString(strconv.FormatUint(sc.duration.Value, 10))
		builder.WriteString(string(sc.duration.Unit))
		builder.WriteByte(')')
	case STATE:
		builder.WriteString("STATE_WINDOW(")
		builder.WriteString(sc.stateColumn)
		builder.WriteByte(')')
	case INTERVAL:
		builder.WriteString("INTERVAL(")
		builder.WriteString(strconv.FormatUint(sc.duration.Value, 10))
		builder.WriteString(string(sc.duration.Unit))
		if sc.offset != nil {
			builder.WriteByte(',')
			builder.WriteString(strconv.FormatUint(sc.offset.Value, 10))
			builder.WriteString(string(sc.offset.Unit))
		}
		builder.WriteByte(')')
		if sc.sliding != nil {
			builder.WriteString(" SLIDING(")
			builder.WriteString(strconv.FormatUint(sc.sliding.Value, 10))
			builder.WriteString(string(sc.sliding.Unit))
			builder.WriteByte(')')
		}
	}
}

func (sc Clause) Name() string {
	return "WINDOW"
}

func (sc Clause) MergeClause(c *clause.Clause) {
	c.Name = ""
	c.Expression = sc
}
