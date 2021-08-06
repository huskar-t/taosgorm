package using_test

import (
	"fmt"
	"github.com/huskar-t/taosgorm/clause/tests"
	"github.com/huskar-t/taosgorm/clause/using"
	"testing"

	"gorm.io/gorm/clause"
)

func TestSetValue(t *testing.T) {
	var (
		results = []struct {
			Clauses []clause.Interface
			Result  []string
			Vars    [][]interface{}
		}{
			{
				Clauses: []clause.Interface{
					clause.Insert{Table: clause.Table{Name: "tb"}},
					using.SetUsing("stb", map[string]interface{}{
						"tag1": 1,
						"tag2": "string",
					}),
				},
				Result: []string{"INSERT INTO tb USING stb(?,?) TAGS(?,?)"},
				Vars: [][]interface{}{
					{"tag1", "tag2", 1, "string"},
				},
			},
		}
	)
	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			tests.CheckBuildClauses(t, result.Clauses, result.Result, result.Vars)
		})
	}
}
