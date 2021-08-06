package expr

import "gorm.io/gorm/clause"

type Expressions []clause.Expression

func (expressions Expressions) Build(builder clause.Builder) {
	for idx, expr := range expressions {
		if idx > 0 {
			builder.WriteByte(' ')
		}
		expr.Build(builder)
	}
}
