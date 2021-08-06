package using

import (
	"gorm.io/gorm/clause"
)

type Clause struct {
	sTable   string
	tagParis map[string]interface{}
}

func (i Clause) Build(builder clause.Builder) {
	builder.WriteString("USING ")
	builder.WriteString(i.sTable)
	var tagNameList = make([]string, 0, len(i.tagParis))
	var tagValueList = make([]interface{}, 0, len(i.tagParis))
	for tagName, tagValue := range i.tagParis {
		tagNameList = append(tagNameList, tagName)
		tagValueList = append(tagValueList, tagValue)
	}
	builder.AddVar(builder, tagNameList)
	builder.WriteString(" TAGS")
	builder.AddVar(builder, tagValueList)
}

func SetUsing(sTable string, tags map[string]interface{}) Clause {
	return Clause{
		sTable:   sTable,
		tagParis: tags,
	}
}

func (i Clause) ADDTagPair(tagName string, tagValue interface{}) Clause {
	i.tagParis[tagName] = tagValue
	return i
}

func (i Clause) Name() string {
	return "USING"
}

func (i Clause) MergeClause(c *clause.Clause) {
	c.Name = ""
	c.Expression = i
}
