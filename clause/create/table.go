package create

import (
	"bytes"
	"gorm.io/gorm/clause"
	"strconv"
)

type Clause struct {
	tables []*Table
}

const (
	STableType = iota + 1
	CommonTableType
)

type Table struct {
	TableType   int
	Table       string
	IfNotExists bool
	STable      string
	Tags        map[string]interface{}
	Column      []*Column
	TagColumn   []*Column
}

func NewTable(name string, ifNotExist bool, column []*Column, Stable string, tags map[string]interface{}) *Table {
	return &Table{
		TableType:   CommonTableType,
		Table:       name,
		IfNotExists: ifNotExist,
		STable:      Stable,
		Tags:        tags,
		Column:      column,
	}
}

func NewStable(name string, ifNotExist bool, column []*Column, tagColumn []*Column) *Table {
	return &Table{
		TableType:   STableType,
		Table:       name,
		IfNotExists: ifNotExist,
		Column:      column,
		TagColumn:   tagColumn,
	}
}

func CreateTable(tables []*Table) Clause {
	return Clause{tables: tables}
}

func (c Clause) AddTables(tables ...*Table) Clause {
	c.tables = append(c.tables, tables...)
	return c
}

type Column struct {
	Name       string
	ColumnType string
	Length     uint64
}

const (
	TimestampType = "TIMESTAMP"
	IntType       = "INT"
	BigIntType    = "BIGINT"
	FloatType     = "FLOAT"
	DoubleType    = "DOUBLE"
	BinaryType    = "BINARY"
	SmallIntType  = "SMALLINT"
	TinyIntType   = "TINYINT"
	BoolType      = "BOOL"
	NCharType     = "NCHAR"
)

func (c *Column) toSql() string {
	b := bytes.NewBufferString("")
	b.WriteString(c.Name)
	b.WriteByte(' ')
	b.WriteString(c.ColumnType)
	if c.ColumnType == NCharType || c.ColumnType == BinaryType {
		b.WriteByte('(')
		b.WriteString(strconv.FormatUint(c.Length, 10))
		b.WriteByte(')')
	}
	return b.String()
}

func (Clause) Name() string {
	return "CREATE TABLE"
}

func (c Clause) Build(builder clause.Builder) {
	for _, table := range c.tables {
		switch table.TableType {
		case CommonTableType:
			builder.WriteString("CREATE TABLE ")
		case STableType:
			builder.WriteString("CREATE STABLE ")
		default:
			return
		}
		if table.IfNotExists {
			builder.WriteString("IF NOT EXIST ")
		}
		builder.WriteString(table.Table)
		if table.TableType == CommonTableType && table.STable != "" {
			builder.WriteString(" USING ")
			builder.WriteString(table.STable)
			tagValueList := make([]interface{}, 0, len(table.Tags))
			index := 0
			builder.WriteByte('(')
			for tag, tagValue := range table.Tags {
				builder.WriteString(tag)
				if index != len(table.Tags)-1 {
					builder.WriteByte(',')
				}
				tagValueList = append(tagValueList, tagValue)
				index += 1
			}
			builder.WriteString(") TAGS ")
			builder.AddVar(builder, tagValueList)
		}
		builder.WriteString(" (")
		for i, column := range table.Column {
			builder.WriteString(column.toSql())
			if i != len(table.Column)-1 {
				builder.WriteByte(',')
			}
		}
		builder.WriteByte(')')
		if table.TableType == STableType {
			builder.WriteString(" TAGS(")
			for i, tags := range table.TagColumn {
				builder.WriteString(tags.toSql())
				if i != len(table.TagColumn)-1 {
					builder.WriteByte(',')
				}
			}
			builder.WriteByte(')')
		}
	}
}

// MergeClause merge order by clauses
func (c Clause) MergeClause(clause *clause.Clause) {
	clause.Name = ""
	clause.Expression = c
}
