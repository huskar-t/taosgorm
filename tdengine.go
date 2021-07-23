package taosgorm

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/taosdata/driver-go/taosSql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

// DriverName is the default driver name for SQLite.
const DriverName = "taosSql"

type Dialect struct {
	DriverName string
	DSN        string
	Conn       gorm.ConnPool
}

func Open(dsn string) gorm.Dialector {
	return &Dialect{DSN: dsn}
}

func (dialect Dialect) Name() string {
	return "tdengine"
}

func (dialect Dialect) Initialize(db *gorm.DB) (err error) {
	if dialect.DriverName == "" {
		dialect.DriverName = DriverName
	}
	db.SkipDefaultTransaction = true
	db.DisableNestedTransaction = true
	db.DisableAutomaticPing = true
	db.DisableForeignKeyConstraintWhenMigrating = true
	if dialect.Conn != nil {
		db.ConnPool = dialect.Conn
	} else {
		db.ConnPool, err = sql.Open(dialect.DriverName, dialect.DSN)
		if err != nil {
			return err
		}
	}
	return
}

func (dialect Dialect) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: "NULL"}
}

func (dialect Dialect) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{migrator.Migrator{Config: migrator.Config{
		DB:                          db,
		Dialector:                   dialect,
		CreateIndexAfterCreateTable: false,
	}}, dialect}
}

func (dialect Dialect) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteByte('?')
}

func (dialect Dialect) QuoteTo(writer clause.Writer, str string) {
	writer.WriteString(str)
	return
}

func (dialect Dialect) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `'`, vars...)
}

func (dialect Dialect) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "bool"
	case schema.Int, schema.Uint:
		sqlType := "bigint"
		switch {
		case field.Size <= 8:
			sqlType = "tinyint"
		case field.Size <= 16:
			sqlType = "smallint"
		case field.Size <= 32:
			sqlType = "int"
		}
		return sqlType
	case schema.Float:
		if field.Size <= 32 {
			return "float"
		}
		return "double"
	case schema.String:
		size := field.Size
		if size == 0 {
			size = 64
		}
		return fmt.Sprintf("NCHAR(%d)", size)
	case schema.Time:
		return "TIMESTAMP"
	case schema.Bytes:
		size := field.Size
		if size == 0 {
			size = 64
		}
		return fmt.Sprintf("BINARY(%d)", size)
	}

	return string(field.DataType)
}

func (dialect Dialect) SavePoint(tx *gorm.DB, name string) error {
	return errors.New("not support")
}

func (dialect Dialect) RollbackTo(tx *gorm.DB, name string) error {
	return errors.New("not support")
}
