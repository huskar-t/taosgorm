# TDengine Gorm 方言

## 使用

```go
package main

import (
	"fmt"
	"github.com/huskar-t/taosgorm"
	"github.com/huskar-t/taosgorm/clause"
	"gorm.io/gorm"
	"time"
)

func main() {
	db, err := gorm.Open(taosgorm.Open("root:taosdata@/tcp(127.0.0.1:6030)/"))
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	//执行自定义语句 demo
	var result []map[string]interface{}
	err = db.Raw("show databases").Scan(&result).Error
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v",result)
	
	//插入时自动建表 demo
	//err = db.Table("t_table").Clauses(clause.SetUsing("s_table", []interface{}{"test1", "test2"})).Create(map[string]interface{}{"ts": time.Now(), "value": 12.0}).Error
	//if err != nil {
	//	panic(err)
	//}
	// //INSERT INTO t_table using s_table tags('test1','test2') (ts,value) VALUES ('2021-07-23 15:56:43.032',12.000000)
	
	//降精度查询 demo
	//now := time.Time{}
	//err = db.Select("avg(value)").Table("t_table").Where("ts > ? and ts <= ?",now.Add(-time.Minute),now).Clauses(clause.SetInterval(1,clause.Day),clause.SetFill(clause.FillNull)).Scan(&result).Error
	// //SELECT avg(value) FROM t_table WHERE ts > '2021-07-23 16:02:56.143' and ts <= '2021-07-23 16:03:56.143' interval(1d) fill(NULL)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v",result)
}
```

## 增加子语句
`using stable tags(?,?,?...)`
使用用例： `db.Clauses(.Clauses(clause.SetUsing("s_table", []interface{}{"test1", "test2"})))`

`interval (?)`
`fill(?)` 
使用用例：`db.Clauses(clause.SetInterval(1,clause.Day),clause.SetFill(clause.FillNull))`
