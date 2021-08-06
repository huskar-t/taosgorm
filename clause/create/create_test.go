package create_test

import (
	"fmt"
	"github.com/huskar-t/taosgorm/clause/create"
	"github.com/huskar-t/taosgorm/clause/tests"
	"testing"

	"gorm.io/gorm/clause"
)

func TestCreateTable(t *testing.T) {
	var (
		results = []struct {
			Clauses []clause.Interface
			Result  []string
			Vars    [][]interface{}
		}{
			{
				[]clause.Interface{create.CreateTable([]*create.Table{
					{
						TableType:   create.CommonTableType,
						Table:       "t_1",
						IfNotExists: true,
						STable:      "st_1",
						Tags: map[string]interface{}{
							"tag_int":    1,
							"tag_string": "string",
						},
						Column: []*create.Column{
							{
								Name:       "ts",
								ColumnType: create.TimestampType,
								Length:     0,
							}, {
								Name:       "c_int",
								ColumnType: create.IntType,
								Length:     0,
							}, {
								Name:       "c_bigint",
								ColumnType: create.BigIntType,
								Length:     0,
							}, {
								Name:       "c_float",
								ColumnType: create.FloatType,
								Length:     0,
							}, {
								Name:       "c_double",
								ColumnType: create.DoubleType,
								Length:     0,
							}, {
								Name:       "c_binary",
								ColumnType: create.BinaryType,
								Length:     128,
							}, {
								Name:       "c_smallint",
								ColumnType: create.SmallIntType,
								Length:     0,
							}, {
								Name:       "c_tinyint",
								ColumnType: create.TinyIntType,
								Length:     0,
							}, {
								Name:       "c_bool",
								ColumnType: create.BoolType,
								Length:     0,
							}, {
								Name:       "c_nchar",
								ColumnType: create.NCharType,
								Length:     128,
							},
						},
					},
				})},
				[]string{
					"CREATE TABLE IF NOT EXIST t_1 USING st_1(tag_int,tag_string) TAGS (?,?) (ts TIMESTAMP,c_int INT,c_bigint BIGINT,c_float FLOAT,c_double DOUBLE,c_binary BINARY(128),c_smallint SMALLINT,c_tinyint TINYINT,c_bool BOOL,c_nchar NCHAR(128))",
					"CREATE TABLE IF NOT EXIST t_1 USING st_1(tag_string,tag_int) TAGS (?,?) (ts TIMESTAMP,c_int INT,c_bigint BIGINT,c_float FLOAT,c_double DOUBLE,c_binary BINARY(128),c_smallint SMALLINT,c_tinyint TINYINT,c_bool BOOL,c_nchar NCHAR(128))",
				},
				[][]interface{}{{1, "string"}, {"string", 1}},
			},
			{
				[]clause.Interface{create.CreateTable(nil).AddTables(&create.Table{
					TableType:   create.CommonTableType,
					Table:       "t_1",
					IfNotExists: true,
					STable:      "st_1",
					Tags: map[string]interface{}{
						"tag_int":    1,
						"tag_string": "string",
					},
					Column: []*create.Column{
						{
							Name:       "ts",
							ColumnType: create.TimestampType,
							Length:     0,
						}, {
							Name:       "c_int",
							ColumnType: create.IntType,
							Length:     0,
						}, {
							Name:       "c_bigint",
							ColumnType: create.BigIntType,
							Length:     0,
						}, {
							Name:       "c_float",
							ColumnType: create.FloatType,
							Length:     0,
						}, {
							Name:       "c_double",
							ColumnType: create.DoubleType,
							Length:     0,
						}, {
							Name:       "c_binary",
							ColumnType: create.BinaryType,
							Length:     128,
						}, {
							Name:       "c_smallint",
							ColumnType: create.SmallIntType,
							Length:     0,
						}, {
							Name:       "c_tinyint",
							ColumnType: create.TinyIntType,
							Length:     0,
						}, {
							Name:       "c_bool",
							ColumnType: create.BoolType,
							Length:     0,
						}, {
							Name:       "c_nchar",
							ColumnType: create.NCharType,
							Length:     128,
						},
					},
				})},
				[]string{
					"CREATE TABLE IF NOT EXIST t_1 USING st_1(tag_int,tag_string) TAGS (?,?) (ts TIMESTAMP,c_int INT,c_bigint BIGINT,c_float FLOAT,c_double DOUBLE,c_binary BINARY(128),c_smallint SMALLINT,c_tinyint TINYINT,c_bool BOOL,c_nchar NCHAR(128))",
					"CREATE TABLE IF NOT EXIST t_1 USING st_1(tag_string,tag_int) TAGS (?,?) (ts TIMESTAMP,c_int INT,c_bigint BIGINT,c_float FLOAT,c_double DOUBLE,c_binary BINARY(128),c_smallint SMALLINT,c_tinyint TINYINT,c_bool BOOL,c_nchar NCHAR(128))",
				},
				[][]interface{}{{1, "string"}, {"string", 1}},
			},
			{
				[]clause.Interface{create.CreateTable([]*create.Table{
					{
						TableType:   create.STableType,
						Table:       "st_1",
						IfNotExists: true,
						Column: []*create.Column{
							{
								Name:       "ts",
								ColumnType: create.TimestampType,
								Length:     0,
							}, {
								Name:       "c_int",
								ColumnType: create.IntType,
								Length:     0,
							}, {
								Name:       "c_bigint",
								ColumnType: create.BigIntType,
								Length:     0,
							}, {
								Name:       "c_float",
								ColumnType: create.FloatType,
								Length:     0,
							}, {
								Name:       "c_double",
								ColumnType: create.DoubleType,
								Length:     0,
							}, {
								Name:       "c_binary",
								ColumnType: create.BinaryType,
								Length:     128,
							}, {
								Name:       "c_smallint",
								ColumnType: create.SmallIntType,
								Length:     0,
							}, {
								Name:       "c_tinyint",
								ColumnType: create.TinyIntType,
								Length:     0,
							}, {
								Name:       "c_bool",
								ColumnType: create.BoolType,
								Length:     0,
							}, {
								Name:       "c_nchar",
								ColumnType: create.NCharType,
								Length:     128,
							},
						},
						TagColumn: []*create.Column{
							{
								Name:       "t_int",
								ColumnType: create.IntType,
								Length:     0,
							},
						},
					},
				})},
				[]string{
					"CREATE STABLE IF NOT EXIST st_1 (ts TIMESTAMP,c_int INT,c_bigint BIGINT,c_float FLOAT,c_double DOUBLE,c_binary BINARY(128),c_smallint SMALLINT,c_tinyint TINYINT,c_bool BOOL,c_nchar NCHAR(128)) TAGS(t_int INT)",
					"CREATE STABLE IF NOT EXIST st_1 (ts TIMESTAMP,c_int INT,c_bigint BIGINT,c_float FLOAT,c_double DOUBLE,c_binary BINARY(128),c_smallint SMALLINT,c_tinyint TINYINT,c_bool BOOL,c_nchar NCHAR(128)) TAGS(t_int INT)",
				},
				nil,
			},
		}
	)
	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			tests.CheckBuildClauses(t, result.Clauses, result.Result, result.Vars)
		})
	}
}
