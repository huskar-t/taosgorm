package main

import (
	"database/sql"
	"github.com/huskar-t/taosgorm"
	"github.com/huskar-t/taosgorm/clause/create"
	"github.com/huskar-t/taosgorm/clause/fill"
	"github.com/huskar-t/taosgorm/clause/using"
	"github.com/huskar-t/taosgorm/clause/window"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"math/rand"
	"time"
)

type Data struct {
	TS    time.Time
	Value float64
}

func main() {
	createDatabase()
	db := connect()
	createSTable(db)
	createTableUsingStable(db)

	now := time.Now()
	randValue := rand.Float64()
	insertData(db, "tb_1", now, randValue)
	t1 := now.Add(time.Second)
	randValue2 := rand.Float64()
	automaticTableCreationWhenInsertingData(db, "tb_2", t1, randValue2)
	tb1Data := queryData(db, "tb_1", now)
	if tb1Data.Value != randValue {
		log.Fatalf("expect value %v got %v", randValue, tb1Data.Value)
	}
	tb2Data := queryData(db, "tb_2", t1)
	if tb2Data.Value != randValue2 {
		log.Fatalf("expect value %v got %v", randValue, tb2Data.Value)
	}
	stbData := queryData(db, "stb_1", now)
	if stbData.Value != randValue {
		log.Fatalf("expect value %v got %v", randValue, stbData.Value)
	}
	t2 := now.Add(time.Second * 2)
	t3 := now.Add(time.Second * 3)
	t4 := now.Add(time.Second * 4)
	v1 := 11
	v2 := 12
	v3 := 13
	automaticTableCreationWhenInsertingMultiData(db, "tb_aggregate", []map[string]interface{}{
		{
			"ts":    t1,
			"value": v1,
		}, {
			"ts":    t2,
			"value": v2,
		}, {
			"ts":    t3,
			"value": v3,
		},
	})
	//aggregate query
	resultAvg := aggregateQuery(db, "tb_aggregate", "avg(value) as v", t1, t3, nil)
	expectAvg := []map[string]interface{}{
		{
			"v": float64(12),
		},
	}
	if !resultMapEqual(expectAvg, resultAvg) {
		log.Fatalf("expect %v got %v", expectAvg, resultAvg)
	}
	windowD, err := window.NewDurationFromTimeDuration(time.Second)
	if err != nil {
		log.Fatal(err)
	}
	resultWindowMax := aggregateQuery(db, "tb_aggregate", "max(value) as v", t1, t4, []clause.Expression{
		window.SetInterval(*windowD),
		fill.SetFill(fill.FillNull),
	})
	expectWindowMax := []map[string]interface{}{
		{
			"ts": now.Add(time.Second),
			"v":  float64(11),
		},
		{
			"ts": now.Add(time.Second * 2),
			"v":  float64(12),
		},
		{
			"ts": now.Add(time.Second * 3),
			"v":  float64(13),
		},
		{
			"ts": now.Add(time.Second * 4),
			"v":  nil,
		},
	}
	if !resultMapEqual(expectWindowMax, resultWindowMax) {
		log.Fatalf("expect %v got %v", expectWindowMax, resultWindowMax)
	}
}

func createDatabase() {
	dsnWithoutDB := "root:taosdata@/tcp(127.0.0.1:6030)/?loc=Local"
	nativeDB, err := sql.Open(taosgorm.DriverName, dsnWithoutDB)
	if err != nil {
		log.Fatalf("connect db error:%v", err)
		return
	}
	_, err = nativeDB.Exec("create database if not exists gorm_test")
	if err != nil {
		log.Fatalf("create database error %v", err)
		return
	}
	_ = nativeDB.Close()
}

func connect() *gorm.DB {
	dsn := "root:taosdata@/tcp(127.0.0.1:6030)/gorm_test?loc=Local"
	db, err := gorm.Open(taosgorm.Open(dsn))
	if err != nil {
		log.Fatalf("unexpected error:%v", err)
	}
	db = db.Debug()
	return db
}

func createSTable(db *gorm.DB) {
	//create stable
	stable := create.NewSTable("stb_1", true, []*create.Column{{
		Name:       "ts",
		ColumnType: create.TimestampType,
	}, {
		Name:       "value",
		ColumnType: create.DoubleType,
	}}, []*create.Column{
		{
			Name:       "tbn",
			ColumnType: create.BinaryType,
			Length:     64,
		},
	})
	err := db.Table("stb_1").Clauses(create.NewCreateTableClause([]*create.Table{stable})).Create(map[string]interface{}{}).Error
	if err != nil {
		log.Fatalf("create sTable error %v", err)
	}
}

func createTableUsingStable(db *gorm.DB) {
	// create table using sTable
	table := create.NewTable("tb_1", true, nil, "stb_1", map[string]interface{}{
		"tbn": "tb_1",
	})
	err := db.Table("tb_1").Clauses(create.NewCreateTableClause([]*create.Table{table})).Create(map[string]interface{}{}).Error
	if err != nil {
		log.Fatalf("create table error %v", err)
	}
}

func insertData(db *gorm.DB, tableName string, ts time.Time, value interface{}) {
	//insert data
	err := db.Table(tableName).Create(map[string]interface{}{
		"ts":    ts,
		"value": value,
	}).Error
	if err != nil {
		log.Fatalf("insert data error %v", err)
	}
}

func automaticTableCreationWhenInsertingData(db *gorm.DB, tableName string, ts time.Time, value interface{}) {
	//automatic table creation when inserting data
	err := db.Table(tableName).Clauses(using.SetUsing("stb_1", map[string]interface{}{
		"tbn": tableName,
	})).Create(map[string]interface{}{
		"ts":    ts,
		"value": value,
	}).Error
	if err != nil {
		log.Fatalf("create table when insert data error %v", err)
	}
}

func queryData(db *gorm.DB, tableName string, ts time.Time) *Data {
	var d Data
	err := db.Table(tableName).Where("ts = ?", ts).Find(&d).Error
	if err != nil {
		log.Fatalf("find data error %v", err)
	}
	return &d
}

func automaticTableCreationWhenInsertingMultiData(db *gorm.DB, tableName string, data []map[string]interface{}) {
	//automatic table creation when inserting data
	err := db.Table(tableName).Clauses(using.SetUsing("stb_1", map[string]interface{}{
		"tbn": tableName,
	})).Create(data).Error
	if err != nil {
		log.Fatalf("create table when insert data error %v", err)
	}
}

func aggregateQuery(db *gorm.DB, tableName string, query string, start, end time.Time, conds []clause.Expression) []map[string]interface{} {
	var result []map[string]interface{}
	err := db.Table(tableName).Select(query).Where("ts >= ? and ts <= ?", start, end).Clauses(conds...).Find(&result).Error
	if err != nil {
		log.Fatalf("aggregate query error %v", err)
	}
	return result
}

func resultMapEqual(m1, m2 []map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false
	}
	for i := range m1 {
		if len(m1[i]) != len(m2[i]) {
			return false
		}

	}
	for i, m := range m1 {
		for s, v := range m {
			_, ok := m2[i][s].(time.Time)
			if ok {
				continue
			}
			if m2[i][s] != v {
				return false
			}
		}
	}
	return true
}
